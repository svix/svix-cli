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

The Svix CLI is available on Linux via:

* The [Snap Store](https://snapcraft.io): `snap install svix`
* The [Arch User Repository (AUR)](https://wiki.archlinux.org/title/Arch_User_Repository): `yay -S svix-cli`
* For Ubuntu/Debian: get the `deb` package from [our Github releases page](https://github.com/svix/svix-cli/releases)
* For Fedora/CentOS: get the `rpm` package from [our Github releases page](https://github.com/svix/svix-cli/releases)


### Pre-built executables

You can download and use our pre-built executables directly from [our releases page](https://github.com/svix/svix-cli/releases), and use them as is without having to install anything.

1. Download and extract the `tar.gz` archive for your operating system.
2. Run the `svix` executable from the command line: `./svix help`.

Note: you may need to allow execution by running `chmod +x svix`.


You can also put the binaries anywhere in your `PATH` so you can run the command from anywhere without needing to provide its full path. On macOS or Linux you can achieve this by moving the executable to `/usr/local/bin` or `/usr/bin`.


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

## Using the `listen` command

The `listen` command creates an on-the-fly publicly accessible URL for use when testing webhooks.

The cli then acts as a proxy, forwarding any requests to the given local URL.
This is useful for testing your webhook server locally without having to open a port or
change any nat configuration on your network.

Example:

`svix listen http://localhost:8000/webhook/`

Output:

```sh
Webhook relay is now listening at
https://api.relay.svix.com/api/v1/receive/q1FB7XNKZTO4s0Tzh5BDTZ7_oktf1NBo/

All requests on this endpoint will be forwarded to your local url:
http://localhost:8080/webhook/
```

The above command will return you a unique URL and forward any POST requests it receives
to `http://localhost:8000/webhook/`.

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
| listen          | Forward webhook requests a local url                       |
| import          | Import data from a file to your Svix Organization          |
| export          | Export data from your Svix Organization to a file          |
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
 2) Install [snapcraft](https://snapcraft.io/docs/installing-snapcraft).
 3) Install goreleaser via the steps [here](https://goreleaser.com/install/).
 4) Build current commit via `goreleaser release --snapshot --skip-publish --rm-dist`.

[release-img]: https://img.shields.io/github/v/release/svix/svix-cli
[release]: https://github.com/svix/svix-cli/releases
[golangci-lint-img]: https://github.com/svix/svix-cli/workflows/go-lint/badge.svg
[golangci-lint]: https://github.com/svix/svix-cli/actions?query=workflow%3Ago-lint
[report-card-img]: https://goreportcard.com/badge/github.com/svix/svix-cli
[report-card]: https://goreportcard.com/report/github.com/svix/svix-cli
