package client

import (
	"crypto/tls"
	"time"

	"rpc-oneway/pkg/resolver"
	"rpc-oneway/pkg/selector"
	"rpc-oneway/protocol"
)

// Option contains all options for creating clients.
type Option struct {
	// Group is used to select the services in the same group. Services set group info in their meta.
	// If it is empty, clients will ignore group.
	Group string

	// Retries retries to send
	Retries int
	// Time to disallow the bad server not to be selected
	TimeToDisallow time.Duration

	// TLSConfig for tcp and quic
	TLSConfig *tls.Config
	// kcp.BlockCrypt
	Block interface{}

	// ConnectTimeout sets timeout for dialing
	ConnectTimeout time.Duration
	// IdleTimeout sets max idle time for underlying net.Conns
	IdleTimeout time.Duration

	// BackupLatency is used for Failbackup mode. rpcx will sends another request if the first response doesn't return in BackupLatency time.
	BackupLatency time.Duration

	// Breaker is used to config CircuitBreaker
	// GenBreaker func() Breaker

	// send heartbeat message to service and check responses
	Heartbeat bool
	// interval for heartbeat
	HeartbeatInterval   time.Duration
	MaxWaitForHeartbeat time.Duration

	// TCPKeepAlive, if it is zero we don't set keepalive
	TCPKeepAlivePeriod time.Duration
	// bidirectional mode, if true serverMessageChan will block to wait message for consume. default false.
	BidirectionalBlock bool

	// alaways use the selected server until it is bad
	Sticky bool

	selector.Selector
	resolver.Resolver

	serverMsgHandler func(protocol.Message)
}
