package framework

import (
	"github.com/JijiaZan/godml/utils"
	"time"
)

const HeartbeatInterval = time.Millisecond * 20

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

type phase int
const (
	INITED = iota
	DATA_ASSIGNED
	DATA_PREPEARED
	TRANNING
	FINISHED
)

type Node struct {
	id int
	role role
	address string
	IsAlive bool
	schAddress string
}

func (node *Node)SendHeartbeat(msg string) {
	args := &HeartbeatArgs{
		ID: node.id,
		Msg: msg,
	}
	reply := &HeartbeatReply{}
	if ok := utils.Call("Scheduler.Heartbeat", args, reply, node.schAddress); !ok {
		utils.DPrintf("Send Heartbeat failed")
	}
}