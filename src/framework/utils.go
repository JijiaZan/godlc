package framework

type DataType int
const (
	Train = iota
	Validation
)

type role int
const (
	SERVER = 0
	WORKER = 1
)