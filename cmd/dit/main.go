// Copyright Dit 2026
// SPDX-License-Identifier: BUSL-1.1

// Package main provides the Dit remote server for connecting d3 CLI to dit-remote-server.
package main

import (
	// Side-effect import: the dit package's init() registers the provider
	// with the remote SDK registry. Keep this as the sole registration site —
	// matches the sibling pattern used by ssh-remote-go and s3-remote-go.
	_ "github.com/ditdotdev/dit-remote-go/dit"

	"github.com/ditdotdev/remote-sdk-go/remote"
)

func main() {
	// Serve the remote plugin
	remote.Serve("dit")
}
