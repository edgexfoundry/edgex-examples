# Application Service Examples

## Overview

This folder contains various examples of Application Services based on the App Functions SDK and Profiles that can be used with App Service Configurable. See the [Application Services](https://docs.edgexfoundry.org/1.2/microservices/application/ApplServices/.) documentation for complete details on Application services.

> *Note: All custom examples have their own `go.mod` file, `Makefile` and in some cases `Dockerfile` (e.g.`simple-filter-xml`). This allows for any example to be copied to your own repository and used as a template in creating your own new custom application service.*

## Build Prerequisites

Please see the [edgex-go README](https://github.com/edgexfoundry/edgex-go/blob/master/README.md).

## Building Examples

Each `custom` example has its own `Makefile` with the `build` target. So simply run `make build` from an example's base folder when exploring that example.

The top level `makefile` is designed to build all the examples under the `custom` folder. Thus the Makefile does not need to be updated when a new example is added to the `custom` folder

​	run `make build` to build all examples.

For simplicity, the executable created for each example is named `app-service` and placed in each examples sub-folder.

## Running an Example

After building examples you simply cd to the folder for the example you want to run and run the executable for that example with or without any of the supported command line options.

The following commands will run the `simple-filter-xml` example

```
cd custom/simple-filter-xml
./app-service
```

## Building App Service Docker Image

The  `aws-export`, `cloud-event`, `secrets` and `simple-filter-xml` examples contain an example `Dockerfile` to demonstrate how to build a **Docker Image** for your Application Service. 

The Makefile in each of those example folders also contains the `docker` target which will build the **Docker Image** for the example.

​	run `make docker`

This command must be run from those example folders, i.e. not the top level Makefile.

> *Note that Application Services no longer use docker profiles. They use Environment Overrides in the docker compose file to make the necessary changes to the configuration for running in Docker. See the **Environment Variable Overrides For Docker** section in [App Service Configurable's README](https://github.com/edgexfoundry/app-service-configurable/blob/master/README.md#environment-variable-overrides-for-docker)* for more details and an example. 

## Profiles

The profiles folder contains example profiles for use with App Service Configurable. 