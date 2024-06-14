package util

import (
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.gobase.git/dt/pool"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
)

var bufPool = pool.NewSyncPool[*netpoll.LinkBuffer](func() any {
	lb := new(netpoll.LinkBuffer)
	return lb
})

// SliceBuf 从Buf中Slice出一个Buf
func SliceBuf(n int, slicer netpoll.Sliceable) (netpoll.Reader, error) {
	buf := bufPool.Get()
	err := slicer.SliceIntoReader(n, buf)
	if err != nil {
		bufPool.Put(buf)
		return nil, err
	}
	return buf, nil
}

// PutBufFromSlice 把Slice出的Buf放回Buf池
func PutBufFromSlice(reader netpoll.Reader) error {
	buf, ok := reader.(*netpoll.LinkBuffer)
	if !ok {
		return fmt.Errorf("reader is not a *netpoll.LinkBuffer")
	}
	err := buf.Skip(buf.Len())
	if err != nil {
		return err
	}
	err = buf.Release()
	if err != nil {
		return err
	}
	bufPool.Put(buf)
	return nil
}
