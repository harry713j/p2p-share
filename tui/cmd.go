package main

import "flag"

var (
	send    = flag.String("send", "", "Sends a particular file defined in the path")
	recieve = flag.String("recieve", "", "Connect with the sender with the provided code and download the file")
)
