# Contributing to Deploy It

## How to contribute to open source on GitHub

Please, read the following article: https://guides.github.com/activities/contributing-to-open-source/, before contributing.

## Benefits for active contributors

**Unlimited apps deployments to your servers and 20% discount on Deployit servers**

## Coding Style

1. All code should be formatted with `gofmt`
2. All code should follow the guidelines covered in Effective Go and Go Code Review Comments.
3. Comment the code. Tell us the why, the history and the context.
4. Variable name length should not be too long.

Great materials to read:
* [Effective Go](https://golang.org/doc/effective_go.html)
* [The Go Blog](https://blog.golang.org)

## Reporting issues

1. Tell us version of Deploy it
2. Include the steps required to reproduce the problem, if possible

## Maintainers

Don't forget to add yourself to [maintainers list](https://github.com/deployithq/deployit/blob/master/CONTRIBUTING.md) of this project in pull request =)

## Building

```bash
$ git clone git@github.com:deployithq/deployit.git
$ cd deployit
$ make build
```

## Starting

### Daemon

```bash
$ make start-daemon
```

### Other commands

```bash
$ deploy it --debug --host http://localhost:3000 --tag latest
```


