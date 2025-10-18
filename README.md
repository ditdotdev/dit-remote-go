# Datadatdat Remote (Go)

Datadatdat remote provider for connecting d3 CLI to datadatdat-remote-server. This provider enables d3 to use datadatdat-remote-server as a remote storage backend, providing a "GitHub for d3" experience with web UI, organizations, authentication, and collaboration features.

## Architecture

This remote provider allows d3 CLI users to:
- Push commits to datadatdat-remote-server via HTTP/HTTPS
- Pull and clone repositories from datadatdat-remote-server
- List commits with metadata and tags
- Use web UI for browsing and collaboration

**Comparison with other remotes:**
- `s3-remote-go` → Direct S3 storage (basic, like git+SSH)
- `ssh-remote-go` → Direct SSH storage (basic)
- `datadatdat-remote-go` → datadatdat-remote-server (GitHub-style, with web UI and collaboration)

## Quick Start

### Installation

The datadatdat remote is included with d3 CLI v1.0.0+. No separate installation needed.

### Usage

```bash
# Add datadatdat remote
d3 remote add origin https://data.datadatdat.io/myorg/myrepo

# Push commits
d3 commit -m "Initial commit"
d3 push origin master

# Clone from datadatdat-remote-server
d3 clone https://data.datadatdat.io/myorg/myrepo

# List commits
d3 remote log origin
```

## Features

- **HTTP/HTTPS protocol** - Works with datadatdat-remote-server APIs
- **Organization/Repository structure** - URL format: `https://host/org/repo`
- **Authentication** - Bearer token support (via environment or config)
- **Metadata and tags** - Full support for commit metadata and tagging
- **gRPC plugin** - Loaded dynamically by d3 CLI via Hashicorp go-plugin

## Configuration

Remote URL format:
```
http://localhost:8080/myorg/myrepo   (development)
https://data.datadatdat.io/myorg/myrepo   (production)
```

Authentication (optional for MVP):
```bash
export DATADATDAT_API_KEY="your-api-token"
```

## Development

See [DEVELOPING.md](DEVELOPING.md) for build, test, and release instructions.

## Contributing

This project follows the Datadatdat community best practices:

  * [Contributing](https://github.com/datadatdat/.github/blob/master/CONTRIBUTING.md)
  * [Code of Conduct](https://github.com/datadatdat/.github/blob/master/CODE_OF_CONDUCT.md)
  * [Community Support](https://github.com/datadatdat/.github/blob/master/SUPPORT.md)

It is maintained by the [Datadatdat community maintainers](https://github.com/datadatdat/.github/blob/master/MAINTAINERS.md)

For more information on how it works, and how to build and release new versions,
see the [Development Guidelines](DEVELOPING.md).

