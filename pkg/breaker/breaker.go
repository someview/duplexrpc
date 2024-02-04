package breaker

// import (
// 	"sync"
// 	"time"
// )

// // 参考资料 circuitBreaker
// // https://talkgo.org/t/topic/3035
// // go-zero
// type Acceptable func(err error) bool

// // 手动回调
// type Promise interface {
// 	// Accept tells the Breaker that the call is successful.
// 	// 请求成功
// 	Accept()
// 	// Reject tells the Breaker that the call is failed.
// 	// 请求失败
// 	Reject(reason string)
// }

// type Breaker interface {
// 	Allow() (Promise, error)
// 	// 熔断器，自动执行上报结果
// 	Do(req func() error) error
// 	// 支持自定义执行结果
// 	DoWithAcceptable(req func() error, acceptable Acceptable) error
// 	// 支持自定义快速执行失败
// 	DoWithFallback(req func() error, fallback func(err error) error) error

// 	// 熔断方法
// 	// fallback - 支持自定义快速失败
// 	// acceptable - 支持自定义判定执行结果
// 	DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
// }

// type googleBreaker struct {
// 	// 敏感度，建议值是1.5-2.0, 这里采取默认值是1.5
// 	k float64
// }

// type rollingWindow struct {
// 	// 滑动窗口的锁
// 	lock sync.RWMutex
// 	// 滑动窗口数量
// 	size int
// 	// 滑动窗口单元时间间隔
// 	interval time.Duration
// 	// lastTime,最后写入桶的时间
// 	lastTime time.Duration
// 	// 游标，用于定位当前应该写入哪个bucket
// 	offSet int

// 	buckets []*Buckets
// }

// type Buckets struct {
// 	// 一个桶标识一个时间间隔
// }

// // 添加数据
// // offset - 游标，定位写入bucket位置
// // v - 行为数据
// func (w *rollingWindow) add(offset int, v float64) {
// 	w.buckets[offset%w.size].add(v)
// }

// // 汇总数据
// // fn - 自定义的bucket统计函数
// func (w *rollingWindow) reduce(start, count int, fn func(b *Bucket)) {
// 	for i := 0; i < count; i++ {
// 		fn(w.buckets[(start+i)%w.size])
// 	}
// }

// // 清理特定bucket
// func (w *window) resetBucket(offset int) {
// 	w.buckets[offset%w.size].reset()
// }

// // 桶
// type Bucket struct {
// 	// 当前桶内值之和
// 	Sum float64
// 	// 当前桶的add总次数
// 	Count int64
// }

// // 向桶添加数据
// func (b *Bucket) add(v float64) {
// 	// 求和
// 	b.Sum += v
// 	// 次数+1
// 	b.Count++
// }

// // 桶数据清零
// func (b *Bucket) reset() {
// 	b.Sum = 0
// 	b.Count = 0
// }
