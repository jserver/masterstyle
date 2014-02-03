package main

import (
	"errors"
    "fmt"
    "log"
    "strings"

    "launchpad.net/goamz/ec2"
)

func GetInstance(name string) (*NamedInstance, error) {
    for _, instance := range instances {
        if name == instance.Name {
            return instance, nil
        }
    }
    return nil, errors.New("Instance not found")
}

func GetInstanceName(instId string) string {
    var name string
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
    header := []string{"Name", "InstId", "State", "Zone", "Type", "bit", "ebs", "DNS"}
    rows := make([]Row, len(instances))
    for idx, instance := range instances {
        tf := map[bool]string{true: "Y", false: "N"}
        bits := map[string]string{"i386": "32", "x86_64": "64"}
        rows[idx] = Row{
            instance.Name,
            instance.InstanceId,
            instance.State.Name,
            instance.AvailZone,
            strings.Split(instance.InstanceType, ".")[1],
            bits[instance.Arch],
            tf[instance.RootDevice == "ebs"],
            instance.DNSName,
        }
    }
    table := Table{header, rows}
    PrintTable(&table)
}

func Reboot(args []string) {
    if len(args) != 1 {
        fmt.Println("Usage: reboot <name>")
        return
    }
    name := args[0]
    instance, err := GetInstance(name)
    if err != nil {
        fmt.Println(err)
        return
    }
    instId := instance.InstanceId

    _, err = conn.RebootInstances(instId)
    if err != nil {
        fmt.Println("AWS ec2 Reboot Fail", err)
    }
}

func Start(args []string) {
    if len(args) != 1 {
        fmt.Println("Usage: start <name>")
        return
    }
    name := args[0]
    instance, err := GetInstance(name)
    if err != nil {
        fmt.Println(err)
        return
    }
    instId := instance.InstanceId

    _, err = conn.StartInstances(instId)
    if err != nil {
        fmt.Println("AWS ec2 Start Fail", err)
    }
}

func Stop(args []string) {
    if len(args) != 1 {
        fmt.Println("Usage: stop <name>")
        return
    }
    name := args[0]
    instance, err := GetInstance(name)
    if err != nil {
        fmt.Println(err)
        return
    }
    instId := instance.InstanceId

    _, err = conn.StopInstances(instId)
    if err != nil {
        fmt.Println("AWS ec2 Stop Fail", err)
    }
}

func Terminate(args []string) {
    if len(args) != 1 {
        fmt.Println("Usage: terminate <name>")
        return
    }
    name := args[0]
    instance, err := GetInstance(name)
    if err != nil {
        fmt.Println(err)
        return
    }
    instId := instance.InstanceId

    _, err = conn.TerminateInstances([]string{instId})
    if err != nil {
        fmt.Println("AWS ec2 Terminate Fail", err)
    }
}

func CreateImage(args []string) {
    if len(args) != 1 {
        fmt.Println("create_image <name>")
        return
    }
    name := args[0]
    instance, err := GetInstance(name)
    if err != nil {
        fmt.Println(err)
        return
    }
    instId := instance.InstanceId
    amiName := AskQuestion("Enter Name: ")
    amiDesc := AskQuestion("Enter Description: ")
    image, err := conn.CreateImage(instId, amiName, amiDesc)
    if err != nil {
        fmt.Println("AWS ec2 create image Fail", err)
    }
    fmt.Println("Newly created image: ", image.ImageId)
}
