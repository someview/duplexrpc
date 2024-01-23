package protocol

type FrameType byte

const (
	InitialType     FrameType = iota
	InitialDoneType           = 1
	DataType                  = 2
	PingType                  = 6
	PongType                  = 7
	CloseType                 = 8
)

// DataFrame [Frame Type (1 byte)] [msgType (1 bytes)] [msgLen (4 length)][flag][optional][payload]
// 需要有一个控制信号, 用于在检测时，完成这个处理过程
type DataFrame struct {
	Header  []byte
	Payload []byte
	Done    chan struct{}
}

func NewDataFrame() DataFrame {
	return DataFrame{}
}
