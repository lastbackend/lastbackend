[![Go Report Card](https://goreportcard.com/badge/github.com/deployithq/deployit)](https://goreportcard.com/report/github.com/deployithq/deployit)
[![GoDoc](https://godoc.org/github.com/deployithq/deployit?status.png)](https://godoc.org/github.com/deployithq/deployit)
[![Travis](https://travis-ci.org/deployithq/deployit.svg?branch=master)](https://travis-ci.org/deployithq/deployit)
[![Join the chat at freenode:deployit](https://img.shields.io/badge/irc-freenode%3A%20%23deployithq--dev-blue.svg)](http://webchat.freenode.net/?channels=%23deployit)
[![Licensed under Apache License version 2.0](https://img.shields.io/github/license/deployithq/deployit.svg?maxAge=2592000)](https://www.apache.org/licenses/LICENSE-2.0)

![alt text](https://deployit.io/images/cdn/logo-purpure.png)

## Deploy It: the command-line toolkit for fast apps deploying

Deploy It is an open-source command-line toolkit and daemon in one application, which allows you to deploy applications to server.

Deploy It fetches code from current directory, request or repo, build it and deploy it to server. 
Deploy It uses powerful containers, that means that your app will be run anywhere, from your development environment on your laptop to any large scale cloud hosting. 
You can run deploy it daemon on the host where you want to deploy your applications (it can be local or remote), run CLI with this host and your applications will be deployed on specified host.  

This project has [Roadmap](https://github.com/deployithq/deployit/blob/master/ROADMAP.md), feel free to offer your features. 

Look at our [Changelog](https://github.com/deployithq/deployit/blob/master/CHANGELOG.md) to see project progress!

We are actively searching for contributors! If you want to help our project and to make developers life easier, please read our **[Contibuting guideliness](https://github.com/deployithq/deployit/blob/master/CONTRIBUTING.md)**.

**We have benefits for active contributors!**

![alt text](https://deployit.io/images/cdn/deployy_2.gif "Image")

___

## Table of contents

1. [Key features](#key_features)
2. [Prerequisites](#prerequisites)
3. [How to install](#how_to_install)
4. [Current CLI commands](#current_cli_commands)
5. [Help](#help)
6. [Maintainers](#maintainers)

___

## <a name="key_features"></a>Key features
1. Fast application deploying to any server
2. Easy application sharing
3. Easy application management
4. Deploying application with url/hub (like docker hub)
5. Deploying scheduling
6. Deploying services like redis, rabbitmq, mysql, etc.

___

## <a name="prerequisites"></a>Prerequisites

### CLI:
- Go 1.6 or higher
- Git

### Daemon:
- Docker
- Go 1.6 or higher
- Git

___

## <a name="how_to_install"></a>How to install

1. Download Deploy it
```bash
$ git clone git@github.com:deployithq/deployit.git
$ cd deployit
$ make build
```

2. Start daemon on host, where you want to deploy your apps
```bash
$ sudo deploy daemon
```

3. Clone sources and run `$ deploy it` command while in sources directory
```bash
$ git clone https://github.com/<username>/<repo>
$ cd <repo>
$ deploy it --host localhost --port 3000 --tag latest
```

___

## <a name="current_cli_commands"></a>Current CLI commands

### Daemon

Install daemon on the host, where you want to deploy your apps

Run `$ sudo deploy daemon`

Daemon flags:
* [--debug] Shows you debug logs
* [--port] Port, which daemon will listen
* [--docker-uri] Docker daemon adress
* [--docker-cert] Docker client certificate
* [--docker-ca] Docker certificate authority that signed the registry certificate
* [--docker-key] Docker client key


### It:

1. Go to folder with your application source code
2. Run `$ deploy it --host localhost --port 3000 --tag latest --log`

What magic is behind `$ deploy it` command:

1. CLI scans all files
2. CLI creates hash table for scanned files
3. CLI packs needed files into tar.gz
4. CLI sends all files to daemon via HTTP
5. DAEMON unpacks tar.gz
6. DAEMON builds unpacked sources
7. DAEMON deploys app to host where daemon is running

## Deploy config

If you want to deploy your application with specific configurations, you can create "deployit.yaml" file, as shown below:

```
env: 
- DEBUG=*
- HOST=localhost
- PORT=3003
memory: 256
ports: 
- 3000
- 9000
volumes:
- /data:/data
- /opt:/opt
```

Configs:
- env: Environments for your application
- memory: Memory limit
- ports: App ports
- volumes: Host storage : App storage

This config is optional. Use it only if you want.

### App start/stop/restart/remove

1. Go to folder with your application source code
2. Run `$ deploy app --host localhost --port 3000 start`

### Common flags

These flags are suitable for all commands except daemon.

Deploy it flags:
* [--debug] Shows you debug logs
* [--tag] Version of your app, examples: "latest", "master", "0.3", "1.9.9", etc.
* [--host] Adress of your host, where daemon is running
* [--port] Port of daemon host
* [--ssl] HTTPS mode if your daemon uses ssl
* [--log] Show build logs

### Future commands

* deploy git
* deploy hub
* deploy app logs
* deploy it at 4:00 pm for 2 hours
* deploy redis/mysql/mongodb/rabbitmq ...

___

## <a name="help"></a>Help

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
```

### Help about other commands:
```bash
$ deploy <command> --help
```

___

## <a name="maintainers"></a>Maintainers

We have separated maintainers page here: [MAINTAINERS.md](https://github.com/deployithq/deployit/blob/master/MAINTAINERS.md)

### Authors

Alexander: https://github.com/undassa

Konstantin: https://github.com/unloop

Bogdan: https://github.com/gofort
