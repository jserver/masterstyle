package main

import (
	"errors"
	"fmt"
	"github.com/jserver/serverstyle/server"
	"io/ioutil"
	"net/rpc"
	"os"
	"strings"
)

func GetAddress(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("No instance name given")
	}
	name := args[0]
	host := instances[name].DNSName

	return fmt.Sprintf("%s:%d", host, 32168), nil
}

func RemoteCall(address string, cmdArgs interface{}, results server.Results, command string) {
	fmt.Println("Calling: " + address)
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	remoteCall := client.Go(command, cmdArgs, results, nil)
	<-remoteCall.Done
	errText := results.GetErr()
	if len(errText) > 0 {
		fmt.Println(">>> [", errText, "]")
	}
	fmt.Println(results.GetOutput())

}

func Update(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdArgs := &server.AptGetUpdateArgs{}
	results := new(server.AptGetUpdateResults)
	command := "AptGetUpdate.Update"

	RemoteCall(address, cmdArgs, results, command)
}

func Upgrade(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdArgs := &server.AptGetUpgradeArgs{}
	results := new(server.AptGetUpgradeResults)
	command := "AptGetUpgrade.Upgrade"

	RemoteCall(address, cmdArgs, results, command)
}

func Install(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) <= 1 {
		fmt.Println("No bundle given")
		return
	}
	bundle := config.Bundles[args[1]]
	packages := strings.Split(bundle, " ")

	cmdArgs := &server.AptGetInstallArgs{packages}
	results := new(server.AptGetInstallResults)
	command := "AptGetInstall.Install"

	RemoteCall(address, cmdArgs, results, command)
}

func Script(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) <= 1 {
		fmt.Println("No script given")
		return
	}
	script := os.ExpandEnv(args[1])
	content, err := ioutil.ReadFile(script)
	if err != nil {
		fmt.Println("Unable to read script")
		return
	}

	script_parts := strings.Split(script, "/")
	script_name := script_parts[len(script_parts)-1]

	cmdArgs := &server.ScriptArgs{script_name, content}
	results := new(server.ScriptResults)
	command := "Script.Runner"

	RemoteCall(address, cmdArgs, results, command)
}

func Test(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) <= 1 {
		fmt.Println("No directories given")
		return
	}

	cmdArgs := &server.TestArgs{args[1:]}
	results := new(server.TestResults)
	command := "Test.Runner"

	RemoteCall(address, cmdArgs, results, command)
}
