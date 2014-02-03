package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"strings"

	"github.com/jserver/serverstyle/server"
)

func GetAddress(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("No instance name given")
	}
	name := args[0]
	instance, err := GetInstance(name)
	if err != nil {
		return "", err
	}
	host := instance.DNSName

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

	if remoteCall.Error != nil {
		fmt.Println(remoteCall.Error.Error())
		return
	}

	errText := results.GetErr()
	if len(errText) > 0 {
		fmt.Println(">>> [", errText, "]")
	}

	errors := results.GetErrors()
	if len(errors) > 0 {
		fmt.Println("-----ERRORS-----")
		fmt.Println(errors)
		fmt.Println("----------------")
	}
}

func Update(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdArgs := &server.AptUpdateArgs{}
	results := new(server.AptUpdateResults)
	command := "AptUpdate.Update"

	RemoteCall(address, cmdArgs, results, command)
}

func Upgrade(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdArgs := &server.AptUpgradeArgs{}
	results := new(server.AptUpgradeResults)
	command := "AptUpgrade.Upgrade"

	RemoteCall(address, cmdArgs, results, command)
}

func GetGroupPackages(name string) (bundles []string) {

	groups := config.Groups[name]
	for _, group := range groups {
		if group.Type == "group" {
			bundles = append(bundles, GetGroupPackages(group.Value)...)
		} else if group.Type == "bundle" {
			packages := config.Bundles[group.Value]
			bundles = append(bundles, packages)
		} else if group.Type == "package" {
			bundles = append(bundles, group.Value)
		}
	}
	return
}

func PPAInstall(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) < 2 {
		fmt.Println("Usage: ppa <box> <name>")
		return
	}

	ppa := config.PPAs[args[1]]

	cmdArgs := &server.PPAInstallArgs{ppa.Name, ppa.Package}
	results := new(server.PPAInstallResults)
	command := "PPAInstall.AddRepo"

	RemoteCall(address, cmdArgs, results, command)

	fmt.Println(">>> Updating")
	Update([]string{args[0]})
	fmt.Println(">>> Installing")
	Install([]string{args[0], "package", ppa.Package})
}

func Install(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) < 3 {
		fmt.Println("Usage: install <box> package(s)/bundle(s)/group(s) <names...>")
		return
	}

	action := args[1]
	names := args[2:]

	bundles := []string{}

	switch action {
	case "package", "packages":
		packages := strings.Join(names, " ")
		bundles = append(bundles, packages)

	case "bundle", "bundles":
		for _, bundle := range names {
			packages := config.Bundles[bundle]
			bundles = append(bundles, packages)
		}

	case "group", "groups":
		for _, name := range names {
			bundles = append(bundles, GetGroupPackages(name)...)
		}

	default:
		fmt.Println("Install action not recognized")
		return
	}

	for _, packages := range bundles {
		cmdArgs := &server.AptInstallArgs{strings.Split(packages, " ")}
		results := new(server.AptInstallResults)
		command := "AptInstall.Install"

		RemoteCall(address, cmdArgs, results, command)
	}
}

func EasyInstall(args []string) {
	address, err := GetAddress(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) <= 1 {
		fmt.Println("No python bundle given")
		return
	}
	bundle := config.PythonBundles[args[1]]
	packages := strings.Split(bundle, " ")

	cmdArgs := &server.EasyInstallArgs{packages}
	results := new(server.EasyInstallResults)
	command := "EasyInstall.Install"

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

	cmdArgs := &server.TestArgs{server.Auth{"user", "pass123"}, args[1:]}
	results := new(server.TestResults)
	command := "Test.Runner"

	RemoteCall(address, cmdArgs, results, command)

	output := results.GetOutput()
	if len(output) > 0 {
		fmt.Println("-----STDOUT-----")
		fmt.Println(output)
		fmt.Println("----------------")
	}
}
