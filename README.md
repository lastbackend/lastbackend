![alt text](https://deployit.co/images/deployit-logo.png "Logo")

## Deploy It: the command-line toolkit for fast apps deploying

Deploy It is an open-source command-line toolkit and daemon in one application, which allows you to deploy any application to any server.

Deploy It fetches code from request or repo, build it and deploy it to any server. 
Deploy It uses powerful Docker containers, that means that your app will be run anywhere, from your development enviroment on your laptop to any large scale cloud hosting. 
You can run your own deployit instance on your laptop for local development and deployment and connect remote deployit daemon to provide remote deploy on any server you want.

![alt text](https://deployit.co/images/cdn/intro-cli.png "Image")

___

## Prerequisites

### CLI:
- Go 1.6 or higher
- Git

### Daemon:
- Docker
- Go 1.6 or higher
- Git

___

## Installation

Coming soon

___

## Contributing

If you want to contribute, please read our contribution guide here: https://github.com/deployithq/deployit/blob/master/CONTRIBUTING.md

___

## Current CLI Commands

### Running Daemon
Run `deploy daemon`

### Deploy it:

1. Go to folder with your application source code
2. Run `deploy it --debug --host http://localhost:3000 --tag latest`

What magic is behind `deploy it` command:

1. CLI scans all files
2. CLI creates hash table for scanned files
3. CLI packs needed files into tar.gz
2. CLI sends all files to daemon via HTTP
3. DAEMON unpacks tar.gz
4. DAEMON builds unpacked sources
5. DAEMON deploys app to host where daemon is running

Deploy it flags:
* [--debug] Shows you debug logs
* [--tag] Version of your app, examples: "latest", "master", "0.3", "1.9.9", etc.
* [--host] Adress of your host, where daemon is running

### Future commands

* deploy url
* deploy it to digital ocean
* deploy it at 4:00 pm for 2 hours
* deploy redis
* deploy search <service>
* deploy app stop/start/restart

___

## Getting help

All information about Deploy It is available via following commands:

### Brief info about all commands
```bash
$ deploy --help
```

### Deploy it command
```bash
$ deploy it --help
```

### Daemon
```bash
$ deploy daemon --help
```16

### Help about other commands:
```bash
$ deploy <command> --help
```

___

## Examples

### Deploying app from sources

1. Starting daemon
```bash
deploy daemon
```

2. Cloning sources and running `deploy it` command
```bash
git clone https://github.com/<username>/<repo>
cd <repo>
deploy it --host http://localhost:3000 --tag latest
```
