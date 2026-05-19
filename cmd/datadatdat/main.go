// Package main provides the Datadatdat remote server for connecting d3 CLI to datadatdat-remote-server.
package main

import (
	// Side-effect import: the datadatdat package's init() registers the provider
	// with the remote SDK registry. Keep this as the sole registration site —
	// matches the sibling pattern used by ssh-remote-go and s3-remote-go.
	_ "github.com/datadatdat/datadatdat-remote-go/datadatdat"

	"github.com/datadatdat/remote-sdk-go/remote"
)

func main() {
	// Serve the remote plugin
	remote.Serve("datadatdat")
}
