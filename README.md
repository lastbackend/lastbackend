![alt text](https://deployit.co/images/deployit-logo.png "Logo")

## Deploy It: the command-line toolkit for fast apps deploying

Deploy It is an open-source command-line toolkit and daemon in one application, which allows you to deploy any application to any server.

Deploy It fetches code from request or repo, build it and deploy it to any server. 
Deploy It uses powerful Docker containers, that means that your app will be run anywhere, from your development enviroment on your laptop to any large scale cloud hosting. 
You can run your own deployit instance on your laptop for local development and deployment and connect remote deployit daemon to provide remote deploy on any server you want.

[Roadmap](https://github.com/deployithq/deployit/blob/master/ROADMAP.md)

[Changelog](https://github.com/deployithq/deployit/blob/master/CHANGELOG.md)

**Contibuting**: https://github.com/deployithq/deployit/blob/master/CONTRIBUTING.md

**We have benefits for active contributors!**

![alt text](https://deployit.co/images/cdn/deployy_2.gif "Image")

___

## Table of contents

1. [Key features](#key_features)
2. [Prerequisites](#prerequisites)
3. [Getting started](#getting_started)
4. [Current CLI Commands](#current_cli_commands)
5. [Help](#help)
6. [Maintainers](#maintainers)

___

## <a name="key_features"></a>Key features
1. Fast application deploying to any server
2. Easy application sharing
3. Managing servers
4. Deploying application with url/hub (like docker hub)
5. Deploying scheduling

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

## <a name="getting_started"></a>Getting started

1. Install Deploy it

2. Start daemon
```bash
deploy daemon
```

3. Clone sources and runn `deploy it` command while in sources directory
```bash
git clone https://github.com/<username>/<repo>
cd <repo>
deploy it --host http://localhost:3000 --tag latest
```

___

## <a name="current_cli_commands"></a>Current CLI Commands

### Running Daemon
Run `deploy daemon`

### Deploy it:

1. Go to folder with your application source code
2. Run `deploy it --debug --host http://localhost:3000 --tag latest`

What magic is behind `deploy it` command:

1. CLI scans all files
2. CLI creates hash table for scanned files
3. CLI packs needed files into tar.gz
4. CLI sends all files to daemon via HTTP
5. DAEMON unpacks tar.gz
6. DAEMON builds unpacked sources
7. DAEMON deploys app to host where daemon is running

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
