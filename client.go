package main

import (
	"flag"
	"fmt"
	"github.com/jserver/serverstyle/server"
	"log"
	"net/rpc"
	"strconv"
)

var (
	host = flag.String("host", "localhost", "Host to Dial")
	port = flag.Int("port", 1234, "Port serverstyle is running on")
	cmd = flag.String("cmd", "test", "Comamnd to run")
)

var (
	args interface{}
	results server.Results
	command string
)

func main() {
	flag.Parse()
	//address := fmt.Sprintf("%s:%d", *host, *port)
	address := *host + ":" + strconv.Itoa(*port) 
	fmt.Println("Calling: " + address)
	packages := flag.Args()

	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	if *cmd == "install" {
		args = &server.AptGetArgs{packages}
		results = new(server.AptGetResults)
		command = "AptGet.Install"

	} else if *cmd == "script" {
		args = &server.ScriptArgs{"script_test.sh", "#!/bin/bash\nls -al\n"}
		results = new(server.ScriptResults)
		command = "Script.Runner"

	} else {
		args = &server.TestArgs{packages}
		results = new(server.TestResults)
		command = "Test.Runner"
	}

	remoteCall := client.Go(command, args, results, nil)
	<-remoteCall.Done
	errText := results.GetErr()
	if len(errText) > 0 {
		fmt.Println(">>> [", errText, "]")
	}
	fmt.Println(results.GetOutput())
}
