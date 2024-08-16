// 用于计算地址的resolver解析器，匹配不同协议, 不同环境下的处理
// 三种环境: k8s、docker、local, 需要使用的解析器是不一样的
// 例如grpc协议,在k8s内,需要的dns resolver是goproxy
// 而docker环境或者local环境,需要的是类似grpc的passthrough的解析器
package discovery

import (
	"fmt"
)

const (
	IMScheme    = "im"
	GrpcScheme  = "grpc"
	LocalScheme = "local"
	DefaultPort = "5002"
)

func ServiceAddress(scheme string, serviceName string) string {
	// pass through
	if scheme == "" {
		return serviceName
	}
	return fmt.Sprintf("%s:///%s", scheme, serviceName)
}
