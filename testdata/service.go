package testdata

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/server"
)

type uServerClient struct {
	server.ServerXClient
}

func (u *uServerClient) SelfServiceType() SubService {
	return (SubService)(nil)
}

func (u *uServerClient) Pub(ctx context.Context, req *PubReq) error {
	return u.ServerXClient.CallOneway(ctx, MethodPub, req)
}

func (u *uServerClient) SubFailed(ctx context.Context, req *SubFailedReq) error {
	return u.ServerXClient.CallOneway(ctx, MethodSubFailed, req)
}

func (u *uServerClient) ExistSub(ctx context.Context, req *ExistSubReq, res *ExistSubRes) error {
	return u.ServerXClient.Call(ctx, MethodExistSub, req, res)
}

var _ ConsumerServiceServer = (*uServerClient)(nil)

func NewConsumerServer(srv server.Server, selfSvcImpl any) ConsumerServiceServer {
	if err := srv.RegisterService(subServiceInfo, selfSvcImpl); err != nil {
		panic(err)
	}
	return &uServerClient{srv.NewServerXClient(pubEmitterInfo)}
}
