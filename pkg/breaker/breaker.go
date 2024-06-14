package breaker

import (
	"sync"
	"time"
)

type Breaker interface {
	// 本次请求是否通过
	Allow() bool
	// 记录失败状态
	Fail()
}

// todo 多核情况下可以更进一步使用shardPbreaker
type breaker struct {
	mu               sync.RWMutex
	failureThreshold uint64        // 失败阈值
	window           time.Duration // 窗口大小

	lastFailureTime time.Time // 最后一次失败时间
	failures        uint64    // 当前失败次数
	open            bool      // 是否是打开状态
}

func NewBreaker(failureThreshold uint64, window time.Duration) *breaker {
	return &breaker{
		failureThreshold: failureThreshold,
		window:           window,
	}
}

func (b *breaker) Allow() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return !b.open
}

// todo 引入时间轮，来控制time.Now()的调用此时
func (b *breaker) Fail() {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	added := b.lastFailureTime.Add(b.window)
	inWindow := now.Before(added)
	// 如果当前时间超出了上一次失败时间加上时间窗口，则重置失败计数和更新熔断器为ready状态
	if !inWindow && b.open {
		b.failures = 0
		b.open = false
	} else {

		// 如果在时间窗口内且失败次数即将超过阈值，则设置熔断器为非ready状态
		if b.failures+1 > b.failureThreshold {
			b.open = true
			return
		}
		// 无论熔断器状态如何，都记录失败
		b.failures++
		b.lastFailureTime = now
	}
}
