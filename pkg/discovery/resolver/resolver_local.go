package resolver

import (
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"log/slog"
	"net"
	"time"
)

type localResolver struct {
	service  string
	nodes    []discovery.Node
	ticker   *time.Ticker
	listener discovery.ServiceListener
}

// NewLocalResolver 相当于直连
//
//	services为本地srv地址
func NewLocalResolver(service string) Resolver {
	r := &localResolver{
		ticker:  time.NewTicker(time.Second * 15),
		service: service,
	}
	return r
}

func (l *localResolver) Start(listener discovery.ServiceListener) {
	l.listener = listener
	l.Refresh()
	host, port, err := parseServiceName(l.service)
	if err != nil {
		panic(fmt.Errorf("service is illegal uri: %v", l.service))
	}
	var nodes []discovery.Node
	for range l.ticker.C {
		nodes, err = l.queryNodes(host, port)
		if err != nil {
			slog.Debug(slogKey, slog.String("error", err.Error()))
			continue
		}
		listener(nodes)
	}
}

func (l *localResolver) queryNodes(host string, port string) ([]discovery.Node, error) {
	ips, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}
	slog.Debug("", slog.Any("instances", ips), slog.String("port", port))
	res := make([]discovery.Node, len(ips))
	var extend string
	if port == "" {
		extend = ""
	} else {
		extend = ":" + port
	}
	for i, ip := range ips {
		res[i] = discovery.NewNode(fmt.Sprintf("%s%s", ip, extend), nil)
	}
	return res, nil
}

func (l *localResolver) Refresh() {
	host, port, err := parseServiceName(l.service)
	if err != nil {
		panic("service is illegal uri")
	}
	nodes, err := l.queryNodes(host, port)
	if err != nil {
		return
	}
	l.listener(nodes)
}

func (l *localResolver) Stop() {
	l.ticker.Stop()
}
