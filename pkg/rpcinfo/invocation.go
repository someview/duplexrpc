package rpcinfo

import (
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/uerror"
	"sync"
	"sync/atomic"
)

var (
	_                   Invocation       = (*invocation)(nil)
	_                   InvocationSetter = (*invocation)(nil)
	invocationPool      sync.Pool
	defaultSeqGenerator = &seqGenerator{}
)

func init() {
	invocationPool.New = newInvocation
}

type seqGenerator struct {
	seq atomic.Int32
}

func (g *seqGenerator) Next() uint32 {
	id := g.seq.Add(1)
	if id == 0 {
		id = g.seq.Add(1)
	}
	return uint32(id)
}

type invocation struct {
	flag

	methodInfo  service.MethodInfo
	serviceImpl any
	packageName string
	serviceName string
	methodName  string
	seqID       uint32
	bizErr      uerror.BizError
	extra       map[string]interface{}
}

func (i *invocation) SetMethodInfo(info service.MethodInfo) {
	i.methodInfo = info
}

func (i *invocation) SetServiceImpl(impl any) {
	i.serviceImpl = impl
}

func (i *invocation) MethodInfo() service.MethodInfo {
	return i.methodInfo
}

func (i *invocation) ServiceImpl() any {
	return i.serviceImpl
}

func (i *invocation) FlagInfo() FlagInfo {
	return &i.flag
}

func (i *invocation) SetBizError(err uerror.BizError) {
	i.bizErr = err
}

func (i *invocation) BizError() uerror.BizError {
	return i.bizErr
}

// NewInvocation creates a new Invocation with the given service, method and optional package.
func NewInvocation(service, method string) *invocation {
	ivk := invocationPool.Get().(*invocation)
	ivk.seqID = defaultSeqGenerator.Next()
	ivk.serviceName = service
	ivk.methodName = method
	return ivk
}

// NewEmptyInvocation to get Invocation for new request in server side
func NewEmptyInvocation() *invocation {
	ivk := invocationPool.Get().(*invocation)
	return ivk
}

func newInvocation() interface{} {
	return &invocation{}
}

// SeqID implements the Invocation interface.
func (i *invocation) SeqID() uint32 {
	return i.seqID
}

// SetSeqID implements the InvocationSetter interface.
func (i *invocation) SetSeqID(seqID uint32) {
	i.seqID = seqID
}

func (i *invocation) PackageName() string {
	return i.packageName
}

func (i *invocation) SetPackageName(name string) {
	i.packageName = name
}

func (i *invocation) ServiceName() string {
	return i.serviceName
}

// SetServiceName implements the InvocationSetter interface.
func (i *invocation) SetServiceName(name string) {
	i.serviceName = name
}

// MethodName implements the Invocation interface.
func (i *invocation) MethodName() string {
	return i.methodName
}

// SetMethodName implements the InvocationSetter interface.
func (i *invocation) SetMethodName(name string) {
	i.methodName = name
}

// BizStatusErr implements the Invocation interface.

func (i *invocation) SetExtra(key string, value interface{}) {
	if i.extra == nil {
		i.extra = map[string]interface{}{}
	}
	i.extra[key] = value
}

func (i *invocation) Extra(key string) interface{} {
	if i.extra == nil {
		return nil
	}
	return i.extra[key]
}

// Reset implements the InvocationSetter interface.
func (i *invocation) Reset() {
	i.zero()
}

// Recycle reuses the invocation.
func (i *invocation) Recycle() {
	i.zero()
	invocationPool.Put(i)
}

func (i *invocation) zero() {
	i.flag.zero()
	i.methodInfo = nil
	i.serviceImpl = nil
	i.seqID = 0
	i.packageName = ""
	i.serviceName = ""
	i.methodName = ""
	i.bizErr = nil
	for key := range i.extra {
		delete(i.extra, key)
	}
}
