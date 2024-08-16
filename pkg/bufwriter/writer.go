package bufwriter

import (
	"time"

	netpoll "github.com/cloudwego/netpoll"
)

func NewBufWriter(typ BufWriterType, bufferThreshold int, maxMsgNum int, delayTime time.Duration, conn netpoll.Connection) BufWriter {
	var wr BatchWriter
	switch typ {
	case ShardQueueType:
		wr = NewShardQueue(bufferThreshold, conn)
	case DelayQueueType:
		wr = NewBatchWriter(bufferThreshold,
			maxMsgNum,
			delayTime,
			conn,
			func(n int, err error) {
				// 一般来说不会发生，这里错误时，连接可以感知
			})
	}
	return wr
}
