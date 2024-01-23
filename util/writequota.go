package util

import (
	"fmt"
	"sync/atomic"
)

/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
// writeQuota is a soft limit on the amount of data a writeWueue can
// schedule before some of it is written To NIC.
type writeQuota struct {
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

func newWriteQuota(sz int64, done <-chan struct{}) *writeQuota {
	w := &writeQuota{
		quota: sz,
		ch:    make(chan struct{}, 1),
		done:  done,
	}
	w.replenish = w.realReplenish
	return w
}

func (w *writeQuota) get(sz int64) error {
	for {
		if atomic.LoadInt64(&w.quota) > 0 {
			atomic.AddInt64(&w.quota, -sz)
			return nil
		}
		select {
		case <-w.ch:
			continue
		case <-w.done:
			return fmt.Errorf("connection closed")
		}
	}
}

func (w *writeQuota) realReplenish(n int) {
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
