package util

import (
	"fmt"

	netpoll "github.com/cloudwego/netpoll"

	"github.com/someview/dt/pool"
)

var bufPool = pool.NewSyncPool[*netpoll.LinkBuffer](func() any {
	lb := new(netpoll.LinkBuffer)
	return lb
})

type Sliceable interface {
	SliceInto(n int, r netpoll.Reader) error
}

// SliceBuf 从Buf中Slice出一个Buf
func SliceBuf(n int, slicer Sliceable) (netpoll.Reader, error) {
	buf := bufPool.Get()
	err := slicer.SliceInto(n, buf)
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
