package registrar

import (
	"bytes"
	"context"
	"github.com/bytedance/sonic"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"io"
	"net/http"
)

type LocalRegistrar struct {
	url string
}

func NewLocalRegistrar(addr string) Registrar {
	return &LocalRegistrar{url: addr + "/im"}
}

func (l *LocalRegistrar) Register(ctx context.Context, service discovery.ServiceInstance) error {
	b, err := sonic.Marshal(service)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)
	resp, err := http.Post(l.url, "application/json", reader)
	if err != nil {
		return err
	}
	p, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result discovery.ServiceInstance
	err = sonic.Unmarshal(p, &result)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalRegistrar) Deregister(ctx context.Context, service discovery.ServiceInstance) error {
	b, err := sonic.Marshal(service)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodDelete, l.url, reader)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	p, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result discovery.ServiceInstance
	err = sonic.Unmarshal(p, &result)
	if err != nil {
		return err
	}
	return nil
}
