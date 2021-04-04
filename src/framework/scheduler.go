package framework

import (
	"sync"
	"io/ioutil"
	"strings"
	"../utils"
	"errors"
)

type WorkerStat struct {
	address string
	IsAlive bool
}

type Scheduler struct {
	mu sync.Mutex

	Name string
	WorkerStats []WorkerStat //true表示存活，

	IsWorking bool
	IsDone bool
}


func MakeScheduler(address string, nWorkers int) *Scheduler{
	sch := &Scheduler{}
	sch.mu = sync.Mutex{}

	sch.WorkerStats = make([]WorkerStat, nWorkers)
 	sch.IsWorking = false
	sch.IsDone = false
	
	utils.Serve(sch, "1234")
	return sch
}

func (sch *Scheduler) AddWorker(args *AddWorkerArgs, reply *AddWorkerReply) error {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if sch.WorkerStats[args.ID].IsAlive == true {
		return errors.New("This worker is already running")
	}

	w := WorkerStat{
		address: args.Address,
		IsAlive: true,
	}

	sch.WorkerStats[args.ID] = w
	utils.DPrintf("The new worker id is: %d", args.ID)
	return nil
}

//之后可以用hdfs的文件指令更换这部分本地文件操作
func (sch *Scheduler) Upload(args *UploadArgs, reply *UploadReply) error {
	fileInfos, err := ioutil.ReadDir(args.Dir)
	if err != nil {
        return err
    }
	
	assignedData := [][]string{}
	for i := 0; i<len(sch.WorkerStats); i++ {
		assignedData = append(assignedData, []string{})
	}

	idx := 0
	for idx < len(fileInfos) {
		for i, w := range(sch.WorkerStats) {
			if idx >= len(fileInfos) {
				break
			}
			if !w.IsAlive {
				continue
			}
			if !strings.HasPrefix(fileInfos[idx].Name(), ".") {
				assignedData[i] = append(assignedData[i], fileInfos[idx].Name())
			}
			idx ++
		}
	}

	for i, w := range(sch.WorkerStats) {
		if !w.IsAlive {
			continue
		}
		adArgs := &AssignDataArgs {
			Dir: args.Dir,
			FileName: assignedData[i],
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
	for _, w := range(sch.WorkerStats) {
		if !w.IsAlive {
			continue
		}
		if ok := utils.Call("Worker.Preprocess", args, reply, w.address); !ok {
			utils.DPrintf("Preprocess data failed")
		}
	}
	return nil
}


func (sch *Scheduler) CheckStats(args *CheckNode)
