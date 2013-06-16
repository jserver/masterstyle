package main

import (
	"fmt"
	"launchpad.net/goamz/ec2"
	"log"
)

func GetInstances() map[string]ec2.Instance {
	resp, err := conn.Instances(nil, nil)
	if err != nil {
		log.Fatal("AWS ec2 GetInstances Fail", err)
	}

	instances = make(map[string]ec2.Instance)

	for resIdx, reservation := range resp.Reservations {
		for idx, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if tag.Key == "Name" {
					name = tag.Value
				}
			}
			if name == "" {
				name = fmt.Sprintf("instance-%d-%d", resIdx, idx)
			}
			instances[name] = instance
		}
	}

	return instances
}

func Status() {
	instances = GetInstances()
	for name, instance := range instances {
		fmt.Printf("%s [%s] %s - %s\n", name, instance.InstanceId, instance.State.Name, instance.DNSName)
	}
}

func Reboot(args []string) {
	if len(args) != 1 {
		fmt.Println("No instance name given")
	}
	name := args[0]
	instId := instances[name].InstanceId

	_, err := conn.RebootInstances(instId)
	if err != nil {
		fmt.Println("AWS ec2 Reboot Fail", err)
	}
}

func Start(args []string) {
	if len(args) != 1 {
		fmt.Println("No instance name given")
	}
	name := args[0]
	instId := instances[name].InstanceId

	_, err := conn.StartInstances(instId)
	if err != nil {
		fmt.Println("AWS ec2 Start Fail", err)
	}
}

func Stop(args []string) {
	if len(args) != 1 {
		fmt.Println("No instance name given")
	}
	name := args[0]
	instId := instances[name].InstanceId

	_, err := conn.StopInstances(instId)
	if err != nil {
		fmt.Println("AWS ec2 Stop Fail", err)
	}
}

func Terminate(args []string) {
	if len(args) != 1 {
		fmt.Println("No instance name given")
	}
	name := args[0]
	instId := instances[name].InstanceId

	_, err := conn.TerminateInstances([]string{instId})
	if err != nil {
		fmt.Println("AWS ec2 Terminate Fail", err)
	}
}
