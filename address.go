package main

import (
	"fmt"
)

func Addresses() {
	resp, err := conn.DescribeAddresses()
	if err != nil {
		fmt.Println("Unable to get Allocated IP Addresses", err)
		return
	}
	for _, address := range resp.Addresses {
		fmt.Println(address.PublicIP, address.InstanceId)
	}
}

func Associate(args []string) {
	if len(args) != 2 {
		fmt.Println("Pass in IP Address and InstanceID!")
	}
	address := args[0]
	instanceId := args[1]

	resp, err := conn.AssociateAddress(address, instanceId)
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
