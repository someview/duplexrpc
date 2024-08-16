package testdata

import (
	"context"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery/resolver"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
)

type PubEmitter interface {
	Pub(ctx context.Context, req *PubReq) error
	SubFailed(ctx context.Context, req *SubFailedReq) error
	ExistSub(ctx context.Context, req *ExistSubReq, res *ExistSubRes) error
}

type SubService interface {
	SubBroadcast(ctx context.Context, req *SubReq) error
	UnSubBroadcast(ctx context.Context, req *UnsubReq) error
	ExistSubscription(ctx context.Context, req *ExistSubReq, res *ExistSubRes) error
	SyncSnapshot(ctx context.Context, req *SyncSnapshotReq) error
}

// ConsumerServiceClient Client应该理解为端点
type ConsumerServiceClient interface {
	SubService
	SelfServiceType() PubEmitter
}

// ConsumerServiceServer self可以考虑
type ConsumerServiceServer interface {
	PubEmitter
	SelfServiceType() SubService
}

const (
	MethodPub       = "Pub"
	MethodSubFailed = "SubFailed"
	MethodExistSub  = "ExistSub"
)

const (
	MethodSubBroadcast      = "SubBroadcast"
	MethodUnSubBroadcast    = "UnsubBroadcast"
	MethodExistSubscription = "ExistSubscription"
	MethodSyncSnapshot      = "SyncSnapshot"
)

// 用户层可以自定NewArg,NewResp来实现内存复用
// todo 提供一个用户层设置内存复用的机会

var subServiceInfo = service.NewServiceInfo("PubService", (*SubService)(nil),
	map[string]service.MethodInfo{
		MethodSubBroadcast:      service.NewMethodInfo(SubBroadcastHandler, service.TArgsFactory[SubReq]{}, nil, true, false),
		MethodUnSubBroadcast:    service.NewMethodInfo(UnsubBroadcastHandler, service.TArgsFactory[UnsubReq]{}, nil, true, false),
		MethodExistSubscription: service.NewMethodInfo(ExistSubscriptionHandler, service.TArgsFactory[ExistSubReq]{}, service.TArgsFactory[ExistSubRes]{}, false, false),
		MethodSyncSnapshot:      service.NewMethodInfo(SyncSnapshotHandler, service.TArgsFactory[SyncSnapshotReq]{}, nil, true, false),
	})

var pubEmitterInfo = service.NewServiceInfo("PubEmitter", (*PubEmitter)(nil),
	map[string]service.MethodInfo{
		MethodPub:       service.NewMethodInfo(ServerPubHandler, service.TArgsFactory[PubReq]{}, nil, true, true),
		MethodSubFailed: service.NewMethodInfo(ServerSubFailedHandler, service.TArgsFactory[SubFailedReq]{}, nil, true, true),
		MethodExistSub:  service.NewMethodInfo(ServerExistSubHandler, service.TArgsFactory[ExistSubReq]{}, service.TArgsFactory[ExistSubRes]{}, false, true),
	})

func SubBroadcastHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(SubService).SubBroadcast(ctx, arg.(*SubReq))
}
func UnsubBroadcastHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(SubService).UnSubBroadcast(ctx, arg.(*UnsubReq))
}

func ExistSubscriptionHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(SubService).ExistSubscription(ctx, arg.(*ExistSubReq), result.(*ExistSubRes))
}

func SyncSnapshotHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(SubService).SyncSnapshot(ctx, arg.(*SyncSnapshotReq))
}

func NewServerConsumerServiceInfo() *service.ServiceInfo {
	serviceName := "Consumer"
	handlerType := (*PubEmitter)(nil)
	methods := map[string]service.MethodInfo{
		MethodPub:       service.NewMethodInfo(ServerPubHandler, service.TArgsFactory[PubReq]{}, nil, true, true),
		MethodSubFailed: service.NewMethodInfo(ServerSubFailedHandler, service.TArgsFactory[SubFailedReq]{}, nil, true, true),
		MethodExistSub:  service.NewMethodInfo(ServerExistSubHandler, service.TArgsFactory[ExistSubReq]{}, service.TArgsFactory[ExistSubRes]{}, false, true),
	}
	return service.NewServiceInfo(serviceName, handlerType, methods)
}

func ServerPubHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(PubEmitter).Pub(ctx, arg.(*PubReq))
}

func ServerSubFailedHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(PubEmitter).SubFailed(ctx, arg.(*SubFailedReq))
}

func ServerExistSubHandler(serviceImpl interface{}, ctx context.Context, arg, result interface{}) error {
	return serviceImpl.(PubEmitter).ExistSub(ctx, arg.(*ExistSubReq), result.(*ExistSubRes))
}

type directResolver struct {
	nodeList []discovery.Node
}

func NewDirectResolver(nodeList ...discovery.Node) resolver.Resolver {
	return &directResolver{nodeList: nodeList}
}

func (d directResolver) Start(listener discovery.ServiceListener) {
	listener(d.nodeList)
}

func (d directResolver) Refresh() {

}

func (d directResolver) Stop() {

}

var _ resolver.Resolver = (*directResolver)(nil)
