package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goamz/s3"
)

func S3Upload() {
	bytes, err := ioutil.ReadFile("/home/joe/golang/bin/serverstyle")
	if err != nil {
		fmt.Println("Unable to open ServerStyle executable", err)
	}
	err = bucket.Put("serverstyle", bytes, "application/octet-stream", s3.PublicRead)
	if err != nil {
		fmt.Println("Unable to PUT ServerStyle executable to S3", err)
	}
}
