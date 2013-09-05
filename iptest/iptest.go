package main

import (
	"fmt"

	"github.com/daaku/go.ip"
	"os"
)

func main() {
	c := ip.Config{
		Device:  "eth0",
		Addr:    "192.168.100.4",
		Mask:    "255.255.255.0",
		Gateway: "192.168.100.1",
	}

	if err := c.Activate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
