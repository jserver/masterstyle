package main

import "fmt"

func help() {
	fmt.Println(`The following commands are available
	exit,q,quit            Leave Program
	
	ls, status             List all machines
	security               List all security groups
	launch <build>         Create a new ec2 instance from specified build
	update <name>          apt-get update
	upgrade <name>         apt-get upgrade
	ppa <name> <ppa>       Add a new ppa repo and install package
	script <name> <script> Run a script on the machine specified
	test <dirs...>         Pass in a few directories to see if ServerStyle is responding

	easy_install <name> <bundle>           Install a python bundle from json config

	install <name> package(s) <packges...> Install one or more packages
	install <name> bundle(s) <bundles...>  Install bundles from json config
	install <name> group(s) <groups...>    Install groups from json config

	EC2 Actions on a specified Machine
	----------------------------------
	reboot              Reboot specified machine
	start               Start specified machine
	stop                Stop specified machine
	terminate           Terminate specified mach
	create_image        Create an AMI off the instance

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
