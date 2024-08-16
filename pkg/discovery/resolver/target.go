package resolver

import (
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type TargetInfo struct {
	serviceName      string
	serviceNamespace string
	port             string
	scheme           string
}

// const IMOnewayDeafultServicePort = "5004"
// const GrpcScheme = "grpc"

const defaultNameSpace = "127.0.0.1"

const ServiceAccountNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

var currentNamespace = LoadNamespace()

// LoadNamespace is used to get the current namespace from the file
func LoadNamespace() string {
	data, err := os.ReadFile(ServiceAccountNamespacePath)
	if err != nil {
		return defaultNameSpace
	}
	return string(data)
}

func parseServiceAddr(service string) (res TargetInfo, err error) {
	parse, err := url.Parse(service)
	if err != nil {
		return res, fmt.Errorf("parse service name error: %v", err)
	}
	// parse必定不能为空,path等同于微服务的serviceName
	if parse.Path == "" {
		service = fmt.Sprintf("%s:///%s", discovery.LocalScheme, service)
		return parseServiceAddr(service)
	}
	// host当做namespace来进行处理
	if parse.Host == "" {
		res.serviceNamespace = currentNamespace
	} else {
		res.serviceNamespace = parse.Host
	}

	temp := strings.Split(strings.TrimPrefix(parse.Path, "/"), ":")
	res.serviceName = temp[0]
	if len(temp) > 1 {
		res.port = temp[1]
	}
	res.scheme = parse.Scheme
	if !validPort(res.port) {
		return res, fmt.Errorf("service port must be empty or integer str")
	}
	return res, nil
}

// serviceName等于于一个url
func parseServiceName(serviceName string) (path string, port string, err error) {
	// 尝试解析为URL
	info, err := parseServiceAddr(serviceName)
	if err != nil {
		return
	}
	return info.serviceName, info.port, nil
}

func validPort(port string) bool {
	if port == "" {
		return true
	}
	_, err := strconv.ParseInt(port, 10, 64)
	return err == nil
}
