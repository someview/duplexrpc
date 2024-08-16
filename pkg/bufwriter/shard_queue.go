// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufwriter

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"

	netpoll "github.com/cloudwego/netpoll"
)

/* DOC:
 * ShardQueue uses the netpoll's nocopy API to merge and send data.
 * The Data Flush is passively triggered by ShardQueue.Add and does not require user operations.
 * If there is an error in the data transmission, the connection will be closed.
 *
 * ShardQueue.Add: add the data to be sent.
 * NewShardQueue: create a queue with netpoll.Connection.
 * ShardSize: the recommended number of shards is 32.
 */
var ShardSize int

func init() {
	ShardSize = runtime.GOMAXPROCS(0)
}

var ErrShardQueueClosed = errors.New("shardQueue has been closed")
var ErrDataEmpty = errors.New("data is empty")

// NewShardQueue
// quota可以限制queue的最大内存配额
// 推荐quota为data大小的10-40倍左右
func NewShardQueue(quota int, conn netpoll.Connection) BatchWriter {
	return newShardQueue(ShardSize, quota, conn)
}

// ShardQueue uses the netpoll's nocopy API to merge and send data.
// The Data Flush is passively triggered by ShardQueue.Add and does not require user operations.
// If there is an error in the data transmission, the connection will be closed.
// ShardQueue.Add: add the data to be sent.
type ShardQueue struct {
	conn      netpoll.Connection
	idx, size int32
	getters   [][]*netpoll.LinkBuffer // len(getters) = size
	swap      []*netpoll.LinkBuffer   // use for swap
	locks     []int32                 // len(locks) = size

	queueTrigger

	fc   FlowControl
	done chan struct{}
}

// here for trigger
type queueTrigger struct {
	trigger  int32
	state    int32 // 0: active, 1: closing, 2: closed
	runNum   int32
	w, r     int32      // ptr of list
	list     []int32    // record the triggered shard
	listLock sync.Mutex // list total lock
}

// NewShardQueue .
func newShardQueue(size int, quota int, conn netpoll.Connection) (queue *ShardQueue) {

	queue = &ShardQueue{
		conn:    conn,
		size:    int32(size),
		getters: make([][]*netpoll.LinkBuffer, size),
		swap:    make([]*netpoll.LinkBuffer, 0, 64),
		locks:   make([]int32, size),
		done:    make(chan struct{}),
	}
	for i := range queue.getters {
		queue.getters[i] = make([]*netpoll.LinkBuffer, 0, 64)
	}

	queue.list = make([]int32, size)

	queue.fc = newWriteQuota(int32(quota), queue.done)

	return queue
}

func (q *ShardQueue) Add(ctx context.Context, lb *netpoll.LinkBuffer) error {
	// check if queue is closed
	if atomic.LoadInt32(&q.state) != active {
		return ErrShardQueueClosed
	}

	// get quota
	err := q.fc.GetWithCtx(ctx, lb.Len())
	if err != nil {
		return err
	}

	// add to queue
	q.add(lb)
	return nil
}

// Add adds to q.getters[shard]
func (q *ShardQueue) add(gts *netpoll.LinkBuffer) {
	if atomic.LoadInt32(&q.state) != active {
		return
	}
	shard := atomic.AddInt32(&q.idx, 1) % q.size
	q.lock(shard)
	trigger := len(q.getters[shard]) == 0
	q.getters[shard] = append(q.getters[shard], gts)
	q.unlock(shard)
	if trigger {
		q.triggering(shard)
	}
}

func (q *ShardQueue) Close() error {
	if !atomic.CompareAndSwapInt32(&q.state, active, closing) {
		return ErrShardQueueClosed
	}

	// 关闭wq
	close(q.done)

	// wait for all tasks finished
	for atomic.LoadInt32(&q.state) != closed {
		if atomic.LoadInt32(&q.trigger) == 0 {
			atomic.StoreInt32(&q.state, closed)
			return nil
		}
		runtime.Gosched()
	}
	return nil
}

// triggering shard.
func (q *ShardQueue) triggering(shard int32) {
	q.listLock.Lock()
	q.w = (q.w + 1) % q.size
	q.list[q.w] = shard
	q.listLock.Unlock()

	if atomic.AddInt32(&q.trigger, 1) > 1 {
		return
	}
	q.foreach()
}

// foreach swap r & w. It's not concurrency safe.
func (q *ShardQueue) foreach() {
	if atomic.AddInt32(&q.runNum, 1) > 1 {
		return
	}
	netpoll.Go(q.processQueue)
}

func (q *ShardQueue) processQueue() {

	var negNum int32 // is negative number of triggerNum
	var n int        // appended bytes
	for triggerNum := atomic.LoadInt32(&q.trigger); triggerNum > 0; {
		q.r = (q.r + 1) % q.size
		shared := q.list[q.r]

		// lock & swap
		q.lock(shared)
		tmp := q.getters[shared]
		q.getters[shared] = q.swap[:0]
		q.swap = tmp
		q.unlock(shared)

		// deal
		n = q.deal(q.swap)
		negNum--
		if triggerNum+negNum == 0 {
			triggerNum = atomic.AddInt32(&q.trigger, negNum)
			negNum = 0
		}
	}
	q.flush()
	q.fc.Release(n)

	// quit & check again
	atomic.StoreInt32(&q.runNum, 0)
	if atomic.LoadInt32(&q.trigger) > 0 {
		q.foreach()
		return
	}
	// if state is closing, change it to closed
	atomic.CompareAndSwapInt32(&q.state, closing, closed)

}

// deal is used to get deal of netpoll.Writer.
func (q *ShardQueue) deal(bufs []*netpoll.LinkBuffer) (n int) {
	writer := q.conn.Writer()

	for _, buf := range bufs {
		n += buf.Len()
		// 此处的writer为conn底层的writer
		err := writer.Append(buf)
		buf.Recycle()
		if err != nil {
			q.conn.Close()
			return
		}
	}
	return
}

// flush is used to flush netpoll.Writer.
func (q *ShardQueue) flush() {
	err := q.conn.Writer().Flush()
	if err != nil {
		q.conn.Close()
		return
	}
}

// lock shard.
func (q *ShardQueue) lock(shard int32) {
	for !atomic.CompareAndSwapInt32(&q.locks[shard], 0, 1) {
		runtime.Gosched()
	}
}

// unlock shard.
func (q *ShardQueue) unlock(shard int32) {
	atomic.StoreInt32(&q.locks[shard], 0)
}
