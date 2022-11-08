# Application Service Examples

## Overview

This folder contains various examples of Application Services based on the `Levski 2.3.x` App Functions SDK and Profiles that can be used with App Service Configurable. See the [Application Services](https://docs.edgexfoundry.org/latest/microservices/application/ApplicationServices/) documentation for complete details on Application services.

> *Note: All custom examples have their own `go.mod` file, `Makefile` and in some cases `Dockerfile` (e.g.`simple-filter-xml`). * These example are **not** designed to be templates for creating a new Application Service.  For this we have the [app-service-template in app-functions-sdk-go](https://github.com/edgexfoundry/app-functions-sdk-go/blob/main/app-service-template/README.md)

## Build Prerequisites

Please see the [edgex-go README](https://github.com/edgexfoundry/edgex-go/blob/master/README.md).

## Building Custom Examples

Each `custom` example has its own `Makefile` with the `build` target. So simply run `make build` from an example's base folder when exploring that example.

The top level `makefile` is designed to build all the examples under the `custom` folder. Thus the Makefile does not need to be updated when a new example is added to the `custom` folder

​	run `make build` to build all examples.

For simplicity, the executable created for each example is named `app-service` and placed in each examples sub-folder.

## Running a Custom Example

After building examples you simply `cd` to the folder for the example you want to run and run the executable for that example with or without any of the supported command line options.

The following commands will run the `simple-filter-xml` example

```
cd custom/simple-filter-xml
./app-service
```

## Building App Service Docker Image

The  `simple-filter-xml` example contains an example `Dockerfile` to demonstrate how to build a **Docker Image** for your Application Service. 

The Makefile in this example also contains the `docker` target which will build the **Docker Image** for the example.

​	run `make docker`

This command must be run from the example folder, i.e. not the top level Makefile.

> *Note that Application Services no longer use docker profiles. They use Environment Overrides in the docker compose file to make the necessary changes to the configuration for running in Docker. See the [Environment Variable Overrides](https://docs.edgexfoundry.org/latest/microservices/application/AdvancedTopics/#environment-variable-overrides) section of the Application Services documentation for more details and an example. 

## Configurable Examples

The `configurable` folder contains example profiles for use with App Service Configurable. 