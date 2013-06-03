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
		args := &server.AptGetArgs{packages}
		results := new(server.AptGetResults)
		aptGetCall := client.Go("AptGet.Install", args, results, nil)
		<-aptGetCall.Done
		if len(results.Err) > 0 {
			fmt.Println(">>> [", results.Err, "]")
		}
		fmt.Println(string(results.Output))
	} else {
		args := &server.TestArgs{packages}
		results := new(server.TestResults)
		testCall := client.Go("Test.Runner", args, results, nil)
		<-testCall.Done
		if len(results.Err) > 0 {
			fmt.Println(">>> [", results.Err, "]")
		}
		fmt.Println(string(results.Output))
	}

}
