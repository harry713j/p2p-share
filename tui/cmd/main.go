package main

import (
	"flag"
	"fmt"

	"github.com/harry713j/p2p-share/tui/internal/command"
)

func main() {

	flag.Parse()
	fmt.Println(*command.Send)
	fmt.Println(*command.Receive)

}
