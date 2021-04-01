package framework

import(
	"os"
	"os/user"
	"strconv"
	//"../prep"
)

const HOST_ADDRESS = "127.0.0.1:1234"

type Worker struct {
	id int
	address string
	IsAlive bool

	prepFun func(sourceDir, targetDir string)
	trainSet []string
	validationSet []string
}

func MakeWorker(port string, pref func(string, string)) *Worker {
	w := &Worker{}
	w.address =  "127.0.0.1:" +  port //ip要自己读取
	w.prepFun = pref

	args := &AddWorkerArgs{w.address}
	reply := &AddWorkerReply{}

	if ok := Call("Scheduler.AddWorker", args, reply, HOST_ADDRESS); !ok {
		DPrintf("add worker fail, please check the scheduler's status")
	} else {
		w.id = reply.ID
		DPrintf("Success! id: %d", reply.ID)
	}

	serve(w, port)
	return w
}

func (w *Worker) AssignData(args *AssignDataArgs, reply *AssignDataReply) error {
	switch args.Dt {
	case Train:
		for _, s := range(args.FileName) {
			w.trainSet = append(w.trainSet, args.Dir + "/" + s)
			DPrintf(args.Dir + "/" + s)
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

	for _, s := range(data) {
		w.prepFun(s, dir)
	}
	
	return nil
}



