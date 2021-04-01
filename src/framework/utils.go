package framework

import(
	"log"
	"net"
	"fmt"
	"net/http"
	"net/rpc"
	"plugin"
)

type DataType int
const (
	Train = iota
	Validation
)

const DEBUG = true

func serve(obj interface{}, port string) {
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

	fmt.Println(err)
	return false
}

func DPrintf(format string, v ...interface{}) {
	if DEBUG {
		log.Printf(format+"\n", v...)
	}
}

func LoadPlugin(filename string) func(string, string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot open plugin %v", filename)
	}
	xprep, err := p.Lookup("Preprocess")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	prepf := xprep.(func(string, string))

	return prepf
}