package kuberesolver

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"rpc-oneway/pkg/resolver"

	jsoniter "github.com/json-iterator/go"
	golog "gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.golog.git"
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

type kubeResolver struct {
	TargetInfo
	context.Context
	context.CancelFunc

	client *http.Client
	notify chan []resolver.ServiceInstance
	err    error
	mu     sync.Mutex

	callCount   int64
	addr        string
	targetPort  string
	serviceName string
	namespace   string
	version     string
	logger      golog.Logger
}

// Refresh implements resolver.Resolver.
// todo根据自身逻辑去实现刷新的逻辑
func (k *kubeResolver) Refresh() {

}

// Start implements resolver.Resolver.
func (k *kubeResolver) Start(listener resolver.ServiceListener) {
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

func (p *kubeResolver) doLongPollingCall(req *ProxyRequest) []resolver.ServiceInstance {
	//等待30s proxy释放连接
	// p.callCount++
	//if p.logger.Enable(golog.LogLevelDebug) {
	//	p.logger.PrintfDebug("proxy访问次数: %d", p.callCount)
	//}
	marshal, err := jsoniter.Marshal(req)
	if err != nil {
		p.logger.PrintfError("layer:grpc,service:%s,desc:marshal req error,alert:code bug\n", p.serviceName, err)
		time.Sleep(time.Second)
		return nil
	}

	reader := bytes.NewReader(marshal)
	resp, err := p.client.Post(p.addr, "application/json", reader)
	if err != nil {
		p.logger.PrintfError("layer:grpc,service:%s,desc:post error %v,alert:service may occur exception, need check proxy and client\n", p.serviceName, err)
		time.Sleep(time.Second)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			if p.logger.Enable(golog.LogLevelDebug) {
				p.logger.PrintfDebug("layer:grpc,service:%v,desc:endpoints无变更\n", p.serviceName)
			}
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			p.logger.PrintlnWarning("layer:grpc,service:%s,desc:res status-%v,reason-%s,alert:service may occur exception, need check proxy and client\n",
				p.serviceName, resp.StatusCode, string(bodyBytes))
		}
		_ = resp.Body.Close()
		time.Sleep(time.Second)
		return nil
	}
	obj := &ProxyResponse{}
	err = jsoniter.NewDecoder(resp.Body).Decode(obj)
	_ = resp.Body.Close()
	if err != nil {
		p.logger.PrintfError("layer:grpc,service:%s,desc:unmarshal res error %v", p.serviceName, err)
		time.Sleep(time.Second)
		return nil
	}
	// https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes
	// https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions
	// 版本号不相等时，向下传递地址更新事件
	if p.logger.Enable(golog.LogLevelDebug) {
		p.logger.PrintlnDebug("layer:grpc,service:", p.serviceName, "desc:服务列表发生变化,当前版本号-",
			obj.ResourceVersion, "历史版本号-", p.version, "端点列表-", obj.Endpoints)
	}
	p.version = obj.ResourceVersion
	//更新服务端地址列表
	p.mu.Lock()
	var serviceInstances []resolver.ServiceInstance

	for _, end := range obj.Endpoints {
		serviceInstances = append(serviceInstances, resolver.ServiceInstance{ // todo 添加负载指标到map
			Endpoint: end.Url,
		})
	}
	p.mu.Unlock()
	return serviceInstances
}

// Stop implements resolver.Resolver.
func (k *kubeResolver) Stop() { // 退出循环
	k.CancelFunc()
}

// 暂时不做配置addr变更的动态监听
// 对于所有的客户端来说，共用一个logger就可以了
func NewKubeResolver(logger golog.Logger, cli *http.Client, registryUrl string, info TargetInfo) resolver.Resolver {
	if info.scheme == LocalScheme {
		//return newLocalDiscovery(info)
		// todo 匹配非k8s环境
		return nil
	}
	logger.PrintlnDebug("info: ", info.port, info.serviceName, info.serviceNamespace)
	res := &kubeResolver{
		TargetInfo: info,
		client:     &http.Client{Timeout: 40 * time.Second, Transport: _defaultTransport},
		logger:     logger,
	}
	res.Context, res.CancelFunc = context.WithCancel(context.Background())
	return res
}
