package main

import (
	"fmt"
	"os"
	"time"
	"github.com/JijiaZan/godml/framework"
	"github.com/JijiaZan/godml/utils"
	"strconv"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "You have to specify the command \n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		conf := utils.LoadGlobalConfig()
		// godml start task_name
		switch os.Args[2] {

		case "scheduler":
			nWorkers := (len(conf.Workers)-1)*2 + 1
			if (len(conf.Servers)-1)*2 > nWorkers {
				nWorkers = (len(conf.Servers)-1)*2
			}
			nWorkers += 1
			//utils.DPrintf("nums: %d", nWorkers)

			sch := framework.MakeScheduler(conf.Scheduler.Address, nWorkers)
			fmt.Println("Scheduler start working")
			for sch.Phase != framework.FINISHED {
				time.Sleep(time.Duration(10) * time.Second)
			}
			fmt.Println("Scheduler stop working")

		case "worker":
			if len(os.Args) < 4 {
				fmt.Fprintf(os.Stderr, "You should provide the id of worker \n")
				os.Exit(1)
			}
			i, _ := strconv.Atoi(os.Args[3])
			if i >= len(conf.Workers) {
				fmt.Fprintf(os.Stderr, "The rank of worker does not exist in worker \n")
				os.Exit(1)
			}

			prepFun := utils.LoadPrepPlugin("../.." + conf.Workers[i].Prepf)

			sAdress := make([]string, len(conf.Servers))
			for i, s := range(conf.Servers) {
				sAdress[i] = s.Address
			}

			w := framework.MakeWorker(conf.Workers[i].Address, conf.Scheduler.Address, i, prepFun, sAdress)
			for w.IsAlive {
				time.Sleep(time.Duration(10) * time.Second)
			}

		case "server":
			if len(os.Args) < 4 {
				fmt.Fprintf(os.Stderr, "You should provide the rank of server \n")
				os.Exit(1)
			}
			i, _ := strconv.Atoi(os.Args[3])
			if i >= len(conf.Servers) {
				fmt.Fprintf(os.Stderr, "The rank of worker does not exist in server \n")
				os.Exit(1)
			}

			nWorker := len(conf.Workers)
			tao, _ := strconv.Atoi(conf.Scheduler.Consistency)
			s := framework.MakeServer(conf.Servers[i].Address, conf.Scheduler.Address, i, nWorker, tao)
			for s.IsAlive {
				time.Sleep(time.Duration(10) * time.Second)
			}
		}

	// case "addWorker":
	// 	// godml addWorker port xx.so
	// 	// xx.so is the plugin of preprocess
	// 	if len(os.Args) < 4 {
	// 		fmt.Fprintf(os.Stderr, "You have to provide the preprocess func \n")
	// 		os.Exit(1)
	// 	}
	// 	prepFun := utils.LoadPrepPlugin(os.Args[3])
	// 	w := framework.MakeWorker(os.Args[2], prepFun)
	// 	fmt.Println("worker start working")
	// 	for !w.IsAlive {
	// 		time.Sleep(time.Duration(10) * time.Second)
	// 	}
	// case "test":

	default:
		fmt.Fprintf(os.Stderr, "command not found \n")
	}
}