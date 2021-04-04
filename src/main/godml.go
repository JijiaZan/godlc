package main

import (
	"fmt"
	"os"
	"time"
	"../framework"
	"../utils"
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
			sch := framework.MakeScheduler(conf.Scheduler.Address, len(conf.Workers))
			fmt.Println("Scheduler start working")
			for !sch.IsDone {
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
				fmt.Fprintf(os.Stderr, "The id of worker does not exist in worker \n")
				os.Exit(1)
			}
			prepFun := utils.LoadPrepPlugin("../.." + conf.Workers[i].Prepf)
			w := framework.MakeWorker(conf.Workers[i].Address, conf.Scheduler.Address, i, prepFun)
			for !w.IsAlive {
				time.Sleep(time.Duration(10) * time.Second)
			}
		//case "server":
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