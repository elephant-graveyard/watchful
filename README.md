# watchful /ˈwɒtʃf(ə)l/

[![License](https://img.shields.io/github/license/homeport/watchful.svg)](https://github.com/homeport/watchful/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/homeport/watchful)](https://goreportcard.com/report/github.com/homeport/watchful)
[![Build Status](https://travis-ci.org/homeport/watchful.svg?branch=develop)](https://travis-ci.org/homeport/watchful)
[![GoDoc](https://godoc.org/github.com/homeport/watchful?status.svg)](https://godoc.org/github.com/homeport/watchful)
[![Release](https://img.shields.io/github/release/homeport/watchful.svg)](https://github.com/homeport/watchful/releases/latest)

![watchful](.docs/logo.png?raw=true "Watchful logo - blue hexagon with a telescope next to logo text simply saying watchful")

## Introducing the watchful

A tool to measure the disruption caused by a change to a Cloud Foundry environment. The most obvious use-case would be the roll-out of an update of Cloud Foundry itself. Usually this requires some or all of the internal micro services to restart. The respective setup with means to achieve some form of high availability will step in to make sure an end-user does not notice the software maintenance. As always, you cannot always make sure there is no flicker or lost HTTP request. The main purpose of this tool is to measure the impact of a maintenance and to report the metrics to the operator. This project is highly influenced by the [uptimer tool](https://github.com/cloudfoundry/uptimer) from the Cloud Foundry community.

_This project is work in progress._

---------

## Contributing

We are happy to have other people contributing to the project. If you decide to do that, here's how to:

- get a Go development environment with version 1.12 or greater
- fork the project
- create a new branch
- make your changes
- open a PR.

Git commit messages should be meaningful and follow the rules nicely written down by [Chris Beams](https://chris.beams.io/posts/git-commit/):

> The seven rules of a great Git commit message
>
> 1. Separate subject from body with a blank line
> 1. Limit the subject line to 50 characters
> 1. Capitalize the subject line
> 1. Do not end the subject line with a period
> 1. Use the imperative mood in the subject line
> 1. Wrap the body at 72 characters
> 1. Use the body to explain what and why vs. how

---------

### Running test cases and binaries generation

There are multiple make targets, but running `all` does everything you want in one call.

```sh
make all
```

---------

### Install it on your machine

If you live dangerously, you can install watchful with this simply one line bash command:

```sh
curl -sL https://raw.githubusercontent.com/homeport/watchful/master/scripts/download-latest.sh | bash
```

---------

### Test it with Linux on your macOS system

Best way is to use Docker to spin up a container:

```sh
docker run \
  --interactive \
  --tty \
  --rm \
  --volume $GOPATH/src/github.com/homeport/watchful:/go/src/github.com/homeport/watchful \
  --workdir /go/src/github.com/homeport/watchful \
  golang:1.12 /bin/bash
```

---------

## Commands

### watchful run

`run` is basically the main starting point for watchful and starts the watchful engine.
It does have some flags tho, to allow CLI configuration:

- `-v|--verbose`: As expected watchful comes with a verbose option. Adding this to the command call will result in
watchful logging successful merkhet tests as well as printing detailed error logs on merkhet fails.

- `-w|--terminalWidth <intValue>`: As go may no always be able to pick up the correct terminal width, you can provide it in this
flag. The provided values specifies the amount of characters per line in your terminal.

- `-l|--language <stringValue>`: If you want to use a different app runtime type, you can specify the apps programming language with this tag. The default is `go`. Right now watchful supports the following languages:
  - `go`

- `-c|--config <stringValue>`: If you do not want to provide a file based config, this parameter also allows you
to pass the config content directly to the CLI, removing the need for a physical copy of it on the disk.

---------

## Configuration

Watchful is obviously highly configurable to fit and test the cloud foundry instance as well as possible.
A sample configuration can also be found here: [config-sample.yml](https://github.com/homeport/watchful/blob/master/config-model.yml). The config file must be located under `./config.yaml`. The configuration is generally split into four sub-parts.

### Cloud Foundry Configuration `cf`

Found under the yaml node `cf`, the cloud foundry configuration hosts following values:

- `domain`: The domain node defines the cloud foundry domain end-point, under which routes will be available.
Please provide the domain including the scheme and any necessary ports, eg: `https://sample-cluster.foo.com`

- `api-endpoint`: The API-Endpoint node specifies the cloud foundry endpoint against watchful can execute setup
commands. This will usually be your domain with the prefix `api`, eg: `https://api.sample-cluster.foo.com`

- `skip-ssl-validation`: This boolean value will simply define whether watchful will use SSL validation when
authenticating against the cloud foundry cluster. eg: `true`

- `username`: The username simply specifies the username used to authenticate against the cloud foundry instance.
When connecting with an API-Key, this username could be for example `apikey`

- `password`: The password specifies the passphrase used with the given username to authenticate against the cloud
foundry instance: Eg: `aHR0cHM6Ly9nb28uZ2wvUGpYamR6`

### Task Configuration `tasks`

Found under the yaml node `tasks`, the tasks configuration is a list of task nodes that define the update tasks
watchful will run against the cloud foundry instance. They are the equivalent to the `while` part of the uptimer
config

Each of the task nodes has the following attributes that you can configure:

- `cmd`: This defines the base command used to execute the command. If you are executing a bash script, this will
usually be `/bin/bash`

- `args`: This node is a string list and defines the arguments passed to the command specified in `cmd`.
When passing a bash script to a `/bin/bash` cmd this may look like this:

```yaml
 - -c
 -  |
   #!/bin/bash

   set -euo pipefail

   echo "Hello, World!"
   ping github.com
```  

This configuration would print `Hello World` as well as ping `github.com`.

- `merkhet-whitelist`: This yaml node can be configured for each task and is a string list. If defined as a node, only the merkhets with the given name will be monitoring the cloud foundry for the time of your task running. This may be helpful if you certainly know that the task will disrupt a feature of the cluster, but don't want that to corrupt the merkhet results.

- `merkhet-blacklist`: This yaml node will only be active if the `merkhet-whitelist` has not been configured. If defined every merkhet will be monitoring, except the ones listed on this blacklist. This node is again a list of strings.

### Merkhet Configuration `merkhets`

Found under the yaml node `merkhets` this list of yaml nodes defines the general set of running merkhets. In here you can configure the threshold for merkhets or deactivate them completely.
If you do not configure a merkhets, it will not be used at all by watchful.

A node inside the `merkhets` node can be configured like this:

- `name`: The name simply defines what merkhet implementation your configuration targets. Currently implemented names
are:

  - app-pushability
  - http-availability
  - cf-log-functionality
  - cf-recent-log-functionality

- `threshold`: The threshold defines how many of the merkhet tests are allowed to fail. This threshold can be either provided as a flat number (eg: `10`) or as a percentage (eg: `50 %`)

- `heartbeat`: This yaml node overwrites the default heartbeat of the merkhet.
You **should not modify** this as long as you don't have a valid use case for it as it may mess with the efficiency of watchful. It is of the typ string and needs a valid time duration specifier, eg: `1s`, `500ms` or `1m30s`

### Logger Configuration `logger-config`

Found under the yaml node `logger-config` the loggers used in watchful can be configured as well.

- `time-location`: This string simply specifies the time zone in which the current system time will be rendered, eg: `UTC`

- `show-logger-name`: If this boolean is set to true, the logger name will be printed to the console. This is generally advised to enable as it allows deeper error tracing, but may be disabled in certain situations.

---------

## Git pre-commit hooks

When working with code, it may not always be the best idea to wait for travis to throw an error if your build failed.
To automize your development workflow, it may be a good idea to use git pre-commit hooks.

These little snippets of code are run prior to a commit and can determine whether your commit should be accepted.
In the case of `watchful`, a pre-commit hook could look something like this, calling both `test` and `analysis` make
targets before a commit.

You can install the default pre-commit hook using this command in your watchful root directory:

```sh
cat <<EOS | cat > .git/hooks/pre-commit && chmod a+rx .git/hooks/pre-commit
#!/usr/bin/env bash

set -euo pipefail
make analysis test

EOS
```

## License

Licensed under [MIT License](https://github.com/homeport/watchful/blob/master/LICENSE)
