package share

type MsgContext interface {
	TraceInfo() []byte
	SetTraceInfo([]byte)
	MsgType() byte
	SetMsgType(byte)
	Metadata() map[string]string
}

type msgContext struct {
}
