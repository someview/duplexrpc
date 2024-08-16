package bufwriter

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"time"

	netpoll "github.com/cloudwego/netpoll"
)

var ErrBatchWriterClosed = errors.New("batchWriter has been closed")

type batchWriter struct {
	bc BatchContainer

	timer      *time.Timer   // flush timer
	flushDelay time.Duration // flush间隔
	notifyChan chan struct{} // flush通知
	isRunning  atomic.Bool

	state atomic.Int32 // 0:active 1:closing 2:closed
	done  chan struct{}

	conn netpoll.Connection

	flushCb func(n int, err error) // flush后执行的回调函数
}

// FlushTo implements BatchWriter.

// NewBatchWriter creates a BatchWriter.
//
//	maxBufferSize: the max size of one batch.
//	flushDelay: the flush delay.
//
// 自动flush数据到connection中
func NewBatchWriter(maxBufferSize int, maxMsgNum int, flushDelay time.Duration,
	conn netpoll.Connection, flushCb func(n int, err error)) BatchWriter {
	done := make(chan struct{}, 1)
	w := &batchWriter{
		flushDelay: flushDelay,
		notifyChan: make(chan struct{}, 1),
		timer:      time.NewTimer(flushDelay),
		done:       done,
		conn:       conn,
		bc:         newBatchContainer(maxBufferSize, maxMsgNum),
		flushCb:    flushCb,
	}
	return w
}

func (b *batchWriter) Add(ctx context.Context, lb *netpoll.LinkBuffer) error {
	if b.state.Load() != active {
		return ErrBatchWriterClosed
	}
	if err := b.bc.Add(ctx, lb); err != nil {
		return err
	}
	b.tryNotifyFlush()
	b.tryRunFlushLoop()
	return nil
}

func (b *batchWriter) Close() error {
	if !b.state.CompareAndSwap(active, closing) {
		return ErrBatchWriterClosed
	}
	// 1. 关闭BatchContainer，唤醒所有AsyncWrite阻塞的协程
	// 2. 关闭b.done，唤醒flush协程
	_ = b.bc.Close()
	close(b.done)

	// wait for all tasks finished
	for b.state.Load() != closed {
		if b.bc.Len() == 0 {
			b.state.Store(closed)
			break
		}

		runtime.Gosched()
	}
	return nil

}

func (b *batchWriter) flushLoop() {

	for !b.bc.IsEmpty() {
		b.waitFlush()
		if err := b.flushData(); err != nil {
			return
		}
	}

	b.isRunning.Store(false)
	// 退出前检查
	if !b.bc.IsEmpty() {
		b.tryRunFlushLoop()
		return
	}

	b.state.CompareAndSwap(closing, closed)
}

func (b *batchWriter) flushData() error {
	n, err := b.bc.FlushTo(b.conn)
	defer b.flushCb(n, err)
	if err != nil {
		_ = b.conn.Close()
	}
	return err
}

// 检查一下是否需要立即flush
func (b *batchWriter) tryNotifyFlush() {
	if b.bc.IsFull() {
		select {
		case b.notifyChan <- struct{}{}:
		default:
		}
	}
}

func (b *batchWriter) waitFlush() {
	b.timer.Reset(b.flushDelay)

	// 三种情况会进行flush
	// 1. notifyFlush 主动通知
	// 2. 到时间
	// 3. 调用Close()
	select {
	case <-b.notifyChan:
	case <-b.timer.C:
	case <-b.done:
	}
}

func (b *batchWriter) tryRunFlushLoop() {
	if b.isRunning.CompareAndSwap(false, true) {
		netpoll.Go(b.flushLoop)
	}
}
