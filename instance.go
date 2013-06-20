package main

import (
	"fmt"
	"launchpad.net/goamz/ec2"
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

func GetInstanceName(instId string) string {
	var name string = "N/A"
	for _, instance := range instances {
		if instId == instance.InstanceId {
			for _, tag := range instance.Tags {
				if tag.Key == "Name" {
					name = tag.Value
				}
			}
		}
	}
	return name
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
			inst := NamedInstance{name, instance}
			instances = append(instances, &inst)
		}
	}

	return instances
}

func Tag() {
	instances := GetInstances()
	answers := make([]Answer, len(instances))
	for idx, value := range instances {
		text := fmt.Sprintf("%s [%s] %s", value.Name, value.InstanceId, value.DNSName)
		answers[idx] = Answer{text, value.InstanceId}
	}
	instanceAnswer := AskMultipleChoice("Instance? ", answers)
	instIds := []string{instanceAnswer}

	line := AskQuestion("Enter Tag Name: ")
	tags := []ec2.Tag{{"Name", line}}

	_, err := conn.CreateTags(instIds, tags)
	if err != nil {
		fmt.Println("Unable to Tag Instance", err)
		return
	}
}

func Status() {
	instances = GetInstances()
	header := []string{"Name", "InstId", "State", "Zone", "Type", "Arch", "RootDevice", "DNS"}
	rows := make([]Row, len(instances))
	for idx, instance := range instances {
		rows[idx] = Row{
			instance.Name,
			instance.InstanceId,
			instance.State.Name,
			instance.AvailZone,
			instance.InstanceType,
			instance.Arch,
			instance.RootDevice,
			instance.DNSName,
		}
	}
	table := Table{header, rows}
	PrintTable(&table)
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
