package main

import (
	"fmt"
	"os"
	"time"
	"../framework"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "You have to specify the command \n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		// godml start task_name
		sch := framework.MakeScheduler(os.Args[2])
		fmt.Println("Scheduler start working")
		for !sch.IsDone {
			time.Sleep(time.Duration(10) * time.Second)
		}
		fmt.Println("Scheduler stop working")
	case "addWorker":
		// godml addWorker port xx.so
		// xx.so is the plugin of preprocess
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "You have to provide the preprocess func \n")
			os.Exit(1)
		}
		prepFun := framework.LoadPlugin(os.Args[3])
		w := framework.MakeWorker(os.Args[2], prepFun)
		fmt.Println("worker start working")
		for !w.IsAlive {
			time.Sleep(time.Duration(10) * time.Second)
		}
	case "test":

	default:
		fmt.Fprintf(os.Stderr, "command not found \n")
	}
}