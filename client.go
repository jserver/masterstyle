package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
	"launchpad.net/goamz/s3"
)

type NamedInstance struct {
	Name string
	ec2.Instance
}

var (
	reader     = bufio.NewReader(os.Stdin)
	home       = os.Getenv("HOME")
	scriptPath string
	config     *Config

	conn      *ec2.EC2
	instances []*NamedInstance

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

		case "security":
			ListSecurityGroups()

		case "tag":
			Tag()

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
		case "ppa":
			PPAInstall(args)
		case "script":
			Script(args)
		case "test":
			Test(args)

		case "easy_install":
			EasyInstall(args)

		case "s3upload":
			S3Upload()

		case "addresses":
			Addresses()
		case "associate":
			Associate()
		case "disassociate":
			Disassociate()

		case "create_image":
			CreateImage(args)

		default:
			fmt.Println("Command Not Found!")
		}
		continue
	}
}
