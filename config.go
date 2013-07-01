package main

import "launchpad.net/goamz/ec2"

type Build struct {
	Size           string
	Image          string
	Key            string
	Zone           string
	SecurityGroups []ec2.SecurityGroup
	UserData       string
	Upgrade        string
	Group          string
	EasyInstall    string
	ScriptName     string
	ScriptAction   string
}

type Project struct {
	Build        string
	UserName     string
	UserFullname string
	UserPassword string
	ScriptName   string
	ScriptAction string
}

type Repository struct {
	KeyServer string
	PublicKey string
	DebLine   string
	Package   string
}

type PPA struct {
	Name, Package string
}

type Group struct {
	Type  string
	Value string
}

type Config struct {
	BucketName      string
	DomainName      string
	KeyPath         string
	PasswordSalt    string
	ScriptPath      string
	Images          map[string]string
	Bundles         map[string]string
	PythonBundles   map[string]string
	Groups          map[string][]Group
	AptRepositories map[string]Repository `json:"Apt:Repositories"`
	PPAs            map[string]PPA 
	Builds          map[string]Build
	Projects        map[string]Project
	DebConfs        map[string][]string
}
