package main

import (
	"fmt"
	"log"
)

func GetInstance(name string) *NamedInstance {
	for _, instance := range instances {
		if name == instance.Name {
			return instance
		}
	}
	return nil
}

func GetInstances() []*NamedInstance {
	resp, err := conn.Instances(nil, nil)
	if err != nil {
		log.Fatal("AWS ec2 GetInstances Fail", err)
	}

	instances = []*NamedInstance{}

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
			inst := NamedInstance{name, &instance}
			instances = append(instances, &inst)
		}
	}

	return instances
}

func Status() {
	instances = GetInstances()
	for _, instance := range instances {
		fmt.Printf("%s [%s (%s)] %s - %s\n", instance.Name, instance.InstanceId, instance.AvailZone, instance.State.Name, instance.DNSName)
	}
}

func Reboot(args []string) {
	if len(args) != 1 {
		fmt.Println("No instance name given")
	}
	name := args[0]
	instance := GetInstance(name)
	instId := instance.InstanceId

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
	instance := GetInstance(name)
	instId := instance.InstanceId

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
	instance := GetInstance(name)
	instId := instance.InstanceId

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
	instance := GetInstance(name)
	instId := instance.InstanceId

	_, err := conn.TerminateInstances([]string{instId})
	if err != nil {
		fmt.Println("AWS ec2 Terminate Fail", err)
	}
}
