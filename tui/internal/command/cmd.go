package command

import "flag"

var (
	Send    = flag.String("send", "", "Sends a particular file defined in the path")
	Recieve = flag.String("recieve", "", "Connect with the sender with the provided code and download the file")
)
