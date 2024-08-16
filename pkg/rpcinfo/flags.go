package rpcinfo

var (
	_ FlagInfo       = (*flag)(nil)
	_ FlagInfoSetter = (*flag)(nil)
)

type flag struct {
	traceInfo []byte
}

func (f *flag) SetTraceInfo(traceInfo []byte) {
	f.traceInfo = traceInfo[:TraceLength]
}

func (f *flag) TraceInfo() []byte {
	return f.traceInfo
}

func (f *flag) zero() {
	f.traceInfo = nil
}
