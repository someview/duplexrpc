package client

import (
	"context"
	"fmt"
	"log"
	"rpc-oneway/protocol"
	"rpc-oneway/server"
	"rpc-oneway/testdata"
	"testing"
	"time"
)

type msgMock struct{}

// MarshalToSizedBuffer implements protocol.SizeableMarshaller.
func (m *msgMock) MarshalToSizedBuffer([]byte) (int, error) {
	return 0, nil
}

// Size implements protocol.SizeableMarshaller.
func (m *msgMock) Size() int {
	return 10
}

// Unmarshal implements protocol.SizeableMarshaller.
func (*msgMock) Unmarshal(data []byte) error {
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}
	return nil
}

var _ protocol.SizeableMarshaller = (*msgMock)(nil)

func TestClient(t *testing.T) {
	srv := server.NewServer()
	// todo 将添加到的消息添加到handler里面, 这里为了避免使用反射，直接使用断言，这样效率会高得多
	srv.AddHandler(1, func(ctx *server.ClientRequestContext, msg any) error {
		fmt.Println("进来,收到消息:", msg)
		return nil
	})
	go func() {
		if err := srv.Serve("tcp", ":8080"); err != nil {
			log.Fatalln("err:", err)
		}
	}()
	//defer func(srv *server.Server, ctx context.Context) {
	//	err := srv.Shutdown(ctx)
	//	if err != nil {
	//		log.Println("err:", err)
	//	}
	//}(srv, context.TODO())

	time.Sleep(time.Second * 1)
	cli := &MuxClient{}
	if err := cli.Connect("tcp", "127.0.0.1:8080"); err != nil {
		log.Fatalln("err:", err)
	}
	if err := cli.Send(context.Background(), 1, &testdata.ClientMessage{
		Header: &testdata.Header{TraceId: "123456789"},
	}); err != nil {
		log.Fatalln("err:", err)
	}

	time.Sleep(time.Minute * 2)
}
