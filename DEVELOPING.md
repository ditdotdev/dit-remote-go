# Project Development

For general information about contributing changes, see the
[Contributor Guidelines](https://github.com/datadatdat/.github/blob/master/CONTRIBUTING.md).

## How it Works

The provider uses the Datadatdat `remote-sdk-go` to provide interfaces for
Datadatdat (d3 CLI) to use. The provider implements the `remote.Remote` interface
and communicates with datadatdat-remote-server via HTTP APIs.

This provider enables d3 to use datadatdat-remote-server as a storage backend,
providing a "GitHub for d3" experience with web UI, organizations, and collaboration.

## Building

Run `go build -v ./...`.

To build the plugin binary:

```bash
go build -o datadatdat ./cmd/datadatdat
```

## Testing

Run `go test -v ./...`.

## Releasing

Push a tag of the form `v<X>.<Y>.<Z>`, and publish the draft release in GitHub.

## Integration with d3 CLI

The d3 CLI loads this provider as a gRPC plugin using Hashicorp's go-plugin.
When you add a remote with `http://` or `https://` scheme, d3 automatically
loads this provider to handle the communication with datadatdat-remote-server.
