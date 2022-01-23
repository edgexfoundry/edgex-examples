[To README](README.md)

# 2. How to install package required for EdgeX development

The Ubuntu server 21.10 is running on the RPI and we can access via SSH. This chapter introduces some items to do before EdgeX installation. 

<br/>

## 2.1 Hostname and Timezone

The RPI has its hostname as "ubuntu", which is fine if there is only one machine has that name. However, as we run more RPIs, there are more chances to make mistakes if all the RPIs have a same host name. To change the name (but pick a editor you like):
```sh
# Login to the RPI via ssh
$ sudo vi /etc/hostname

# To take effect of the new name
$ sudo reboot # Then, login again after 1~2 minutes
```

Also, the RPI has its timezone as UTC so that scheduled tasks based on local time will not work as expected. To change the timezone:
```sh
# Please pick an area in the list of tzdata
$ sudo dpkg-reconfigure tzdata

# To confirm the new timezone
$ timedatectl
        Local time: Sun 2020-11-01 11:26:43 PST     
    Universal time: Sun 2020-11-01 19:26:43 UTC     
```

<br/>

## 2.3 Update package list

Ubuntu's package management system is based on Debian and Canonical manages Ubuntu's package list. To update the RPI:
```sh
$ sudo apt update
$ sudo apt upgrade
```

<br/>

## 2.4 Install basic packages

Before EdgeX installation, we need to install some basic packages in advance:
```sh
sudo apt install -y \
    jq \
    vim \
    git \
    tmux \
    curl \
    tree \
    make \
    libzmq3-dev \
    gnupg-agent \
    build-essential \
    ca-certificates \
    apt-transport-https \
    software-properties-common
```

<br/>

## 2.5 Install Go SDK and Delve

Go is a programming language used for EdgeX development and Delve is a debugger for Go development. Both will be used to build and develop EdgeX services later in this tutorial. To install Go SDK and Delve:
```sh
# Go v1.15.3 is being installed here because it is the latest stable version as of today but please check it from https://golang.org/dl/

$ cd ~
$ sudo mkdir /usr/local/go # This may exist

# Download the SDK
$ wget https://go.dev/dl/go1.18.5.linux-arm64.tar.gz

# Extract and place the SDK under /usr/local/go
$ sudo bash -c "tar -xf go1.18.5.linux-arm64.tar.gz --strip-components=1 -C /usr/local/go"

# Go needs a directory for its libraries
$ mkdir ~/go

# Config the bashrc for the Go SDK and libraries' path
$ echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> ~/.bashrc
$ echo "export GOROOT=/usr/local/go" >> ~/.bashrc
$ echo "export GOPATH=$HOME/go" >> ~/.bashrc
$ source ~/.bashrc

# To confirm the configurations
$ go version
go version go1.18.5 linux/arm64

# To install Delve
$ go get -u github.com/go-delve/delve/cmd/dlv

# To confirm Delve's version
$ dlv version
Delve Debugger
Version: 1.8.0
Build: $Id: 6a6c9c332d5354ddf1f8a2da3cc477bd18d2be53 $
```

<br/>

## 2.6 Install Docker and Docker-compose

Docker is a containerization platform/tool. EdgeX' core services are conveniently packaged as docker containers so that we can leverage Docker to run EdgeX. To install Docker and Docker-compose:
```sh
# Install Docker
$ sudo apt install -y docker.io

# To confirm the versions installed 
$ docker -v
Docker version 20.10.7, build 20.10.7-0ubuntu5.1
$ docker-compose -v
docker-compose version 1.27.4, build unknown

# Enable and start the Docker daemon
$ sudo systemctl enable docker
$ sudo systemctl start docker

# Add the current user to the Docker group
$ sudo usermod -aG docker ${LOGNAME}

# Reboot to take effect
$ sudo reboot
```

<br/>

So, now the required packages for EdgeX development are ready!

<br/>

---

Next: [How to install EdgeX](30_install_edgex.md)
