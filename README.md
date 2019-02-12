# watchful /ˈwɒtʃf(ə)l/

[![License](https://img.shields.io/github/license/homeport/watchful.svg)](https://github.com/homeport/watchful/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/homeport/watchful)](https://goreportcard.com/report/github.com/homeport/watchful)
[![Build Status](https://travis-ci.org/homeport/watchful.svg?branch=develop)](https://travis-ci.org/homeport/watchful)
[![GoDoc](https://godoc.org/github.com/homeport/watchful?status.svg)](https://godoc.org/github.com/homeport/watchful)
[![Release](https://img.shields.io/github/release/homeport/watchful.svg)](https://github.com/homeport/watchful/releases/latest)

## Introducing the watchful

A tool to measure the disruption caused by a change to a Cloud Foundry environment. The most obvious use-case would be the roll-out of an update of Cloud Foundry itself. Usually this requires some or all of the internal micro services to restart. The respective setup with means to achieve some form of high availability will step in to make sure an end-user does not notice the software maintenance. As always, you cannot always make sure there is no flicker or lost HTTP request. The main purpose of this tool is to measure the impact of a maintenance and to report the metrics to the operator. This project is highly influenced by the [uptimer tool](https://github.com/cloudfoundry/uptimer) from the Cloud Foundry community.

_This project is work in progress._

## Contributing

We are happy to have other people contributing to the project. If you decide to do that, here's how to:

- get a Go development environment with version 1.11 or greater
- fork the project
- create a new branch
- make your changes
- open a PR.

Git commit messages should be meaningful and follow the rules nicely written down by [Chris Beams](https://chris.beams.io/posts/git-commit/):
> The seven rules of a great Git commit message
> 1. Separate subject from body with a blank line
> 1. Limit the subject line to 50 characters
> 1. Capitalize the subject line
> 1. Do not end the subject line with a period
> 1. Use the imperative mood in the subject line
> 1. Wrap the body at 72 characters
> 1. Use the body to explain what and why vs. how

### Running test cases and binaries generation

There are multiple make targets, but running `all` does everything you want in one call.

```sh
make all
```

### Test it with Linux on your macOS system

Best way is to use Docker to spin up a container:

```sh
docker run \
  --interactive \
  --tty \
  --rm \
  --volume $GOPATH/src/github.com/homeport/watchful:/go/src/github.com/homeport/watchful \
  --workdir /go/src/github.com/homeport/watchful \
  golang:1.11 /bin/bash
```

### Git pre-commit hooks
When working with code, it may not always be the best idea to wait for travis to throw an error if your build failed.
To automize your development workflow, it may be a good idea to use git pre-commit hooks. 

These little snippets of code are run prior to a commit and can determine whether your commit should be accepted.
In the case of watchful, a pre-commit hook could look something like this, calling both test and analysis make 
targets before a commit.

You can install the default pre-commit hook using this command in your watchful root directory:
```sh 
cat pre-commit.sh > .git/hooks/pre-commit && chmod u+x .git/hooks/pre-commit
```

## License

Licensed under [MIT License](https://github.com/homeport/watchful/blob/master/LICENSE)
