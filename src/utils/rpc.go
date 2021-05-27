package utils

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	//"google.golang.org/grpc"
	//"github.com/JijiaZan/godml/pyserver"
	//"github.com/JijiaZan/godml/framework"
)

func Serve(obj interface{}, port string) {
	rpc.Register(obj)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":" + port)
	if e != nil {
		log.Fatal("listen error:", e)
	} else {
		DPrintf("listenning")
	}
	go http.Serve(l, nil)
}

func Call(rpcname string, args interface{}, reply interface{}, address string) bool {
	c, err := rpc.DialHTTP("tcp", address)

	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	log.Fatal("listen error:", err)
	return false
}

// func gServe(obj &framework.Worker, port string) {
// 	lis, err := net.Listen("tcp",  ":" + port)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	pyserver.RegisterWorkerServer(s, obj)
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }