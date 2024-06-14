// 用于计算地址的resolver解析器，匹配不同协议, 不同环境下的处理
// 三种环境: k8s、docker、local, 需要使用的解析器是不一样的
// 例如grpc协议,在k8s内,需要的dns resolver是goproxy
// 而docker环境或者local环境,需要的是类似grpc的passthrough的解析器
package discovery

import (
	"fmt"
	"strings"

	"gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.goconfig.git/setting"
)

type Scheme string

const (
	IMScheme    Scheme = "im"
	GrpcScheme  Scheme = "grpc"
	LocalScheme Scheme = "local"
	IMPort             = 5004
	GrpcPort           = 5002
	DefaultPort        = 5004
)

// Scheme用于表示域名解析器的分类
func (s Scheme) IsValid() bool {
	if s == IMScheme || s == GrpcScheme {
		return true
	}
	if s == "" || s == LocalScheme {
		return true
	}
	return false
}

func (s Scheme) Port() int {
	if s == IMScheme {
		return IMPort
	}
	if s == GrpcScheme {
		return GrpcPort
	}
	return DefaultPort
}

// getAddressInDocker
// 在容器内时, 根据服务名获取对方地址, 通过k8s提供的DNS功能, 加上默认的端口.
func ContainerAddress(port int, serviceName, productCode, scheme string) string { // 看k8s ip中不包含host
	ip := getDomainForDocker(serviceName, productCode)
	return fmt.Sprintf("%s:///%s:%d", scheme, ip, port)
}

func PassThroughAddress(port int, serviceName, scheme string) string {
	return fmt.Sprintf("%s:///%s:%d", scheme, serviceName, port)
}

// 获取docker中的服务
func getDomainForDocker(serviceName, productCode string) string {
	return strings.ToLower(productCode + "-" + serviceName)
}

// Scheme用于表示域名解析的分类
func ServiceAddress(
	env setting.RuntimeEnvironment, productCode string,
	scheme Scheme, serviceName string, port int) string { // 看k8s ip中不包含host
	if scheme == "" {
		scheme = LocalScheme
	}
	if !scheme.IsValid() {
		panic("scheme is not allowed")
	}

	// 容器环境下,端口采用协议默认的端口
	if env == setting.K8s || env == setting.Docker {
		port = scheme.Port()
		return ContainerAddress(port, serviceName, productCode, string(scheme))
	}
	// 否则的话，采用passthrough的resolver
	return PassThroughAddress(port, serviceName, string(scheme))
}

// 服务端监听的端口
func ServerPort(env setting.RuntimeEnvironment, scheme Scheme, port int) int {
	if scheme == "" {
		scheme = LocalScheme
	}
	if env == setting.K8s || env == setting.Docker {
		return scheme.Port()
	}
	return port
}
