[![Go Report Card](https://goreportcard.com/badge/github.com/lastbackend/lastbackend)](https://goreportcard.com/report/github.com/lastbackend/lastbackend)
[![GoDoc](https://godoc.org/github.com/lastbackend/lastbackend?status.png)](https://godoc.org/github.com/lastbackend/lastbackend)
[![Travis](https://travis-ci.org/lastbackend/lastbackend.svg?branch=master)](https://travis-ci.org/lastbackend/lastbackend)
[![Gitter](https://badges.gitter.im/lastbackend/lastbackend.svg)](https://gitter.im/lastbackend/lastbackend?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Licensed under Apache License version 2.0](https://img.shields.io/github/license/lastbackend/lastbackend.svg?maxAge=2592000)](https://www.apache.org/licenses/LICENSE-2.0)
[![StackShare](https://img.shields.io/badge/tech-stack-0690fa.svg?style=flat)](https://stackshare.io/last-backend/last-backend)

![alt text](docs/assets/preview.png?raw=true "Image")


### The Open-source container management platform

Last.Backend container management platform is the new and modern open-source container management system with service discovery, overlay networks and more..

## Table of contents

1. [About](#about-container-platform)
2. [Key Features](#a-namekey_featuresakey-features)
3. [Runtime design principles](#runtime-design-principles)
4. [Endpoint Interface](#endpoint-interface)
5. [Current State](#current-state)
6. [How to get started](#how-to-get-started)
7. [Maintainers](#maintainers)
8. [Roadmap](#roadmap)
9. [Community](#community)
10. [Authors](#repository-owners)
11. [License](#license)

___

## About container platform

Container platform is an open-source system for automating deployment, scaling, and management of containerized applications. 
It contains a set of technologies like: container scheduling, service discovery, overlay network, container runtime, 
container images runtime, load-balancing and more. All of these technologies provided out of the box and don't requiring additional research and setup.

![alt text](docs/assets/structure.png?raw=true "Image")


Platform structure:

- REST API server
- State controller
- Resources Scheduler
- Node agent
- Service discovery
- Ingress proxy
- Log collector
- IPVS-based load-balancer
- VxLAN-based overlay network
- Storage runtime driver
- Container runtime
- Container images runtime

___

## <a name="key_features"></a>Key features

1. Fast application deploying to any server
2. Easy application sharing
3. Easy application management
4. Deploying application by url
5. Deploying scheduling
6. Deploying stateful services.
7. Developer-friendly CLI


## Runtime design principles
Our design principles allows us to create extendable and powerful system. We separated runtime into particular package and use interfaces to add ability to extend supported technologies.
By default Last.Backend operate with this runtimes:
- CRI - container runtime interface: docker by default
- CII - container image interface: docker by default
- CSI - container storage interface: host directory by default
- CNI - container network interface: vxlan by default
- CPI - container proxy interface: IPVS by default

All these runtimes are documented in runtime section, where are described all methods, types and algorithms.

## Endpoint interface
The main endpoint to manage cluster is REST API interface.
Our team use swagger for generation API documentation. To create swagger spec, just execute ``` make swagger-spec``` command in root of repository.


You can use REST API in these options:

- directly with CURL or another apps
- using Last.Backend CLI (located in separate repo [lastbackend/cli](https://github.com/lastbackend/cli))
- for building custom go apps - you can use golang client located in `pgk/api/client` package

## Current state

Current version is in public beta stage and include:
- cluster management
- node management
- overlay network based on vxlan
- internal endpoints for pods balancing based on IPVS
- ingress server based on haproxy
- internal discovery server
- services management with basic blue/green deployments
- volumes management

All of these functionality is under active test now, so don't surprised by frequent PRs please.

Join us in Gitter [![Gitter](https://badges.gitter.im/lastbackend/lastbackend.svg)](https://gitter.im/lastbackend/lastbackend?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
This project has [Roadmap](ROADMAP.md), feel free to offer your features. 

We are actively searching for contributors! If you want to help our project and to make developers life easier, please read our **[Contibuting guideliness](http://docs.lastbackend.com/#_contributing)**.



## How to get started

If you want to dive into project, the best place to start - is our **[documentation](http://docs.lastbackend.com/#_getting_started)**.


## Maintainers

We have separated [maintainers page](https://github.com/lastbackend/lastbackend/blob/master/MAINTAINERS.md)

**Note**: Last.Backend is under active development stage and our team is working day and night to make it better.
Your suggestions, comments and contributions will be very helpful for us!


## Roadmap

For details on our planned features and future direction please refer to our [roadmap](ROADMAP.md).


### Community

Join us on social media:
 - [Twitter](https://twitter.com/LastBackend)
 - [Facebook](https://www.facebook.com/lastbackend)
 - [Stackshare](https://stackshare.io/last-backend/last-backend)
 - [AngelList](https://angel.co/last-backend)
 - [LinkedIn](https://www.linkedin.com/company/last-backend)
 - [question@lastbackend.com](mailto:question@lastbackend.com)


### Repository owners

The repository is owned by Last.Backend LLC company.

The main contributors are:
- Alexander: https://github.com/undassa
- Konstantin: https://github.com/unloop

---

## License

Origin is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/).
