package rpcinfo

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/internal"
	"sync"
)

var rpcInfoPool sync.Pool

type rpcInfo struct {
	endpointCtx     context.Context
	invocation      Invocation
	interactionMode InteractionMode
}

func (r *rpcInfo) InteractionMode() InteractionMode {
	return r.interactionMode
}

func (r *rpcInfo) EndpointContext() context.Context {
	return r.endpointCtx
}

func (r *rpcInfo) Invocation() Invocation { return r.invocation }

func (r *rpcInfo) zero() {
	r.endpointCtx = nil
	r.invocation = nil
	r.interactionMode = 0
}

func (r *rpcInfo) Recycle() {
	if v, ok := r.invocation.(internal.Reusable); ok {
		v.Recycle()
	}
	r.zero()
	rpcInfoPool.Put(r)
}

func init() {
	rpcInfoPool.New = newRPCInfo
}

func newRPCInfo() interface{} {
	return &rpcInfo{}
}

func NewRPCInfo(interactionMode InteractionMode, ink Invocation) RPCInfo {
	r := rpcInfoPool.Get().(*rpcInfo)
	r.invocation = ink
	r.interactionMode = interactionMode
	return r
}

func NewEmptyRPCInfo(ink Invocation) RPCInfo {
	r := rpcInfoPool.Get().(*rpcInfo)
	r.invocation = ink
	return r
}
