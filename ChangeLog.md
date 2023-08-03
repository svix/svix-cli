# Changelog

## Version 0.20.0
* Update Svix lib (includes fixes and improvements)

## Version 0.19.0
* Listen: attempt reconnecting on any kind of error, not just simple connection issues

## Version 0.18.3
* Further improvements to reconnection and error handling

## Version 0.18.2
* Fix bug in Svix Play causing reconnection issues (revert the change from 0.17.1)

## Version 0.18.1
* Snapcraft: fix releasing to snapcraft

## Version 0.18.0
* Update Svix libs (expose new properties)
* Add new app-portal endpoint and deprecate dashboard-access

## Version 0.17.2
* CI: fix build on Windows

## Version 0.17.1
* Fix bug in Svix Play causing reconnection issues

## Version 0.17.0
* Login: Support setting API url
* Bump svix-libs to 0.52.0 & support new features

## Version 0.16.0
* Add bindings for integrations API

## Version 0.15.1
* Verify: improved error message for invalid signatures
* Improved detection of stdin data when piping values to commands

## Version 0.15.0
* Make dashboard access command name uniform with lib names
* Verify: improve error output & flag descriptions
* Bump cobra to 1.3.0, replace completion command
* Pretty: disable color by default on windows
* Bump svix-libs to 0.41.1 & support new features
* Verify: Improve error message for invalid timestamp

## Version 0.14.0
* Add support for non-POST http methods on play

## Version 0.13.0
* Add play.svix.com support

## Version 0.12.0
* Update `--color` flag to enum `auto|always|never`
* Add `authentication logout` command for invalidating dashboard tokens
* Add `import` and `export` commands for quickly adding event-types
* Bump svix-libs t0 v0.24.0 (adds rate limit support, and disabled value to endpoints)
  
## Version 0.11.1
* Fix bug where application update expected the wrong number of arguments

## Version 0.11.0
* Add `listen` command for testing webhook receivers locally
* Bump svix-libs to v0.18.0, allows webhook verifying with prefixed secrets
* Don't html escape when printing json output, fixes dashboard urls
* Initial snap release!

## Version 0.10.0
* print response body to stdin on api errors

## Version 0.9.0
* adds `open` command for quickly opening Svix's documentation in your browser (thanks @codepope!)

## Version 0.8.1
* remove erroneous printing to stderr
* official release to brew and scoop

## Version 0.8.0
* add `completion` function to generate shell completion scripts

## Version 0.7.0
* enable colorized output by default if stdout is a TTY

## Version 0.6.1
* fix bug in `endpoint create` which caused a panic

## Version 0.6.0
* Change verify to use flags for non payload data instead of positional arguments

## Version 0.5.0
* Initial Public Release :rocket:

## Version 0.0.0 (Initial release)
* Allows interaction with the Svix API
* Allows validation of webhook payloads
