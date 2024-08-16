package server

import (
	"context"
	"errors"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/metadata"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/middleware"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
)

var _ ServerXClient = (*serverXClient)(nil)

type serverXClient struct {
	targetServiceInfo *service.ServiceInfo
	transHandler      remote.TransHandler
	opt               *option
	callFunc          middleware.GenericCall
}

func NewServerXClient(targetServiceInfo *service.ServiceInfo, transHandler remote.TransHandler, opt *option) ServerXClient {
	c := &serverXClient{
		targetServiceInfo: targetServiceInfo,
		transHandler:      transHandler,
		opt:               opt,
	}
	c.callFunc = middleware.Chain(opt.mws...)(c.callCommon)
	return c
}

func (s *serverXClient) Call(ctx context.Context, methodName string, req any, resp any) (err error) {
	ri := rpcinfo.NewRPCInfo(rpcinfo.Unary, rpcinfo.NewInvocation(s.targetServiceInfo.ServiceName, methodName))
	defer func() {
		ri.Recycle()
		metadata.RecycleContext(ctx)
	}()
	return s.callFunc(ctx, req, resp, ri)
}

func (s *serverXClient) CallOneway(ctx context.Context, methodName string, req any) (err error) {
	ri := rpcinfo.NewRPCInfo(rpcinfo.Oneway, rpcinfo.NewInvocation(s.targetServiceInfo.ServiceName, methodName))
	defer func() {
		ri.Recycle()
		metadata.RecycleContext(ctx)
	}()
	return s.callFunc(ctx, req, nil, ri)
}

func (s *serverXClient) callCommon(ctx context.Context, args any, result any, ri rpcinfo.RPCInfo) (err error) {
	end, ok := metadata.ExtractEndpoint(ctx)
	if !ok {
		return errors.New("missing endpoint, use metadata.WithEndpoint")
	}
	sendMsg := remote.NewMessage(ri, s.opt.remoteOpt.Codec)
	defer sendMsg.Recycle()
	sendMsg.SetData(args)
	err = s.transHandler.Write(ctx, end, sendMsg)
	if err != nil {
		return err
	}
	if ri.InteractionMode() == rpcinfo.Oneway {
		return nil
	}
	recvMsg := remote.NewMessage(ri, s.opt.remoteOpt.Codec)
	defer recvMsg.Recycle()
	recvMsg.SetData(result)
	return s.transHandler.Read(ctx, end, recvMsg)
}
