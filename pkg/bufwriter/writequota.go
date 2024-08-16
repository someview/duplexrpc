package bufwriter

import (
	"context"
	"errors"
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

// writeQuota is a soft limit on the amount of data a writer can
// schedule before some of it is written out.
type writeQuota struct {
	quota int32
	// get waits on read from when quota goes less than or equal to zero.
	// replenish writes on it when quota goes positive again.
	ch chan struct{}
	// done is triggered in error case.,单向关闭信号
	done <-chan struct{}
	// replenish is called by loopyWriter to give quota back to.
	// It is implemented as a field so that it can be updated
	// by tests.
	replenish func(n int)
}

var ErrNoQuota = errors.New("no quota available")

func newWriteQuota(sz int32, done <-chan struct{}) *writeQuota {
	w := &writeQuota{
		quota: sz,
		ch:    make(chan struct{}, 1),
		done:  done,
	}

	w.replenish = w.realReplenish

	return w
}

// getWithCtx 尝试获取写入配额，如果配额不足，则等待直到有足够的配额或上下文被取消。
// 参数:
//
//	ctx: 上下文对象，用于检测操作是否应该因为超时或取消而取消。
//	sz: 需要获取的配额大小
//
// 返回值:
//
//	错误: 如果上下文被取消或连接关闭，则返回相应的错误。
func (w *writeQuota) getWithCtx(ctx context.Context, sz int32) error {
	// 不断尝试获取配额直到成功或上下文被取消。
	for {
		// 检查当前配额是否足够。
		if atomic.LoadInt32(&w.quota) > 0 {
			// 减少配额并尝试发送一个空结构体到通道，表示配额已被使用。
			atomic.AddInt32(&w.quota, -sz)
			select {
			// 如果通道不忙，发送一个空结构体。
			case w.ch <- struct{}{}:
			// 如果通道已满，直接继续，不需要等待。
			default:
			}
			// 配额获取成功，返回nil。
			return nil
		}
		// 检查上下文是否已被取消。
		select {
		case <-ctx.Done():
			// 如果上下文被取消，返回相应的错误。
			return ctx.Err()
		// 等待配额通道释放配额。
		case <-w.ch:
			continue
		// 如果done通道关闭，表示连接已关闭，返回错误。
		case <-w.done:
			return fmt.Errorf("writequota closed")
		}
	}
}

func (w *writeQuota) realReplenish(n int) {
	sz := int32(n)
	a := atomic.AddInt32(&w.quota, sz)
	b := a - sz
	// 增加前配额不足,增加后配额充足的情况下，通知获取配额方
	if b <= 0 && a > 0 {
		select {
		case w.ch <- struct{}{}:
		default:
		}
	}
}

func (w *writeQuota) GetWithCtx(ctx context.Context, sz int) error {
	return w.getWithCtx(ctx, int32(sz))
}

func (w *writeQuota) TryGet(sz int) bool {
	if atomic.LoadInt32(&w.quota) > 0 {
		atomic.AddInt32(&w.quota, int32(-sz))
		select {
		case w.ch <- struct{}{}:
		default:
		}
		return true
	}
	return false
}

func (w *writeQuota) Release(sz int) {
	w.realReplenish(sz)
}

func (w *writeQuota) Available() int {
	return int(atomic.LoadInt32(&w.quota))
}
