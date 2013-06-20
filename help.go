package main

import "fmt"

func help() {
	fmt.Println(`The following commands are available
	exit,q,quit         Leave Program
	
	ls, status          List all machines
	launch              Create a new ec2 instance from specified build
	update              apt-get update
	upgrade             apt-get upgrade
	install             Install a bundle from json config
	script              Run a script on the machine specified
	test                Pass in a few directories to see if ServerStyle is responding
	easy, easy_install  Install a python bundle from json config

	EC2 Actions on a specified Machine
	----------------------------------
	reboot              Reboot specified machine
	start               Start specified machine
	stop                Stop specified machine
	terminate           Terminate specified mach

	S3 Actions
	----------
	s3upload            Upload ServerStyle to s3 for new machines

	IP Addresses
	------------
	addresses           List all IP Addresses allocated
	associate           Attach an instance to one of your allocated addresses
	disassociate        Release an address from an instance
	`)
}


