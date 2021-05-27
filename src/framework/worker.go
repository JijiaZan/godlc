package framework

import(
	"os"
	"os/user"
	"strconv"
	"github.com/JijiaZan/godml/utils"
	"time"
	"context"
	"github.com/JijiaZan/godml/pyserver"

	"google.golang.org/grpc"
	"net"
	"log"
)

//const HOST_ADDRESS = "127.0.0.1:1234"


type Worker struct {
	pyserver.UnimplementedWorkerServer
	// id int
	// address string
	// IsAlive bool
	// schAddress string
	Node

	// 服务器地址
	serverAdress []string

	// 预处理
	prepFun func(sourceDir, targetDir string)
	trainSet []string
	validationSet []string

	// 神经网络参数
	layers []*pyserver.Layer
	lastEpochLayers []*pyserver.Layer
}

func MakeWorker(address string, schAddress string, rank int, pref func(string, string), sAdress []string) *Worker {
	w := &Worker{}
	w.address =  address
	w.schAddress = schAddress
	w.role = WORKER
	w.IsAlive = true

	w.prepFun = pref
	w.serverAdress = sAdress
	for _, s := range(w.serverAdress) {
		utils.DPrintf(s)
	}
	
	args := &AddNodeArgs{w.address, w.role, rank}
	reply := &AddNodeReply{}
	if ok := utils.Call("Scheduler.AddNode", args, reply, w.schAddress); !ok {
		utils.DPrintf("add worker fail, please check the scheduler's status")
	} else {
		utils.DPrintf("add worker successfully")
		w.id = reply.ID
		utils.DPrintf("ID: %d", w.id)
	}

	port := utils.GetPort(w.address)
	utils.Serve(w, port)

	// 发送心跳
	go func() {
		for {
			time.Sleep(HeartbeatInterval)
			w.SendHeartbeat("")
		}
	}()
	
	//python服务
	nextPort := utils.GetNextPort(w.address)
	go w.GServe(nextPort)

	return w
}

func (w *Worker) AssignData(args *AssignDataArgs, reply *AssignDataReply) error {
	switch args.Dt {
	case Train:
		for _, s := range(args.FileName) {
			w.trainSet = append(w.trainSet, args.Dir + "/" + s)
			utils.DPrintf(args.Dir + "/" + s)
		}
	case Validation:
		for _, s := range(args.FileName) {
			w.validationSet = append(w.validationSet, args.Dir + "/"  + s)
		}
	}
	return nil
}

func (w *Worker) Preprocess(args *PreprocessArgs, reply *PreprocessReply) error {
	u, _ := user.Current()
	dir := u.HomeDir + "/godml"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	
	dir += "/preprocessed_" + strconv.Itoa(w.id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	var data []string
	if args.Dt == Train {
		data = w.trainSet
		dir += "/train/"
	} else {
		data = w.validationSet
		dir += "/validation/"
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	go func() {
		for _, s := range(data) {
			w.prepFun(s, dir)
		}
	}()
	return nil
}

func (w *Worker) Push(ctx context.Context, in *pyserver.PushRequest) (*pyserver.PushReply, error) {

	args := &PushArgs{}
	reply := &PushReply{}

	args.Epoch = int(in.GetEpoch())
	args.Rank = (w.id - 1) / 2

	utils.DPrintf("epoch: %+v", in.Epoch)
	utils.DPrintf("len: %+v", len(in.Layers))

	if (in.Epoch == -1) {

		args.Gradients = in.Layers
		// 推送完整的
	} else {
		//计算梯度
		gradients := make([]*pyserver.Layer, len(in.Layers))
		for i, layer := range(in.Layers) {
			gradients[i] = utils.CalGradient(w.lastEpochLayers[i], layer)
		}
		args.Gradients = gradients
	}

	w.lastEpochLayers = in.Layers // 迭代上一次
	
	if ok := utils.Call("Server.Push", args, reply, w.serverAdress[0]); !ok {
		utils.DPrintf("Push server failed")
	} else {
		utils.DPrintf("Push server successfully")
	}

	return &pyserver.PushReply{Success: false}, nil
}

func (w *Worker) Pull(ctx context.Context, in *pyserver.PullRequest) (*pyserver.PullReply, error) {
	utils.DPrintf("pull epoch: %d",in.Epoch)

	args := &PullArgs{}
	reply := &PullReply{}

	layers := &pyserver.PullReply{}

	if ok := utils.Call("Server.Pull", args, reply, w.serverAdress[0]); !ok {
		utils.DPrintf("Pull server failed")
	} else {
		utils.DPrintf("Pull server successfully")
		layers.Layers = reply.Weights
	}
	return layers, nil
} 

func (w *Worker) GServe(port string) {
	lis, err := net.Listen("tcp",  ":" + port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pyserver.RegisterWorkerServer(s, w)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}