package protocol

import "sync/atomic"

var defaultSeqGenerator = &seqGenerator{}

type seqGenerator struct {
	seq atomic.Int32
}

func (g *seqGenerator) Next() uint32 {
	id := g.seq.Add(1)
	if id == 0 {
		id = g.seq.Add(1)
	}
	return uint32(id)
}

// GenSeqId NextSeqId 使用的是全局生成器
func GenSeqId() uint32 {
	return defaultSeqGenerator.Next()
}

func NewSeqIdGenerator() SeqIdGenerator {
	return &seqGenerator{}
}
