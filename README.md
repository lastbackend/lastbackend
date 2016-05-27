![alt text](https://deployit.co/images/deployit-logo.png "Logo")

## Deploy It: the command-line toolkit for fast apps deploying

Deploy It is an open-source command-line toolkit and daemon in one application, which allows you to deploy any application to any server.

Deploy It fetches code from request or repo, build it and deploy it to any server. 
Deploy It uses powerful Docker containers, that means that your app will be run anywhere, from your development enviroment on your laptop to any large scale cloud hosting. 
You can run your own deployit instance on your laptop for local development and deployment and connect remote deployit daemon to provide remote deploy on any server you want.

![alt text](https://deployit.co/images/cdn/intro-cli.png "Image")
___

## Building Deploy It

### Prerequisites for CLI:
- Go 1.6 or higher
- Git

### Prerequisites for Daemon:
- Docker
- Go 1.6 or higher
- Git

### Building
```bash
git clone git@github.com:deployithq/deployit.git
cd deployit
go build -o /opt/bin/deploy
```
___

## Current CLI Commands

### Running Daemon
Run `deploy daemon --debug`

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

## Roadmap
- [ ] Deploy app to host [CLI, DAEMON]
- [ ] Deploy app with configs from yaml file [CLI, DAEMON] 
- [ ] Delete app from host [CLI, DAEMON]
- [ ] App logs streaming [CLI, DAEMON]
- [ ] App start/stop/restart [CLI, DAEMON]
- [ ] Deploy app with git url [CLI, DAEMON]
- [ ] Deploy app with hub url (like Docker Hub) [CLI, DAEMON]
- [ ] Add Digital Ocean host [CLI]
- [ ] Delete Digital Ocean host [CLI]
- [ ] Digital Ocean host start/stop/restart [CLI]
- [ ] Deploy to Digital Ocean host [CLI]
- [ ] Deploy scheduling (at 4:00 pm and for 2 hours) [CLI, DAEMON]
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
deploy daemon --debug
```

2. Cloning sources and running `deploy it` command
```bash
git clone https://github.com/<username>/<repo>
cd <repo>
deploy it --host http://localhost:3000 --tag latest
```
