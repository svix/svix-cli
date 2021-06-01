# SvixCLI

![GitHub release (latest by date)](https://img.shields.io/github/v/release/svixhq/svix-cli)

A CLI to interact with the Svix API.

**With the CLI, you can:**

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

# Add the `--help` flag for information about the available commands
svix [command] --help
```

## Commands

The Svix CLI supports the following commands:
 - init            Interactively configure your Svix API credentials
 - application     List, create & modify applications
 - authentication  Manage authentication tasks such getting dashboard urls
 - endpoint        List, create & modify endpoints
 - event-type      List, create & modify event types
 - message         List & create messages
 - message-attempt List, lookup & resend message attempts
 - version         Get the version of the Svix CLI
 - help            Help about any command

## Documentation

For a more information, checkout our [API reference](https://docs.svix.com)

### Development

#### Building the current commit

This project uses (goreleaser)[https://github.com/goreleaser/goreleaser/].
1) Install goreleaser via the steps (here)[https://goreleaser.com/install/]
2) Build current commit via `goreleaser release --snapshot --skip-publish`