package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
	"launchpad.net/goamz/s3"
	"log"
	"os"
	"strconv"
	"strings"
)

type Build struct {
	Size           string
	Image          string
	Key            string
	Zone           string
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
	BucketName      string
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
	scriptPath string
	config Config

	conn      *ec2.EC2
	instances map[string]ec2.Instance

	bucket *s3.Bucket
)

func AskQuestion(question string) string {
	for {
		fmt.Print(question)
		bytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal("Unable to read answer!")
		}
		if len(bytes) == 0 {
			continue
		}
		return string(bytes)
	}
}

type Answer struct {
	Text  string
	Value string
}

func AskMultipleChoice(question string, answers []Answer) string {
	for {
		for idx, answer := range answers {
			fmt.Printf("%d) %s\n", idx+1, answer.Text)
		}
		fmt.Print(question)
		bytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal("Unable to read multiple choice answer!")
		}
		if len(bytes) == 0 {
			continue
		}
		num, err := strconv.Atoi(string(bytes))
		if err != nil {
			fmt.Println("Unable to convert to int")
			continue
		}
		if num > 0 && num <= len(answers) {
			return answers[num-1].Value
		}
	}
}

func help() {
	fmt.Println(`The following commands are available
	exit,q,quit         Leave Program
	
	ls, status          List all machines
	launch              Create a new ec2 instance from specified build
	update              apt-get update
	upgrade             apt-get upgrade
	install             Install a bundle from json config
	script              Run a script on the machine specified
	test                Pass in a few directories to see if ServerStyle is responding
	easy, easy_install  Install a python bundle from json config

	EC2 Actions on a specified Machine
	----------------------------------
	reboot              Reboot specified machine
	start               Start specified machine
	stop                Stop specified machine
	terminate           Terminate specified mach

	S3 Actions
	----------
	s3upload            Upload ServerStyle to s3 for new machines

	IP Addresses
	------------
	addresses           List all IP Addresses allocated
	associate           Attach an instance to one of your allocated addresses
	disassociate        Release an address from an instance
	`)
}

func main() {
	// Config File
	data, err := ioutil.ReadFile(home + "/.serverstyle/config.json")
	if err != nil {
		log.Fatal("Unable to read config file")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Unable to parse config:", err)
	}

	// ScriptPath
	scriptPath = os.ExpandEnv(config.ScriptPath)

	// AWS Auth
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal("AWS AUTH FAIL!")
	}

	// EC2 Connection
	conn = ec2.New(auth, aws.USEast)

	// S3 Connection/Bucket
	s3conn := s3.New(auth, aws.USEast)
	bucket = s3conn.Bucket(config.BucketName)

	// Instances
	instances = GetInstances()

	// Command Loop
	for {
		line := AskQuestion("(MasterStyle): ")
		parts := strings.Split(line, " ")
		cmd := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		switch cmd {
		case "h", "help":
			help()

		case "exit", "q", "quit":
			fmt.Println("Bye-Bye")
			return

		case "launch":
			Launch(args)

		case "ls", "status":
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

		case "easy", "easy_install":
			EasyInstall(args)

		case "s3upload":
			S3Upload()

		case "addresses":
			Addresses()
		case "associate":
			Associate(args)
		case "disassociate":
			Disassociate(args)

		default:
			fmt.Println("Command Not Found!")
		}
		continue
	}
}
