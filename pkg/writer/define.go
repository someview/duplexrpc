package writer

type Writer int

const (
	ShardQueue Writer = iota + 1
	BatchWriter
)
