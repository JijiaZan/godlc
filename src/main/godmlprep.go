package main

import (
	"fmt"
	"os"
	"../framework"
)

const ADDRESS = "127.0.0.1:1234"

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "You have to specify the command \n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload":
		// godmlprep upload train/validation dirOfData
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "You have to specify the dir of dataset \n")
			os.Exit(1)
		}

		args := &framework.UploadArgs{}
		args.Dir = os.Args[3]
		reply := &framework.UploadReply{}

		switch os.Args[2] {
		case "train":
			args.Dt = framework.Train
		case "validation":
			args.Dt = framework.Validation
		default:
			fmt.Fprintf(os.Stderr, "Wrong data type\n")
		}

		if ok := framework.Call("Scheduler.Upload", args, reply, ADDRESS); !ok {
			framework.DPrintf("upload file failed")
			os.Exit(1)
		} else {
			framework.DPrintf("Assign data successfully")
		}
	case "preprocess":
		//godmlprep preprocess train/validation
		args := &framework.PreprocessArgs{}
		reply := &framework.PreprocessReply{}

		switch os.Args[2] {
		case "train":
			args.Dt = framework.Train
		case "validation":
			args.Dt = framework.Validation
		default:
			fmt.Fprintf(os.Stderr, "Wrong data type\n")
		}

		if ok := framework.Call("Scheduler.Preprocess", args, reply, ADDRESS); !ok {
			framework.DPrintf("Preprocess data failed")
			os.Exit(1)
		} else {
			framework.DPrintf("Start preprocessing data")
		}

	default:
		fmt.Fprintf(os.Stderr, "command not found \n")
	}

}