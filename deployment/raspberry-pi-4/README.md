# EdgeX on Raspberry Pi

Raspberry Pi (RPI) is a series of well-known single board computers (SBC). The credit card sized computer can run Linux as its operating system so that it can be a useful platform for Internet of Things (IoT) projects. 

Although there are various IoT frameworks, EdgeX framework is a solid candidate for IoT projects because of its architecture, workflow, and productivity. Its microservices architecture is very flexible and thus new features can be easily adapted. The framework is built using modern web patterns from the beginning so that developers can leverage their knowledge about the technologies. 

A couple of months ago, EdgeX foundation released a new version called Geneva (each version has a name from cities). Also Canonnical announced a new Ubuntu release 20.10 Groovy Gorilla and it embraces RPI as one of the major IoT targets. This news made me think about the combination of RPI, Ubuntu, and EdgeX. Running EdgeX on RPI is not something new because Golang and Docker are the cross platform tools used for EdgeX code base but this time the OS is Ubuntu-ARM64.

This short tutorial requires some familiarity with the technologies and tools (vi, curl, sed, and jq) but it is progressivly written and easy to follow.

<br/>

## References

Raspberry Pi:
- https://www.raspberrypi.org/
- https://projects.raspberrypi.org/en/pathways/getting-started-with-raspberry-pi

Ubuntu on ARM:
- https://ubuntu.com/download/iot
- https://ubuntu.com/tutorials/how-to-install-ubuntu-on-your-raspberry-pi

EdgeX: 
- https://www.edgexfoundry.org/
- https://github.com/edgexfoundry
- https://docs.edgexfoundry.org/2.1/walk-through/Ch-Walkthrough/

<br/>

## Prerequisites

- An Ubuntu installed PC (will be called **the host** in this series)
- A Raspberry Pi 4B (4GB or 8GB RAM version is recommended)
- A USB-C power adapter
- A micro SD card and reader
- An ethernet cable

<br/>

## Index

- [How to install Ubuntu on RPI](10_install_ubuntu.md)
- [How to install packages required for EdgeX development](20_install_packages.md)
- [How to install and launch EdgeX](30_install_edgex.md)
- [How to develop custom device and app services](40_custom_device_services.md)

