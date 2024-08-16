package main

import (
	"context"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/testdata"
)

type mockSubService struct {
	existSubRes testdata.ExistSubRes
	subReq      testdata.SubReq
}

func (c *mockSubService) SetExistSubscriptionRes(res *testdata.ExistSubRes) {
	if res != nil {
		c.existSubRes = *res
	}
}

func (c *mockSubService) GetSubReq() *testdata.SubReq {
	return &c.subReq
}

func (c *mockSubService) SubBroadcast(ctx context.Context, req *testdata.SubReq) error {
	c.subReq.ConnId = req.ConnId
	c.subReq.SubscriptionId = req.SubscriptionId
	return nil
}

func (c *mockSubService) UnSubBroadcast(ctx context.Context, req *testdata.UnsubReq) error {
	return nil
}

func (c *mockSubService) ExistSubscription(ctx context.Context, req *testdata.ExistSubReq, res *testdata.ExistSubRes) error {
	res.Exists = c.existSubRes.Exists
	return nil
}

func (c *mockSubService) SyncSnapshot(ctx context.Context, req *testdata.SyncSnapshotReq) error {
	return nil
}

var _ testdata.SubService = (*mockSubService)(nil)

type mockPubListener struct {
	testdata.PubReq
	ServerExistSubReq *testdata.ExistSubReq
	ServerExistSubRes *testdata.ExistSubRes
}

func (m *mockPubListener) ExistSub(ctx context.Context, req *testdata.ExistSubReq, res *testdata.ExistSubRes) error {
	m.ServerExistSubReq = req
	res.Exists = m.ServerExistSubRes.Exists
	return nil
}

func (m *mockPubListener) GetServerPubReq() *testdata.PubReq {
	return &m.PubReq
}
func (m *mockPubListener) GetServerExistSubReq() *testdata.ExistSubReq {
	return m.ServerExistSubReq
}

func (m *mockPubListener) SetServerExistSubRes(res *testdata.ExistSubRes) {
	if res != nil {
		m.ServerExistSubRes = res
	}
}

func (m *mockPubListener) Pub(ctx context.Context, req *testdata.PubReq) error {
	m.PubReq = *req
	return nil
}

func (m *mockPubListener) SubFailed(ctx context.Context, req *testdata.SubFailedReq) error {
	return nil
}

var _ testdata.PubEmitter = (*mockPubListener)(nil)
