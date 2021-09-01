# Changelog

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
