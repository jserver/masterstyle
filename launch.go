package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goamz/ec2"
	"log"
)

func Launch(args []string) {
	if len(args) != 1 {
		fmt.Println("no build given!")
		return
	}
	build := args[0]

	fmt.Print("Enter Tag Name: ")
	line, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal("Unable to read tag value!")
	}
	tags := []ec2.Tag{{"Name", string(line)}}

	options := ec2.RunInstances{
		ImageId:            config.Images[config.Builds[build].Image],
		InstanceType:       config.Builds[build].Size,
		KeyName:            config.Builds[build].Key,
		SecurityGroups:     config.Builds[build].SecurityGroups,
	}

	if config.Builds[build].Placement != "" {
		options.PlacementGroupName = config.Builds[build].Placement
	}

	if config.Builds[build].UserData != "" {
		userData, err := ioutil.ReadFile(home + "/.clifford.d/" + config.Builds[build].UserData)
		if err != nil {
			log.Fatal("unable to read file")
		}
		options.UserData = userData
	}

	resp, err := conn.RunInstances(&options)
	if err != nil {
		log.Fatal("AWS ec2 Run Instances Fail", err)
	}

	instIds := make([]string, len(resp.Instances))
	for idx, instance := range resp.Instances {
		fmt.Println("Now running", instance.InstanceId)
		instIds[idx] = instance.InstanceId
	}
	_, err = conn.CreateTags(instIds, tags)
	if err != nil {
		fmt.Println("Error Creating Tags: ", err)
	}

	fmt.Println("Make sure you terminate instances to stop the cash flow.")
}
