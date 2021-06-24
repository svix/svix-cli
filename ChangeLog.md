# Changelog

## Version 0.11.0
* add `listen` command for testing webhook receivers locally
* bump svix-libs to v0.18.0, allows webhook verifying with prefixed secrets
* initial snap release!

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
