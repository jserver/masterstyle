package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goamz/ec2"
)

func Launch(args []string) {
	if len(args) != 1 {
		fmt.Println("no build given!")
		return
	}
	build := args[0]

	line := AskQuestion("Enter Tag Name: ")
	tags := []ec2.Tag{{"Name", line}}

	options := ec2.RunInstances{
		ImageId:        config.Images[config.Builds[build].Image],
		InstanceType:   config.Builds[build].Size,
		KeyName:        config.Builds[build].Key,
		SecurityGroups: config.Builds[build].SecurityGroups,
	}

	if config.Builds[build].Zone != "" {
		options.AvailZone = config.Builds[build].Zone
	}

	if config.Builds[build].UserData != "" {
		userData, err := ioutil.ReadFile(scriptPath + config.Builds[build].UserData)
		if err != nil {
			fmt.Println("Unable to read UserData script", err)
			return
		}
		options.UserData = userData
	}

	resp, err := conn.RunInstances(&options)
	if err != nil {
		fmt.Println("AWS ec2 Run Instances Fail: ", err)
		return
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
