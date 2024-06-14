package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"io"
	"net/http"
)

func main() {
	s := &service{endpoint: make(map[string]discovery.ServiceInstance)}
	http.HandleFunc("/im", s.route)
	panic(http.ListenAndServe(":5004", nil))
}

type service struct {
	endpoint map[string]discovery.ServiceInstance
}

func (s *service) route(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("新请求", req.Method)
	switch req.Method {
	case http.MethodGet:
		s.handleGet(resp, req)
	case http.MethodPost:
		s.handlePost(resp, req)
	case http.MethodDelete:
		s.handleDelete(resp, req)
	}
}

func (s *service) handleGet(resp http.ResponseWriter, req *http.Request) {
	var err error
	var p []byte
	defer func() {
		if err != nil {
			b, _ := sonic.Marshal(err.Error())
			resp.Write(b)
			return
		} else {
			resp.Write(p)
			return
		}
	}()
	nodes := make([]discovery.ServiceInstance, len(s.endpoint))
	for _, v := range s.endpoint {
		nodes = append(nodes, v)
	}
	p, err = sonic.Marshal(nodes)
}

func (s *service) handlePost(resp http.ResponseWriter, req *http.Request) {
	var err error
	var p []byte
	defer func() {
		if err != nil {
			b, _ := sonic.Marshal(err.Error())
			resp.Write(b)
			return
		} else {
			resp.Write(p)
			return
		}
	}()
	p, err = io.ReadAll(req.Body)
	if err != nil {
		return
	}
	var result discovery.ServiceInstance
	err = sonic.Unmarshal(p, &result)

	if err != nil {
		return
	}
	fmt.Println("服务注册：", result.Address())
	s.endpoint[result.Address()] = result
}

func (s *service) handleDelete(resp http.ResponseWriter, req *http.Request) {
	var err error
	var p []byte
	defer func() {
		if err != nil {
			b, _ := sonic.Marshal(err.Error())
			resp.Write(b)
			return
		} else {
			resp.Write(p)
			return
		}
	}()
	p, err = io.ReadAll(req.Body)
	if err != nil {
		return
	}
	var result discovery.ServiceInstance
	err = sonic.Unmarshal(p, &result)
	if err != nil {
		return
	}
	fmt.Println("服务注销：", result.Address())
	delete(s.endpoint, result.Address())
}
