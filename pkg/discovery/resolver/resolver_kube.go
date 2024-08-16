package resolver

import (
	"bytes"
	"context"
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	//golog "gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.golog.git"
)

var _dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

var _defaultTransport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	DialContext:           _dialer.DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	MaxIdleConnsPerHost:   100,
}

const slogKey = "layer:rpc"

type kubeResolver struct {
	TargetInfo
	context.Context
	context.CancelFunc

	client *http.Client
	notify chan []discovery.ServiceInstance
	err    error
	mu     sync.Mutex

	callCount   int64
	addr        string
	targetPort  string
	serviceName string
	namespace   string
	version     string

	logger      *slog.Logger
	serviceAttr slog.Attr
}

// Refresh implements resolver.Resolver.
// todo根据自身逻辑去实现刷新的逻辑
func (k *kubeResolver) Refresh() {

}

// Start implements resolver.Resolver.
func (k *kubeResolver) Start(listener discovery.ServiceListener) {
	req := &ProxyRequest{
		Service:         k.serviceName,
		Namespace:       k.namespace,
		ResourceVersion: k.version,
		PortName:        k.scheme,
	}
	for {
		select {
		case <-k.Context.Done():
			return
		default:
			req.ResourceVersion = k.version
			res := k.doLongPollingCall(req) // doCall应该控制访问次数，否则会有风暴问题
			if res != nil {
				listener(res)
			}
		}
	}
}

func (k *kubeResolver) doLongPollingCall(req *ProxyRequest) []discovery.Node {

	marshal, err := jsoniter.Marshal(req)
	if err != nil {
		k.logger.Error(slogKey, k.serviceAttr, slog.String("error", err.Error()),
			slog.String("alert", "service may occur exception, need check proxy and client"))
		time.Sleep(time.Second)
		return nil
	}

	reader := bytes.NewReader(marshal)
	resp, err := k.client.Post(k.addr, "application/json", reader)
	if err != nil {
		k.logger.Error(slogKey, k.serviceAttr, slog.String("error", err.Error()),
			slog.String("alert", "service may occur exception, need check proxy and client"))
		time.Sleep(time.Second)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			k.logger.Debug(slogKey, slog.String("service", k.serviceName), slog.String("desc", "endpoints无变更"))
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			k.logger.Warn(slogKey, k.serviceAttr,
				slog.String("desc", fmt.Sprintf("status-%v,reason-%s", resp.StatusCode, string(bodyBytes))),
				slog.String("alert", "service may occur exception, need check proxy and client"))
		}
		_ = resp.Body.Close()
		time.Sleep(time.Second)
		return nil
	}
	obj := &ProxyResponse{}
	err = jsoniter.NewDecoder(resp.Body).Decode(obj)
	_ = resp.Body.Close()
	if err != nil {
		k.logger.Error(slogKey,
			slog.String("service", k.serviceName), slog.String("error", err.Error()),
			slog.String("alert", "service may occur exception, need check proxy and client"),
		)
		time.Sleep(time.Second)
		return nil
	}
	// https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes
	// https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions
	// 版本号不相等时，向下传递地址更新事件
	// todo 实现一个规范化的slog的接口
	k.logger.Debug(slogKey, slog.String("service", k.serviceName),
		slog.String("当前版本号", obj.ResourceVersion),
		slog.String("历史版本号", k.version), slog.Any("端点列表", obj.Endpoints))

	k.version = obj.ResourceVersion
	//更新服务端地址列表
	k.mu.Lock()
	var serviceInstances []discovery.Node

	for _, end := range obj.Endpoints {
		serviceInstances = append(serviceInstances, discovery.NewNode(end.Url, end.Metadata))
	}
	k.mu.Unlock()
	return serviceInstances
}

// Stop implements resolver.Resolver.
func (k *kubeResolver) Stop() { // 退出循环
	k.CancelFunc()
}

// 暂时不做配置addr变更的动态监听
// 对于所有的客户端来说，共用一个logger就可以了
func NewKubeResolver(cli *http.Client, registryUrl string, info TargetInfo) Resolver {

	// logger.PrintlnDebug("info: ", info.port, info.serviceName, info.serviceNamespace)
	res := &kubeResolver{
		TargetInfo:  info,
		client:      &http.Client{Timeout: 40 * time.Second, Transport: _defaultTransport},
		logger:      slog.Default(),
		serviceAttr: slog.String("service", info.serviceName),
	}

	res.Context, res.CancelFunc = context.WithCancel(context.Background())
	return res
}
