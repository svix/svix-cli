# Svix CLI

[![GitHub release (latest by date)][release-img]][release]
[![GolangCI][golangci-lint-img]][golangci-lint]
[![Go Report Card][report-card-img]][report-card]

A CLI to interact with the Svix API.

**With the Svix CLI, you can:**

- Interact with the Svix CLI
- Validate Webhook payloads


## Installation

### macOS

The Svix CLI is available on macOS via [Homebrew](https://brew.sh/):

```sh
brew install svix/svix/svix
```

### Windows

The Svix CLI is available on Windows via the [Scoop](https://scoop.sh/) package manager:

```sh
scoop bucket add svix https://github.com/svix/scoop-svix.git
scoop install svix
```

### Linux
Via snap
```bash
snap install svix-cli
```
The latest binary release
```bash
bash <(curl -sL https://github.com/svix/svix-cli/releases/download/latest/svix.tar.gz | tar xz && mv svix-cli /usr/bin/)
```

### From source
Install
```bash
git clone https://github.com/svix-cli/svix-cli.git
cd svix-cli
```
Build
```
make
```
For more information on building
```bash
make help
```
### Other
If you are on another platform such as linux or just rather not use a package manager, you can download a binary from our Github releases and use it straight away without having to install any additional dependencies.
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

curl -sL  | tar zx && sudo mv ./light-ca /usr/bin/light-ca

## Quick Start

```sh
# Set your Auth Token temporarily via the SVIX_AUTH_TOKEN environment variable
export SVIX_AUTH_TOKEN=<MY-AUTH-TOKEN>
# or to persistently store your auth token in a config file run
svix login # interactively configure your Svix API credentials

# Create an Application with the name "Demo"
svix application create '{ "name": "demo" }'
# or pipe in some json
echo '{ "name": "demo" }' | svix application create
# or use the convinence cli flags
svix application create --data-name demo

# List Applications
svix application list --limit 2 --iterator some_iterator 
```


## Commands

The Svix CLI supports the following commands:
| Command         | Description                                                |
| --------------- | ---------------------------------------------------------- |
| login           | Interactively configure your Svix API credentials          |
| application     | List, create & modify applications                         |
| authentication  | Manage authentication tasks such as getting dashboard URLs |
| endpoint        | List, create & modify endpoints                            |
| event-type      | List, create & modify event types                          |
| message         | List & create messages                                     |
| message-attempt | List, lookup & resend message attempts                     |
| verify          | Verify the signature of a webhook message                  |
| open            | Quickly open Svix pages in your browser                    |
| completion      | Generate completion script                                 |
| version         | Get the version of the Svix CLI                            |
| help            | Help about any command                                     |


## Shell Completions

Shell completion scripts are provided for Bash, Zsh, fish, & PowerShell.

To generate a script for your shell type `svix completion <SHELL NAME>`.

For detailed instructions on configuring completions for your shell run `svix completion --help`.


## Documentation

For a more information, checkout our [API reference](https://docs.svix.com).


### Development

#### Building the current commit

This project uses [goreleaser](https://github.com/goreleaser/goreleaser/).
 1) Install [go](https://golang.org/doc/install).
 2) Install goreleaser via the steps [here](https://goreleaser.com/install/).
 3) Build current commit via `goreleaser release --snapshot --skip-publish --rm-dist`.

Alternatively, use make:
```
make
```
For more info:
```
make help
```

[release-img]: https://img.shields.io/github/v/release/svix/svix-cli
[release]: https://github.com/svix/svix-cli/releases
[golangci-lint-img]: https://github.com/svix/svix-cli/workflows/go-lint/badge.svg
[golangci-lint]: https://github.com/svix/svix-cli/actions?query=workflow%3Ago-lint
[report-card-img]: https://goreportcard.com/badge/github.com/svix/svix-cli
[report-card]: https://goreportcard.com/report/github.com/svix/svix-cli
