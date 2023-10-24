package main

import (
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/cmd"
)

func main() {
	clients.Init()
	cmd.Execute()
}
