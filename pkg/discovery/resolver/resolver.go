package resolver

import (
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"net/http"
)

func NewResolver(cli *http.Client, registryUrl string) Resolver {
	info, _ := parseServiceAddr(registryUrl)
	// todo 匹配非k8s环境
	if info.scheme == string(discovery.LocalScheme) {
		return NewLocalResolver(info.serviceName + ":" + info.port)
	} else {
		return NewKubeResolver(cli, registryUrl, info)
	}
}
