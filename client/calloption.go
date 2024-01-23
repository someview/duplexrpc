package client

type CallOptions struct {
	TraceId []byte
	SpanId  []byte
	AltKey  string // 负载均衡指定的key
}

type CallOption func(o *CallOptions)

func WithTrace(traceId []byte, SpanId []byte) CallOption {
	return func(o *CallOptions) {
		o.TraceId = traceId
		o.SpanId = SpanId
	}
}
