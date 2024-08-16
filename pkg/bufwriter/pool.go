package bufwriter

import (
	"math/bits"
	"sync"

	netpoll "github.com/cloudwego/netpoll"
)

const (
	levelNum = 5
	baseBit  = 4
	minSize  = 1 << baseBit
)

// index contains []byte which cap is 1<<index
var caches [levelNum]sync.Pool

func init() {
	for i := 0; i < levelNum; i++ {
		size := calcSize(i)
		caches[i].New = func() interface{} {
			buf := make([]*netpoll.LinkBuffer, 0, size)
			return buf
		}
	}
}

func calcSize(level int) int {
	return 1 << (baseBit + level)
}

// calculates which pool to get from
// res: 0-maxlevel maxlevel + 1
func calcLevel(size int) int {
	if size <= minSize {
		return 0
	}
	if isPowerOfTwo(size) {
		return bsr(size) - baseBit
	}
	return bsr(size) - baseBit + 1
}
func bsr(x int) int {
	return bits.Len(uint(x)) - 1
}

func isPowerOfTwo(x int) bool {
	return (x & (-x)) == x
}

// MallocBufSlice 根据指定参数分配固定容量的linkBuffer slice
func MallocBufSlice(size int) []*netpoll.LinkBuffer {
	index := calcLevel(size)
	if index >= levelNum {
		return make([]*netpoll.LinkBuffer, 0, size)
	}
	var ret = caches[index].Get().([]*netpoll.LinkBuffer)
	return ret
}

// FreeBufSlice Free should be called when the buf is no longer used.
func FreeBufSlice(buf []*netpoll.LinkBuffer) {
	index := calcLevel(cap(buf))
	if index >= levelNum {
		return
	}
	buf = buf[:0]
	//nolint:staticcheck
	caches[index].Put(buf)
}
