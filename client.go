package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
	"log"
	"os"
	"strings"
)

type Build struct {
	Size           string
	Image          string
	Key            string
	Placement      string
	SecurityGroups []ec2.SecurityGroup
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

type PPA struct {
	Package string
	Source  string
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
	PPAs            []PPA
	Groups          map[string]string
	Builds          map[string]Build
	Projects        map[string]Project
	DebConfs        map[string][]string
}

var (
	reader = bufio.NewReader(os.Stdin)
	home   = os.Getenv("HOME")
	config Config

	conn *ec2.EC2
	instances map[string]ec2.Instance
)

func main() {

	// Config File
	fmt.Println("Reading config file")
	data, err := ioutil.ReadFile(home + "/.clifford.json")
	if err != nil {
		log.Fatal("Unable to read file")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Unable to parse config:", err)
	}

	// EC2 Connection
	fmt.Println("Establishing connection with aws")
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal("AWS AUTH FAIL!")
	}
	conn = ec2.New(auth, aws.USEast)

	// Instances
	fmt.Println("Getting Instances")
	instances = GetInstances()

	// Command Loop
	for {
		print("(MasterStyle): ")
		bytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal("Unable to read command!")
		}
		if len(bytes) == 0 {
			continue
		}
		line := string(bytes)

		parts := strings.Split(line, " ")
		cmd := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		switch cmd {
			case "help":
				fmt.Println("Help Message")

			case "exit":
				fallthrough
			case "quit":
				fmt.Println("Bye-Bye")
				return

			case "launch":
				Launch(args)

			case "status":
				Status()

			case "reboot":
				Reboot(args)
			case "start":
				Start(args)
			case "stop":
				Stop(args)
			case "terminate":
				Terminate(args)

			case "update":
				Update(args)
			case "upgrade":
				Upgrade(args)
			case "install":
				Install(args)
			case "script":
				Script(args)
			case "test":
				Test(args)

			default:
				fmt.Println("Command Not Found!")
		}
		continue
	}
}
