package testdata

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/client"
)

type uClient struct {
	cli client.XClient
}

func (c *uClient) SelfServiceType() PubEmitter {
	return (PubEmitter)(nil)
}

func (c *uClient) ServiceType() PubEmitter {
	return (PubEmitter)(nil)
}

func (c *uClient) SubBroadcast(ctx context.Context, req *SubReq) error {
	return c.cli.CallOneway(ctx, MethodSubBroadcast, req)
}

func (c *uClient) UnSubBroadcast(ctx context.Context, req *UnsubReq) error {
	return c.cli.CallOneway(ctx, MethodUnSubBroadcast, req)
}

func (c *uClient) ExistSubscription(ctx context.Context, req *ExistSubReq, res *ExistSubRes) error {
	return c.cli.Call(ctx, MethodExistSubscription, req, res)
}

func (c *uClient) SyncSnapshot(ctx context.Context, req *SyncSnapshotReq) error {
	return c.cli.CallOneway(ctx, MethodSyncSnapshot, req)
}

var _ ConsumerServiceClient = (*uClient)(nil)

func NewConsumerServiceClient(dest string, selfSvcImpl any, opt ...client.Option) ConsumerServiceClient {
	cli, err := client.NewDuplexClient(dest, subServiceInfo, pubEmitterInfo, selfSvcImpl, opt...)
	if err != nil {
		panic(err)
	}
	return &uClient{cli: cli}
}
