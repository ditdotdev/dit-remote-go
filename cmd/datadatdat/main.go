// Package main provides the Datadatdat remote server for connecting d3 CLI to datadatdat-remote-server.
package main

import (
	"github.com/datadatdat/datadatdat-remote-go/datadatdat"
	"github.com/datadatdat/remote-sdk-go/remote"
)

func main() {
	// Register the datadatdat provider
	remote.Register(datadatdat.NewProvider())

	// Serve the remote plugin
	remote.Serve("datadatdat")
}
