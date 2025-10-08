package main

import (
	"flag"
	"fmt"

	"github.com/harry713j/p2p-share/tui/internal/command"
	"github.com/harry713j/p2p-share/tui/internal/service"
	"github.com/harry713j/p2p-share/tui/internal/util"
)

func main() {
	flag.Parse()

	filePath := *command.Send
	u := util.NewUtility()
	port := u.GetDynamicPort()
	err := service.SendFile(filePath, fmt.Sprint(port))

	if err != nil {
		fmt.Printf("File send failed: %v\n", err)
		return
	}

	code := *command.Receive
	err = service.ReceiveFile(code)

	if err != nil {
		fmt.Printf("File download failed: %v\n", err)
	}
}
