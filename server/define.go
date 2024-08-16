package server

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"
)

type ServerXClient interface {
	Call(ctx context.Context, methodName string, req any, resp any) (err error)
	CallOneway(ctx context.Context, methodName string, req any) (err error)
}

type Server interface {
	RegisterService(svcInfo *service.ServiceInfo, handler interface{}) error
	GetServiceInfos() map[string]*service.ServiceInfo
	NewServerXClient(svcInfo *service.ServiceInfo) ServerXClient // 创建客户端向服务端通讯的请求的client
	Run() error
	Stop() error
}
