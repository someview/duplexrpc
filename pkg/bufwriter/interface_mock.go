// Code generated by MockGen. DO NOT EDIT.
// Source: define.go
//
// Generated by this command:
//
//	mockgen -source=define.go -destination ./interface_mock.go -package bufwriter
//

// Package bufwriter is a generated GoMock package.
package bufwriter

import (
	context "context"
	net "net"
	reflect "reflect"
	time "time"

	netpoll "github.com/cloudwego/netpoll"
	gomock "go.uber.org/mock/gomock"
)

// MockBufWriter is a mock of BufWriter interface.
type MockBufWriter struct {
	ctrl     *gomock.Controller
	recorder *MockBufWriterMockRecorder
}

// MockBufWriterMockRecorder is the mock recorder for MockBufWriter.
type MockBufWriterMockRecorder struct {
	mock *MockBufWriter
}

// NewMockBufWriter creates a new mock instance.
func NewMockBufWriter(ctrl *gomock.Controller) *MockBufWriter {
	mock := &MockBufWriter{ctrl: ctrl}
	mock.recorder = &MockBufWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBufWriter) EXPECT() *MockBufWriterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockBufWriter) Add(ctx context.Context, lb *netpoll.LinkBuffer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, lb)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockBufWriterMockRecorder) Add(ctx, lb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockBufWriter)(nil).Add), ctx, lb)
}

// Close mocks base method.
func (m *MockBufWriter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBufWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBufWriter)(nil).Close))
}

// MockBufFlusher is a mock of BufFlusher interface.
type MockBufFlusher struct {
	ctrl     *gomock.Controller
	recorder *MockBufFlusherMockRecorder
}

// MockBufFlusherMockRecorder is the mock recorder for MockBufFlusher.
type MockBufFlusherMockRecorder struct {
	mock *MockBufFlusher
}

// NewMockBufFlusher creates a new mock instance.
func NewMockBufFlusher(ctrl *gomock.Controller) *MockBufFlusher {
	mock := &MockBufFlusher{ctrl: ctrl}
	mock.recorder = &MockBufFlusherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBufFlusher) EXPECT() *MockBufFlusherMockRecorder {
	return m.recorder
}

// FlushTo mocks base method.
func (m *MockBufFlusher) FlushTo(arg0 netpoll.Connection) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushTo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FlushTo indicates an expected call of FlushTo.
func (mr *MockBufFlusherMockRecorder) FlushTo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushTo", reflect.TypeOf((*MockBufFlusher)(nil).FlushTo), arg0)
}

// MockFlowControl is a mock of FlowControl interface.
type MockFlowControl struct {
	ctrl     *gomock.Controller
	recorder *MockFlowControlMockRecorder
}

// MockFlowControlMockRecorder is the mock recorder for MockFlowControl.
type MockFlowControlMockRecorder struct {
	mock *MockFlowControl
}

// NewMockFlowControl creates a new mock instance.
func NewMockFlowControl(ctrl *gomock.Controller) *MockFlowControl {
	mock := &MockFlowControl{ctrl: ctrl}
	mock.recorder = &MockFlowControlMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlowControl) EXPECT() *MockFlowControlMockRecorder {
	return m.recorder
}

// Available mocks base method.
func (m *MockFlowControl) Available() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Available")
	ret0, _ := ret[0].(int)
	return ret0
}

// Available indicates an expected call of Available.
func (mr *MockFlowControlMockRecorder) Available() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Available", reflect.TypeOf((*MockFlowControl)(nil).Available))
}

// GetWithCtx mocks base method.
func (m *MockFlowControl) GetWithCtx(ctx context.Context, sz int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithCtx", ctx, sz)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetWithCtx indicates an expected call of GetWithCtx.
func (mr *MockFlowControlMockRecorder) GetWithCtx(ctx, sz any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithCtx", reflect.TypeOf((*MockFlowControl)(nil).GetWithCtx), ctx, sz)
}

// Release mocks base method.
func (m *MockFlowControl) Release(sz int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Release", sz)
}

// Release indicates an expected call of Release.
func (mr *MockFlowControlMockRecorder) Release(sz any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockFlowControl)(nil).Release), sz)
}

// TryGet mocks base method.
func (m *MockFlowControl) TryGet(sz int) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryGet", sz)
	ret0, _ := ret[0].(bool)
	return ret0
}

