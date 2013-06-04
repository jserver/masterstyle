package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jserver/serverstyle/server"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strings"
)

type Config struct {
    DomainName      string                       `json:"domain_name"`
	KeyPath         string                       `json:"key_path"`
	PasswordSalt    string                       `json:"password_salt"`
	ScriptPath      string                       `json:"script_path"`
	Images          map[string]string            `json:"images"`
	Bundles         map[string]string            `json:"bundles"`
	PythonBundles   map[string]string            `json:"python_bundles"`
	AptRepositories map[string]map[string]string `json:"apt_repositories"`
	PPAs            map[string]string            `json:"personal_package_archives"`
	Groups          map[string]string            `json:"groups"`
	Builds          map[string]map[string]string `json:"builds"`
	Projects        map[string]map[string]string `json:"projects"`
	DebConf         map[string]map[string]string `json:"debconf"`
}

var (
	host = flag.String("host", "localhost", "Host to Dial")
	port = flag.Int("port", 1234, "Port serverstyle is running on")

	home = os.Getenv("HOME")
	config Config

	command string
	cmdArgs interface{}
	results server.Results
)

func main() {
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Println("Calling: " + address)
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("No Command Given")
	}
	cmd := args[0]		


	data, err := ioutil.ReadFile(home + "/.clifford.json")
	if err != nil {
		log.Fatal("unable to read file")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("unable to parse config:", err)
	}

	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	if cmd == "install" {
		if len(args) != 2 {
			log.Fatal("No bundle given")
		}
		bundle := config.Bundles[args[1]]
		packages := strings.Split(bundle, " ")
		
		cmdArgs = &server.AptGetArgs{packages}
		results = new(server.AptGetResults)
		command = "AptGet.Install"

	} else if cmd == "script" {
		cmdArgs = &server.ScriptArgs{"script_test.sh", "#!/bin/bash\nls -al\n"}
		results = new(server.ScriptResults)
		command = "Script.Runner"

	} else if cmd == "test" {
		packages := []string{"A", "B", "C"}
		cmdArgs = &server.TestArgs{packages}
		results = new(server.TestResults)
		command = "Test.Runner"
	} else {
		log.Fatal("Command Unknown")
	}

	remoteCall := client.Go(command, cmdArgs, results, nil)
	<-remoteCall.Done
	errText := results.GetErr()
	if len(errText) > 0 {
		fmt.Println(">>> [", errText, "]")
	}
	fmt.Println(results.GetOutput())
}
