# deployit: the command-line toolkit for apps deploying #

DeployIT is an open-source command-line toolkit to deploy apps and manage infrastructure.

DeployIT fetch code from request or repo, build it and deploy to any server connected. 
DeployIT uses powerful Docker containers, that means that your app will be run anywhere, from your development enviroment on your laptop to any large scale cloud hosting. 
You can run your own deployit instance on your laptop for local development and deployment and connect remote deployit daemon to provide remote deploy on any server you want.

-------------------------------------------------------------------------------

# deployit: structure #

The deployit toolkit contains several components its based on:
- cli:   command line interface for communication from terminal
- app:   the main package, that handle all application management logic
- repo:  works with git repos for implementing build from repo and autobuild features
- node:  node management runtime, with connected cloud integrations
- build: builder package runtime for manage build application process

-------------------------------------------------------------------------------

# deployit: daemon #

Deployit daemon is daemon tool manages and provide all necessary logic with your apps and servers. It was created and adopted to teams and support many cool features like:
- users and team management
- nodes and clusters management
- build on special prepared servers
- apps backups and rollbacks
- slack and other integrations
- autodeploy after code changes

Deployit will be use remote platform and resources by default after you'll be logged in.

-------------------------------------------------------------------------------

# deployit: hello world #

That is no way to do it simpler.

```bash
$ deploy it
```

What magic is behind this 2 words:
1. deployit daemon scans files in app directory and store information about sources in .deployit folder
2. after scanning deamon create tar.gz arhieve and send in to remote daemon
3. remote daemon store arhieve to network storage and run build process
4. after successful build daemon export root context and send it to server
5. remote daemon on server extract arhieve and start runC container in provided archieve

-------------------------------------------------------------------------------

# deployit: api #

This section shows some basic deployit cli operations you can do.


## getting help ##

receive information about available operations with deployit command line tool   

```bash
$ deploy help 
```

-------------------------------------------------------------------------------

## deploy app ##

you can deploy current directory sources or remote git repository from github, bitbucket or gitlab. There is no sence what deployit will deploy.
to deploy current directory use ```it``` as argument

```bash
$ deploy it  
```

After executing this command, deploy it deamon will create a special container on prepared server and deploy sources to it. You can also select specific cloud provider to use for deploy.
You can get list of supported remote cloud executing this command 

```bash
$ deploy to help
$ deploy it to digitalocean --args
```
For first droplet creating it will ask you for an cloud provider access. You can open link to use oauth2 authentication or just enter digitalocean access token.
If you already have connected server to deployit you can set it for next deploy by setting its hostname like:

```bash
$ deploy it to droplet-00
```

-------------------------------------------------------------------------------

## manage app ##

You can start, stop, restart, resize, remove and receive logs from running app executing same name commands.

-------------------------------------------------------------------------------

## manage services ##

Deployit cli can run any service for you. Just type 

```bash
$ deploy services list 
```

to see available for deploy services. 
To create a new one service just type

```bash
$ deploy service redis 
```

and after few minutes you'll receive dns and port for running service. Like app deploy you can set specific cloud or server for deploy or create new one.
