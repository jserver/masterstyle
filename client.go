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

type Build struct {
	Size           string
	Image          string
	Key            string
	Zone           string
	SecurityGroups []string
	UserData       string
	Upgrade        string
	Group          string
	EasyInstall    string
	ScriptName     string
	ScriptAction   string
}

type Project struct {
	Build        string
	UserName     string
	UserFullname string
	UserPassword string
	ScriptName   string
	ScriptAction string
}

type Repository struct {
	KeyServer string
	PublicKey string
	DebLine   string
	Package   string
}

type Config struct {
	DomainName      string
	KeyPath         string
	PasswordSalt    string
	ScriptPath      string
	Images          map[string]string
	Bundles         map[string]string
	PythonBundles   map[string]string
	AptRepositories map[string]Repository `json:"Apt:Repositories"`
	PPAs            map[string]string
	Groups          map[string]string
	Builds          map[string]Build
	Projects        map[string]Project
	DebConfs        map[string][]string
}

var (
	host = flag.String("host", "localhost", "Host to Dial")
	port = flag.Int("port", 1234, "Port serverstyle is running on")

	home   = os.Getenv("HOME")
	config Config

	command string
	cmdArgs interface{}
	results server.Results
)

func main() {
	flag.Parse()
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

	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Println("Calling: " + address)
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
		if len(args) != 2 {
			log.Fatal("No script given")
		}
		script := os.ExpandEnv(args[1])
		content, err := ioutil.ReadFile(script)
		if err != nil {
			log.Fatal("unable to read file")
		}

		script_parts := strings.Split(script, "/")
		script_name := script_parts[len(script_parts)-1]

		cmdArgs = &server.ScriptArgs{script_name, content}
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
