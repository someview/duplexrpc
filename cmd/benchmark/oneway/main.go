package main

import (
	"context"
	"fmt"
	"github.com/jhue58/latency"
	"github.com/jhue58/latency/buckets"
	"github.com/panjf2000/ants/v2"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/client"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/remote"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/rpcinfo"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/server"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testdata"
	"sync"
	"time"
)

var subReqPool = sync.Pool{New: func() any { return new(testdata.SubReq) }}

type subReqFactory struct {
}

func (s subReqFactory) New() any {
	return subReqPool.Get()
}

func (s subReqFactory) Recycle(a any) {
	subReqPool.Put(a)
}

func main() {
	testdata.SetSubServiceArgGetter(testdata.MethodSubBroadcast, subReqFactory{}, nil)
	addr := "127.0.0.1:8085"
	count := 1000000
	wg := sync.WaitGroup{}
	svr := server.NewServer(server.WithAddress(addr), server.WithRemoteOption(remote.WithOnRPCDone(func(ctx context.Context, req, res any, ri rpcinfo.RPCInfo) {
		wg.Done()
	})))
	_ = testdata.NewConsumerServer(svr, new(mockSubService))
	go func() {
		err := svr.Run()
		if err != nil {
			panic(err)
		}
	}()

	cli := testdata.NewConsumerServiceClient("consumerService:5002", new(mockPubListener), client.WithResolver(testdata.NewDirectResolver(
		discovery.NewNode(addr, nil))))

	time.Sleep(3 * time.Second)

	reporter := latency.NewLateReporter(buckets.NewBucketsRecorder())
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			fmt.Println(reporter.Report())
		}
	}()
	pool, _ := ants.NewPool(100)

	start := time.Now()
	defer func() {
		fmt.Printf("%d条Oneway,用时:%s\n", count, time.Since(start).String())
		fmt.Println(reporter.Report())
	}()

	task := func() {
		start, end := reporter.Alloc()
		start()
		req := subReqPool.Get().(*testdata.SubReq)
		defer func() {
			subReqPool.Put(req)
		}()
		err := cli.SubBroadcast(context.TODO(), req)
		if err != nil {
			panic(err)
		}
		end()
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		pool.Submit(task)
		//task()
	}
	wg.Wait()

}
