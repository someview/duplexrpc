package util

import (
	"fmt"
	"rpc-oneway/protocol"
	"sync/atomic"
)

// WriteQuota /*
// WriteQuota is a soft limit on the amount of data a writeWueue can
// schedule before some of it is written To NIC.
type WriteQuota struct {
	quota int64
	// get waits on read from when quota goes less than or equal to zero.
	// replenish writes on it when quota goes positive again.
	ch chan struct{}
	// done is triggered in error case.
	done <-chan struct{}
	// replenish is called by loopyWriter to give quota back to.
	// It is implemented as a field so that it can be updated
	// by tests.
	replenish func(n int)
}

func NewWriteQuota(sz int64, done chan struct{}) *WriteQuota {
	w := &WriteQuota{
		quota: sz,
		ch:    make(chan struct{}, 1),
		done:  done,
	}
	w.replenish = w.realReplenish
	return w
}

func (w *WriteQuota) get(stream protocol.Stream) error {
	for {
		if atomic.LoadInt64(&w.quota) > 0 {
			atomic.AddInt64(&w.quota, -stream.Size)
			return nil
		}
		select {
		case <-w.ch:
			continue
		case <-stream.Done:
			return fmt.Errorf("stream closed")
		}
	}
}

func (w *WriteQuota) realReplenish(n int) {
	sz := int64(n)
	a := atomic.AddInt64(&w.quota, sz)
	b := a - sz
	if b <= 0 && a > 0 {
		select {
		case w.ch <- struct{}{}:
		default:
		}
	}
}
