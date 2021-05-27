package framework

import(
	"github.com/JijiaZan/godml/utils"
	"github.com/JijiaZan/godml/pyserver"
	"time"
	"sync"
)

type Server struct {
	Node

	nWorker int

	weights []*pyserver.Layer

	vLock sync.Mutex
	vectorClock []int
	tao int
	writingRank int

}

func MakeServer(address string, schAddress string, rank int, numOfWorker int, tao int) *Server{
	s := &Server{}

	s.address = address
	s.schAddress = schAddress
	s.role = SERVER
	s.IsAlive = true

	s.vLock = sync.Mutex{}

	s.nWorker = numOfWorker
	s.tao = tao
	s.vectorClock = make([]int, s.nWorker)

	args := &AddNodeArgs{s.address, s.role, rank}
	reply := &AddNodeReply{}
	if ok := utils.Call("Scheduler.AddNode", args, reply, s.schAddress); !ok {
		utils.DPrintf("add server fail, please check the scheduler's status")
	} else {
		utils.DPrintf("add server successfully")
		s.id = reply.ID
		utils.DPrintf("ID: %d", s.id)
	}

	port := utils.GetPort(s.address)
	utils.Serve(s, port)

	go func() {
		for {
			time.Sleep(HeartbeatInterval)
			s.SendHeartbeat("")
		}
	}()

	return s
}

func (s *Server) Pull(args *PullArgs, reply *PullReply) error {
	reply.Weights = s.weights
	return nil
}

func (s *Server) Push(args *PushArgs, reply *PushReply) error {
	if args.Epoch == -1 {
		if s.weights == nil {
			s.vLock.Lock()
			defer s.vLock.Unlock()
			s.weights = args.Gradients
		} else {
			return nil
		}
	} else {
		s.Updata(args)
	}
	return nil
}

func (s *Server) Updata(args *PushArgs) {
	for {
		// double check
		//utils.DPrintf("out: %d - %d < %d", args.Epoch, s.tao, s.getOldestEpoch())
		if args.Epoch - s.tao <= s.getOldestEpoch() {
			s.vLock.Lock()
			//utils.DPrintf("get lock")
			//utils.DPrintf("in: %d - %d < %d", args.Epoch, s.tao, s.getOldestEpoch())
			if args.Epoch - s.tao <= s.getOldestEpoch() {
				// 更新网络
				for i := 0; i < len(args.Gradients); i++ {
					s.weights[i] = utils.CalGradient(s.weights[i], args.Gradients[i])
				}
				s.vectorClock[args.Rank] = args.Epoch
				s.vLock.Unlock()
				break
			}
			s.vLock.Unlock()
			//utils.DPrintf("release lock")
		}
		time.Sleep(time.Millisecond * 20)
	} 
}

func (s *Server) getOldestEpoch() int {
	min := int(^uint(0) >> 1);
	for _, n := range(s.vectorClock) {
		if n < min {
			min = n
		}
	}
	return min
}

