package framework

import "github.com/JijiaZan/godml/pyserver"

type AddNodeArgs struct {
	Address string
	Role role
	Rank int
}
type AddNodeReply struct {
	ID int
}


type UploadArgs struct {
	Dir string
	Dt DataType
}
type UploadReply struct {}


type AssignDataArgs struct {
	Dir string
	FileName []string
	Dt DataType
}
type AssignDataReply struct {}


type PreprocessArgs struct {
	Dt DataType
}
type PreprocessReply struct{}


type HeartbeatArgs struct{
	ID int
	Msg string
}
type HeartbeatReply struct{}

type PushArgs struct{
	Gradients []*pyserver.Layer
	Epoch int
	Rank int
}
type PushReply struct{}

type PullArgs struct{}
type PullReply struct{
	Weights []*pyserver.Layer
}
