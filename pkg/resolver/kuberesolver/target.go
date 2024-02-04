package kuberesolver

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type TargetInfo struct {
	serviceName      string
	serviceNamespace string
	port             string
	scheme           string
}

const GrpcScheme = "grpc"
const LocalScheme = "local"

const defaultServicePort = "5002"
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

func parseServiceAddr(fullServiceName string) (res TargetInfo, err error) {
	parse, err := url.Parse(fullServiceName)
	if err != nil {
		return res, fmt.Errorf("parse service name error: %v", err)
	}
	if !(parse.Scheme == GrpcScheme || parse.Scheme == LocalScheme) {
		return res, fmt.Errorf("parse service schema error: %v", err)
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
	} else {
		res.port = defaultServicePort
	}
	res.scheme = parse.Scheme
	return res, nil
}
