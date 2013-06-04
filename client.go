package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jserver/serverstyle/server"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
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

func launch() {
    auth, err := aws.EnvAuth()
    if err != nil {
        log.Fatal("AWS AUTH Fail")
    }
    e := ec2.New(auth, aws.USEast)

	secGroups := []ec2.SecurityGroup{{Name: "django-dev"}, {Name: "serverstyle"}}

    options := ec2.RunInstances{
        ImageId:            "ami-f3d1bb9a",
        InstanceType:       "t1.micro",
		KeyName:            "jserver",
		PlacementGroupName: "us-east-1b",
		SecurityGroups:     secGroups,
    }
    resp, err := e.RunInstances(&options)
    if err != nil {
		log.Fatal("AWS ec2 Run Instances Fail")
    }

    for _, instance := range resp.Instances {
        println("Now running", instance.InstanceId)
    }
    println("Make sure you terminate instances to stop the cash flow.")
}

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

	if cmd == "launch" {
		fmt.Println("GETTING READY TO LAUNCH")
		launch()
		fmt.Println("DID WE LAUNCH")

	} else if cmd == "install" {
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