// TryGet indicates an expected call of TryGet.
func (mr *MockFlowControlMockRecorder) TryGet(sz any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryGet", reflect.TypeOf((*MockFlowControl)(nil).TryGet), sz)
}

// MockBatchContainer is a mock of BatchContainer interface.
type MockBatchContainer struct {
	ctrl     *gomock.Controller
	recorder *MockBatchContainerMockRecorder
}

// MockBatchContainerMockRecorder is the mock recorder for MockBatchContainer.
type MockBatchContainerMockRecorder struct {
	mock *MockBatchContainer
}

// NewMockBatchContainer creates a new mock instance.
func NewMockBatchContainer(ctrl *gomock.Controller) *MockBatchContainer {
	mock := &MockBatchContainer{ctrl: ctrl}
	mock.recorder = &MockBatchContainerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchContainer) EXPECT() *MockBatchContainerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockBatchContainer) Add(ctx context.Context, lb *netpoll.LinkBuffer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, lb)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockBatchContainerMockRecorder) Add(ctx, lb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockBatchContainer)(nil).Add), ctx, lb)
}

// Close mocks base method.
func (m *MockBatchContainer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBatchContainerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBatchContainer)(nil).Close))
}

// FlushTo mocks base method.
func (m *MockBatchContainer) FlushTo(arg0 netpoll.Connection) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushTo", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FlushTo indicates an expected call of FlushTo.
func (mr *MockBatchContainerMockRecorder) FlushTo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushTo", reflect.TypeOf((*MockBatchContainer)(nil).FlushTo), arg0)
}

// IsEmpty mocks base method.
func (m *MockBatchContainer) IsEmpty() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty.
func (mr *MockBatchContainerMockRecorder) IsEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockBatchContainer)(nil).IsEmpty))
}

// IsFull mocks base method.
func (m *MockBatchContainer) IsFull() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFull")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFull indicates an expected call of IsFull.
func (mr *MockBatchContainerMockRecorder) IsFull() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFull", reflect.TypeOf((*MockBatchContainer)(nil).IsFull))
}

// Len mocks base method.
func (m *MockBatchContainer) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockBatchContainerMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockBatchContainer)(nil).Len))
}

// TryAdd mocks base method.
func (m *MockBatchContainer) TryAdd(lb *netpoll.LinkBuffer) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryAdd", lb)
	ret0, _ := ret[0].(bool)
	return ret0
}

// TryAdd indicates an expected call of TryAdd.
func (mr *MockBatchContainerMockRecorder) TryAdd(lb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryAdd", reflect.TypeOf((*MockBatchContainer)(nil).TryAdd), lb)
}

// MockBatchWriter is a mock of BatchWriter interface.
type MockBatchWriter struct {
	ctrl     *gomock.Controller
	recorder *MockBatchWriterMockRecorder
}

// MockBatchWriterMockRecorder is the mock recorder for MockBatchWriter.
type MockBatchWriterMockRecorder struct {
	mock *MockBatchWriter
}

// NewMockBatchWriter creates a new mock instance.
func NewMockBatchWriter(ctrl *gomock.Controller) *MockBatchWriter {
	mock := &MockBatchWriter{ctrl: ctrl}
	mock.recorder = &MockBatchWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchWriter) EXPECT() *MockBatchWriterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockBatchWriter) Add(ctx context.Context, lb *netpoll.LinkBuffer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, lb)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockBatchWriterMockRecorder) Add(ctx, lb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockBatchWriter)(nil).Add), ctx, lb)
}

// Close mocks base method.
func (m *MockBatchWriter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBatchWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBatchWriter)(nil).Close))
}

// MockConnection is a mock of Connection interface.
type MockConnection struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionMockRecorder
}

// MockConnectionMockRecorder is the mock recorder for MockConnection.
type MockConnectionMockRecorder struct {
	mock *MockConnection
}

// NewMockConnection creates a new mock instance.
func NewMockConnection(ctrl *gomock.Controller) *MockConnection {
	mock := &MockConnection{ctrl: ctrl}
	mock.recorder = &MockConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnection) EXPECT() *MockConnectionMockRecorder {
	return m.recorder
}

// AddCloseCallback mocks base method.
func (m *MockConnection) AddCloseCallback(callback netpoll.CloseCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCloseCallback", callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCloseCallback indicates an expected call of AddCloseCallback.
func (mr *MockConnectionMockRecorder) AddCloseCallback(callback any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCloseCallback", reflect.TypeOf((*MockConnection)(nil).AddCloseCallback), callback)
}

// Close mocks base method.
func (m *MockConnection) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnectionMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnection)(nil).Close))
}

