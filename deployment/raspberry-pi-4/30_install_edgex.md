[To README](README.md)

# 3. How to install EdgeX 

The tools to run EdgeX services are ready. The EdgeX stack consists of many docker containers but users don't need to launch one by one since the EdgeX team offers Docker-compose files per each release iteration. 

<br/>

## 3.1 Launch EdgeX core services

The compose files can be found from a repository - [edgexfoundry/edgex-compose](https://github.com/edgexfoundry/edgex-compose). In this chapter, the Jakarta version will be used as it is the latest stable version. To launch containers:
```sh
$ cd ~

# Let's store all repositories in a directory to organize the home directory.
$ mkdir repo
$ cd repo

# Clone the edgex-compose repository
$ git clone https://github.com/edgexfoundry/edgex-compose
$ git fetch --all
$ git checkout jakarta
$ ls
GOVERNANCE.md
LICENSE
Makefile
OWNERS.md
README.md
compose-builder
docker-compose-arm64.yml
docker-compose-no-secty-arm64.yml
docker-compose-no-secty-with-app-sample-arm64.yml
docker-compose-no-secty-with-app-sample.yml
docker-compose-no-secty.yml
docker-compose-portainer.yml
docker-compose-with-app-sample-arm64.yml
docker-compose-with-app-sample.yml
docker-compose.yml
taf


# There are several compose files but we only need one to launch for our purpose. 
# - ARM64 version should be used for RPI. 
# - Security is out of scope in this tutorial. 
# With these criteria, we will use "docker-compose-no-secty-arm64.yml". 

# This command launches the stack but might take couple minutes depends on the network.
$ docker-compose -f docker-compose-no-secty-arm64.yml up -d
...
Creating edgex-ui-go       ... done
Creating edgex-redis       ... done
Creating edgex-core-consul ... done
Creating edgex-support-notifications ... done
Creating edgex-support-scheduler     ... done
Creating edgex-kuiper                ... done
Creating edgex-core-metadata         ... done
Creating edgex-core-command          ... done
Creating edgex-core-data             ... done
Creating edgex-device-rest           ... done
Creating edgex-device-virtual        ... done
Creating edgex-app-rules-engine      ... done
Creating edgex-sys-mgmt-agent        ... done


# Once launching is done, let's check what are up and running. Some columns are removed.
$ docker ps --format "table {{.ID}}\t{{.Image}}\t{{.Status}}"
CONTAINER ID   IMAGE                                               STATUS
7d041dd2cde5   edgexfoundry/sys-mgmt-agent-arm64:2.1.0             Up 4 minutes
d405bf5834f5   edgexfoundry/device-virtual-arm64:2.1.0             Up 4 minutes
4eacedaff009   edgexfoundry/app-service-configurable-arm64:2.1.0   Up 4 minutes
95a08253847e   edgexfoundry/device-rest-arm64:2.1.0                Up 4 minutes
23110c85c80c   edgexfoundry/core-data-arm64:2.1.0                  Up 4 minutes
8b432e47ea9f   edgexfoundry/core-command-arm64:2.1.0               Up 4 minutes
183dec617f70   edgexfoundry/core-metadata-arm64:2.1.0              Up 4 minutes
8d3945944f74   lfedge/ekuiper:1.3.1-alpine                         Up 4 minutes
47e71a310782   edgexfoundry/support-notifications-arm64:2.1.0      Up 4 minutes
404a9ee4f501   edgexfoundry/support-scheduler-arm64:2.1.0          Up 4 minutes
2859afc4b612   consul:1.10.3                                       Up 4 minutes
70d3a90e88a4   redis:6.2.6-alpine                                  Up 4 minutes
284ff5b3fa92   edgexfoundry/edgex-ui-arm64:2.1.0                   Up 4 minutes
```

<br/>

The EdgeX structure diagram clearly shows the purpose of each service:

![EdgeX 2.1 Architecture Diagram](./assets/EdgeX_architecture.png)

There are the core services in the middle. Devices services will talk to the hardwares. Supporting services will inject rules and run actions scheduled. Application services will interact with frontend or external cloud services. All the well designed services are just launched with the one line of command!

<br/>

## 3.2 Test EdgeX services with Curl 

Although the services are launched well, it is worth to test the servcies before writing any custom service.

Curl is a command line tool of *nix systems to transfer data to a given URL and the basic tool to ping EdgeX services:
```
$ curl http://localhost:59880/api/v2/ping
{"apiVersion":"v2","timestamp":"Mon Jan 10 22:45:05 UTC 2022"}

$ curl http://localhost:59881/api/v2/ping
{"apiVersion":"v2","timestamp":"Mon Jan 10 22:45:32 UTC 2022"}

$ curl http://localhost:59882/api/v2/ping
{"apiVersion":"v2","timestamp":"Mon Jan 10 22:45:56 UTC 2022"}
```

Also docker-compose can be used to monitor logs:
```sh
$ docker-compose -f docker-compose-no-secty-arm64.yml logs -f
```

<br/>

## 3.3 (Optional) Monitor EdgeX services with Portainer 

A local web service "Portainer" can be launched to monitor Docker services but it is not a command line tool so that let's use a web browser from the host machine. Before launch it, we need to edit a line in the compose file:
```sh
# Open and edit the "docker-compose-portainer.yml" 
# and update a line under the ports section from -"127.0.0.1:9000:9000" to - "9000:9000" to allow accesses from other machines.
$ vi docker-compose-portainer.yml

# Then launch it with this command.
$ docker-compose -f docker-compose-portainer.yml up -d
```

<br/>

Now we can access to the portainer web service with RPI's IP address and port number 9000 from the host's browser. For the first time access, a new password needs to be set. With the Portainer UI, we can monitor the log and interact with each service. 

![Portainer](./assets/portainer.png)

<br/>

EdgeX stack is up and running so that we can start making our own custom device and app services. 

<br/>

---

Next: [How to develop custom device services](40_custom_device_services.md)
