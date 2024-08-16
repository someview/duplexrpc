package bufwriter

import (
	"context"
	"sync"
	"sync/atomic"

	netpoll "github.com/cloudwego/netpoll"
)

type batchContainer struct {
	maxMsgNum int32
	length    atomic.Int32

	mu   sync.Mutex
	bufs []*netpoll.LinkBuffer

	fc   FlowControl
	done chan struct{}
}

func newBatchContainer(maxBufferSize int, maxMsgNum int) BatchContainer {
	done := make(chan struct{}, 1)
	b := &batchContainer{
		maxMsgNum: int32(maxMsgNum),
		fc:        newWriteQuota(int32(maxBufferSize), done),
		bufs:      MallocBufSlice(maxMsgNum),
		done:      done,
	}

	return b
}

func (b *batchContainer) Len() int {
	return int(b.length.Load())
}

func (b *batchContainer) IsEmpty() bool {
	return b.length.CompareAndSwap(0, 0)
}

func (b *batchContainer) IsFull() bool {
	return b.length.Load() >= b.maxMsgNum || b.fc.Available() <= 0
}

func (b *batchContainer) Add(ctx context.Context, lb *netpoll.LinkBuffer) error {
	err := b.fc.GetWithCtx(ctx, lb.Len())
	if err != nil {
		return err
	}
	b.append(lb)
	return nil
}

func (b *batchContainer) TryAdd(lb *netpoll.LinkBuffer) bool {

	if !b.fc.TryGet(lb.Len()) {
		return false
	}
	b.append(lb)
	return true
}

func (b *batchContainer) append(lb *netpoll.LinkBuffer) {
	b.mu.Lock()
	b.bufs = append(b.bufs, lb)
	b.mu.Unlock()
	b.length.Add(1)
}

func (b *batchContainer) FlushTo(conn netpoll.Connection) (int, error) {
	wr := conn.Writer()
	mallocBufs := MallocBufSlice(int(b.maxMsgNum))
	b.mu.Lock()
	bufs := b.bufs
	b.bufs = mallocBufs
	b.mu.Unlock()
	var count int
	for i := 0; i < len(bufs); i++ {
		buf := bufs[i]
		b.length.Add(-1)
		count += buf.Len()
		_ = wr.Append(buf)
		buf.Recycle()
	}
	defer FreeBufSlice(bufs)
	err := wr.Flush()
	b.fc.Release(count)
	return count, err
}

// Close 并发不安全，方法仅能调用1次,多次调用会panic
func (b *batchContainer) Close() error {
	close(b.done)
	return nil
}
