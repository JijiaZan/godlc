package framework

import (
	"sync"
	"io/ioutil"
	"strings"
	"../utils"
	"errors"
	//"strconv"
	//"time"
)

type NodeStat struct {
	address string
	IsAlive bool
}

type Scheduler struct {
	mu sync.Mutex

	Name string
	NodeStats []NodeStat //true表示存活，

	IsWorking bool
	IsDone bool
}


func MakeScheduler(address string, nNodes int) *Scheduler{
	sch := &Scheduler{}
	sch.mu = sync.Mutex{}

	sch.NodeStats = make([]NodeStat, nNodes)
 	sch.IsWorking = false
	sch.IsDone = false
	
	utils.Serve(sch, "1234")
	return sch
}

func (sch *Scheduler) AddNode(args *AddNodeArgs, reply *AddNodeReply) error {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	id := (int)(args.Role) + 2 * args.Rank

	if sch.NodeStats[id].IsAlive == true {
		return errors.New("This worker is already running")
	}
	w := NodeStat{
		address: args.Address,
		IsAlive: true,
	}

	sch.NodeStats[id] = w
	utils.DPrintf("The new worker id is: %d", id)
	reply.ID = id
	return nil
}

//之后可以用hdfs的文件指令更换这部分本地文件操作
func (sch *Scheduler) Upload(args *UploadArgs, reply *UploadReply) error {
	fileInfos, err := ioutil.ReadDir(args.Dir)
	if err != nil {
        return err
    }
	
	assignedData := make([][]string, (len(sch.NodeStats)+1) / 2)

	idx := 0
	for idx < len(fileInfos) {
		for i, w := range(sch.NodeStats) {
			if idx >= len(fileInfos) {
				break
			}
			if i & 1 == 0 || !w.IsAlive {
				continue
			}
			if !strings.HasPrefix(fileInfos[idx].Name(), ".") {
				assignedData[i/2] = append(assignedData[i/2], fileInfos[idx].Name())
			}
			idx ++
		}
	}

	for i, w := range(sch.NodeStats) {
		if i & 1 == 0 || !w.IsAlive {
			continue
		}
		adArgs := &AssignDataArgs {
			Dir: args.Dir,
			FileName: assignedData[i/2],
			Dt: args.Dt,
		}
		adReply := &AssignDataReply{}
		if ok := utils.Call("Worker.AssignData", adArgs, adReply, w.address); !ok {
			utils.DPrintf("Assign data failed")
		}
	}

	return nil
}

func (sch *Scheduler) Preprocess( args *PreprocessArgs, reply *PreprocessReply) error {
	for id, w := range(sch.NodeStats) {
		if id & 1 == 0 || !w.IsAlive {
			continue
		}
		if ok := utils.Call("Worker.Preprocess", args, reply, w.address); !ok {
			utils.DPrintf("Preprocess data failed")
		}
	}
	return nil
}


// func (sch *Scheduler) CheckStats(args *CheckNode)
//func (sch *Scheduler) HeatBeat(args *HBArgs)