// IsActive mocks base method.
func (m *MockConnection) IsActive() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsActive")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsActive indicates an expected call of IsActive.
func (mr *MockConnectionMockRecorder) IsActive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsActive", reflect.TypeOf((*MockConnection)(nil).IsActive))
}

// LocalAddr mocks base method.
func (m *MockConnection) LocalAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// LocalAddr indicates an expected call of LocalAddr.
func (mr *MockConnectionMockRecorder) LocalAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalAddr", reflect.TypeOf((*MockConnection)(nil).LocalAddr))
}

// Read mocks base method.
func (m *MockConnection) Read(b []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockConnectionMockRecorder) Read(b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockConnection)(nil).Read), b)
}

// Reader mocks base method.
func (m *MockConnection) Reader() netpoll.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader")
	ret0, _ := ret[0].(netpoll.Reader)
	return ret0
}

// Reader indicates an expected call of Reader.
func (mr *MockConnectionMockRecorder) Reader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reader", reflect.TypeOf((*MockConnection)(nil).Reader))
}

// RemoteAddr mocks base method.
func (m *MockConnection) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockConnectionMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockConnection)(nil).RemoteAddr))
}

// SetDeadline mocks base method.
func (m *MockConnection) SetDeadline(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDeadline", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDeadline indicates an expected call of SetDeadline.
func (mr *MockConnectionMockRecorder) SetDeadline(t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*MockConnection)(nil).SetDeadline), t)
}

// SetIdleTimeout mocks base method.
func (m *MockConnection) SetIdleTimeout(timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetIdleTimeout", timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetIdleTimeout indicates an expected call of SetIdleTimeout.
func (mr *MockConnectionMockRecorder) SetIdleTimeout(timeout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIdleTimeout", reflect.TypeOf((*MockConnection)(nil).SetIdleTimeout), timeout)
}

// SetOnRequest mocks base method.
func (m *MockConnection) SetOnRequest(on netpoll.OnRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOnRequest", on)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOnRequest indicates an expected call of SetOnRequest.
func (mr *MockConnectionMockRecorder) SetOnRequest(on any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnRequest", reflect.TypeOf((*MockConnection)(nil).SetOnRequest), on)
}

// SetReadDeadline mocks base method.
func (m *MockConnection) SetReadDeadline(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline.
func (mr *MockConnectionMockRecorder) SetReadDeadline(t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockConnection)(nil).SetReadDeadline), t)
}

// SetReadTimeout mocks base method.
func (m *MockConnection) SetReadTimeout(timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadTimeout", timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadTimeout indicates an expected call of SetReadTimeout.
func (mr *MockConnectionMockRecorder) SetReadTimeout(timeout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadTimeout", reflect.TypeOf((*MockConnection)(nil).SetReadTimeout), timeout)
}

// SetWriteDeadline mocks base method.
func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline.
func (mr *MockConnectionMockRecorder) SetWriteDeadline(t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*MockConnection)(nil).SetWriteDeadline), t)
}

// SetWriteTimeout mocks base method.
func (m *MockConnection) SetWriteTimeout(timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteTimeout", timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteTimeout indicates an expected call of SetWriteTimeout.
func (mr *MockConnectionMockRecorder) SetWriteTimeout(timeout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteTimeout", reflect.TypeOf((*MockConnection)(nil).SetWriteTimeout), timeout)
}

// SliceIntoReader mocks base method.
func (m *MockConnection) SliceIntoReader(n int, r netpoll.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SliceIntoReader", n, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// SliceIntoReader indicates an expected call of SliceIntoReader.
func (mr *MockConnectionMockRecorder) SliceIntoReader(n, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SliceIntoReader", reflect.TypeOf((*MockConnection)(nil).SliceIntoReader), n, r)
}

// Write mocks base method.
func (m *MockConnection) Write(b []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockConnectionMockRecorder) Write(b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockConnection)(nil).Write), b)
}

// Writer mocks base method.
func (m *MockConnection) Writer() netpoll.Writer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer")
	ret0, _ := ret[0].(netpoll.Writer)
	return ret0
}

// Writer indicates an expected call of Writer.
func (mr *MockConnectionMockRecorder) Writer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockConnection)(nil).Writer))
}
