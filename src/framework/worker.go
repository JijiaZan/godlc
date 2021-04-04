package framework

import(
	"os"
	"os/user"
	"strconv"
	"../utils"
)

//const HOST_ADDRESS = "127.0.0.1:1234"

type Worker struct {
	id int
	address string
	IsAlive bool
	schAddress string

	prepFun func(sourceDir, targetDir string)
	trainSet []string
	validationSet []string
}

func MakeWorker(address string, schAddress string, rank int, pref func(string, string)) *Worker {
	w := &Worker{}
	w.address =  address//ip要自己读取
	w.prepFun = pref
	w.schAddress = schAddress

	args := &AddNodeArgs{w.address, WORKER, rank}
	reply := &AddNodeReply{}

	port := utils.GetPort(w.address)
	utils.Serve(w, port)

	if ok := utils.Call("Scheduler.AddNode", args, reply, w.schAddress); !ok {
		utils.DPrintf("add worker fail, please check the scheduler's status")
	} else {
		utils.DPrintf("add worker successfully")
		w.id = reply.ID
		utils.DPrintf("ID: %d", w.id)
	}


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



