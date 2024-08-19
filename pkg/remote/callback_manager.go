package remote

import (
	netpoll "github.com/cloudwego/netpoll"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/someview/dt/pool"
)

var readerChanPool = pool.NewSyncPool[chan netpoll.Reader](func() any {
	return make(chan netpoll.Reader, 1)
})

type callBackManager struct {
	mp cmap.ConcurrentMap[uint32, chan netpoll.Reader]
}

func newCallBackManager() CallBackManager {
	return &callBackManager{
		mp: cmap.NewWithCustomShardingFunction[uint32, chan netpoll.Reader](func(key uint32) uint32 {
			return key
		}),
	}

}

func (c *callBackManager) Set(seqID uint32, reader chan netpoll.Reader) {
	c.mp.Set(seqID, reader)
}

func (c *callBackManager) Delete(seqID uint32) {
	c.mp.Remove(seqID)
}

func (c *callBackManager) Load(seqID uint32) (reader chan netpoll.Reader, ok bool) {
	return c.mp.Get(seqID)
}
func (c *callBackManager) LoadAndDelete(seqID uint32) (reader chan netpoll.Reader, ok bool) {
	return c.mp.Pop(seqID)
}

func (c *callBackManager) IsEmpty() bool {
	return c.mp.IsEmpty()
}
