package main

import (
	"errors"
	"fmt"
	"launchpad.net/goamz/ec2"
)

func GetAddresses() ([]ec2.Address, error) {
	resp, err := conn.DescribeAddresses(nil, nil)
	if err != nil {
		return nil, errors.New("Unable to get IP Addresses")
	}
	return resp.Addresses, nil
}

func Addresses() {
	addresses, err := GetAddresses()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addresses {
		fmt.Println(address.PublicIP, address.InstanceId)
	}
}

func Associate(args []string) {
	addresses, err := GetAddresses()
	if err != nil {
		fmt.Println(err)
		return
	}
	answers := make([]Answer, len(addresses))
	for idx, address := range addresses {
		answers[idx] = Answer{address.PublicIP, address.PublicIP}
	}
	addressAnswer := AskMultipleChoice("IP Address? ", answers)

	instanceMap := GetInstances()
	answers = make([]Answer, len(instanceMap))
	idx := 0
	for key, value := range instanceMap {
		answers[idx] = Answer{key, value.InstanceId}
	}
	instanceAnswer := AskMultipleChoice("Instance? ", answers)

	resp, err := conn.AssociateAddress(addressAnswer, instanceAnswer)
	if err != nil {
		fmt.Println("Unable to Associate IP Address", err)
		return
	}
	fmt.Println(resp.Return)
}

func Disassociate(args []string) {
	if len(args) != 1 {
		fmt.Println("Pass in IP Address!")
	}
	address := args[0]

	resp, err := conn.DisassociateAddress(address)
	if err != nil {
		fmt.Println("Unable to Disassociate IP Address", err)
		return
	}
	fmt.Println(resp.Return)
}
