# SvixCLI

![GitHub release (latest by date)](https://img.shields.io/github/v/release/svixhq/svix-cli) [![GolangCI][golangci-lint-img]][golangci-lint] [![Go Report Card][report-card-img]][report-card]

A CLI to interact with the Svix API.

**With the Svix CLI, you can:**

- Interact with the Stix CLI
- Validate Webhook payloads


## Installation

Homebrew Tap & Scoop Bucket Coming soon.

For now, you can download the binary executable and use it straight away without having to install any additional dependencies.
1) Find the latest release, download the tar.gz file for your given operating system and extract it.
2) Inside you'll find the `svix` executable which you can run directly (Note: you may need to allow execution via `chmod +x <PATH_TO_SVIX_EXE>`),

On macOS or Linux, you can move this file to `/usr/local/bin` or `/usr/bin` locations to have it be runnable from anywhere; or place it anywhere and add it to your path (ex. `export PATH=$PATH:<PATH_TO_SVIX_EXE>`) Otherwise, simply `cd` to the folder where you extracted the tar.gz file and run it with ./svix.


## Usage

Installing the Svix CLI provides access to the `svix` command.

```sh
sivx [command]

# Run `svix help` for information about the available commands
svix help

# or add the `--help` flag to any command for a more detailed description and list of flags
svix [command] --help
```


## Quick Start

```sh
# Set your Auth Token temporarily via the SVIX_AUTH_TOKEN environment variable
export SVIX_AUTH_TOKEN=<MY-AUTH-TOKEN>
# or to persistently store your auth token in a config file run
svix config # interactively configure your Svix API credentials

# Create an Application with the name "Demo"
svix application create '{ "name": "demo" }'
# or pipe in some json
echo '{ "name": "demo" }' | svix application create
# or use the convinence cli flags
svix application create --name demo

# List Applications
svix application list --limit 2 --iterator some_iterator 
```


## Commands

The Svix CLI supports the following commands:
| Command         | Description                                             |
| --------------- | ------------------------------------------------------- |
| login           | Interactively configure your Svix API credentials       |
| application     | List, create & modify applications                      |
| authentication  | Manage authentication tasks such getting dashboard urls |
| endpoint        | List, create & modify endpoints                         |
| event-type      | List, create & modify event types                       |
| message         | List & create messages                                  |
| message-attempt | List, lookup & resend message attempts                  |
| verify          | Verify the signature of a webhook message               |
| version         | Get the version of the Svix CLI                         |
| help            | Help about any command                                  |


## Documentation

For a more information, checkout our [API reference](https://docs.svix.com).


### Development

#### Building the current commit

This project uses [goreleaser](https://github.com/goreleaser/goreleaser/)
 1) Install [go](https://golang.org/doc/install).
 2) Install goreleaser via the steps [here](https://goreleaser.com/install/).
 3) Build current commit via `goreleaser release --snapshot --skip-publish`.

[release-img]: https://img.shields.io/github/v/release/svixhq/svix-cli
[golangci-lint-img]: https://github.com/svixhq/svix-cli/workflows/go-lint/badge.svg
[golangci-lint]: https://github.com/svixhq/svix-cli/actions?query=workflow%3Ago-lint
[report-card-img]: https://goreportcard.com/badge/github.com/svixhq/svix-cli
[report-card]: https://goreportcard.com/report/github.com/svixhq/svix-cli