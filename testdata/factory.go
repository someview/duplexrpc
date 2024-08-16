package testdata

import "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/service"

func SetSubServiceArgGetter(methodName string, req, res service.ArgFactory) {
	oldInfo := subServiceInfo.MethodInfo(methodName)

	newInfo := service.NewMethodInfo(oldInfo.Handler(), req, res, oldInfo.OneWay(), oldInfo.ServerCall())

	subServiceInfo.Methods[methodName] = newInfo

}
