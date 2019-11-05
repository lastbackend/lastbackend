# Deploy the Last.Backend cluster via Ansible.


## Last.Backend Ansible Playbook

The goal is easily install a Last.Backend cluster components on machines running:

- [X] Debian 
- [x] Ubuntu  
- [x] Raspbian [ test mode ]
- [ ] CentOS

on processor architecture:

- [X] x64
- [X] arm64
- [X] armhf

## System requirements:

Deployment environment must have Ansible 2.4.0+
Master and workers must have passwordless SSH access

## Usage

Add the system information gathered above into a file called hosts.ini. For example:

```text
[master]
10.0.0.1

[worker]
10.0.0.[10:11]

[lastbackend-cluster:children]
master
node
```

Start provisioning of the cluster using the following command:

```bash
ansible-playbook cluster.yml
```

## Connect 

To get access to your **Last.Backend** cluster just execute:

```bash
scp 10.0.0.1:/var/lib/lastbackend/server/access-token
lb -H 10.0.0.1 node list

``` 