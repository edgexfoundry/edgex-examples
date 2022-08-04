[To README](README.md)

# 4. How to develop custom device service

Now that we have got the core services for EdgeX running on our RPI we can look at how to implement our own device service.

The EdgeX community offers C and go device service SDKs, we will use go for this exercise.
<br/>

## 4.1 How to use EdgeX device service SDK

This chapter shows how to make an example device service based on EdgeX device service SDK, which offers basic capability to communicate with EdgeX core services. We will implement an "echo" feature so that it returns strings as same as it receives. With this example, readers can find API entries and handlers of EdgeX core services. 

Below resources can be used to learn more about the device service SDK:
- https://docs.edgexfoundry.org/2.1/microservices/device/sdk/Ch-DeviceSDK/
- https://docs.edgexfoundry.org/2.1/getting-started/Ch-GettingStartedSDK-Go/
- https://github.com/edgexfoundry/device-sdk-go

To make our own device service, readers should:
- Clone EdgeX device service SDK
- Relocate files
- Edit configuration and Go files
- Compile, launch, and test

<br/>

### 4.1.1 Create Stub Service

For exercise purposes create folder called `repo` with a subfolder `device-echo` to hold our service.  Other examples in this repository have better project structures to follow but for convenience and ease of discussion we will implement the entire service in main.go.  But first some setup is needed.

# Initialize Module and Install Device SDK

This can be done with a few simple commands run in our service directory `device-echo`
```shell
$ go mod init main
go: creating new go.mod: module main
$ go get github.com/edgexfoundry/device-sdk-go/v2@v2.1.0
go get: added github.com/edgexfoundry/device-sdk-go/v2 v2.1.0
```

# Create Service Stub

In `main.go` add the following code using your favorite text editor.  This will have build errors until the driver is implemented in the next step if you are using an IDE.

```go
package main

import (
	"fmt"
	"github.com/edgexfoundry/device-sdk-go/v2"
	sdkModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"time"
)

const (
	serviceName string = "device-echo"
)

func main() {
	sd := echoDriver{}
	startup.Bootstrap(serviceName, device.Version, &sd)
}

```

### 4.1.2 Add Driver

We then need to add a driver as specified in the above links.  This will be for a simple service that accepts string values in PUT commands, and returns the last value received on GETs.  It will also send the current value along the EdgeX message bus every 5 seconds.  The driver implementation is below, I placed it inline in main.go:

```go
type echoDriver struct {
	lc         logger.LoggingClient
	asyncCh    chan<- *sdkModels.AsyncValues
	deviceCh   chan<- []sdkModels.DiscoveredDevice
	echoString string
}

// Initialize performs protocol-specific initialization for the device
// service.
func (ed *echoDriver) Initialize(
	lc logger.LoggingClient,
	asyncCh chan<- *sdkModels.AsyncValues,
	deviceCh chan<- []sdkModels.DiscoveredDevice) error {

	ed.lc = lc
	ed.asyncCh = asyncCh
	ed.deviceCh = deviceCh

	ed.echoString = ""
	go ed.Echo()

	return nil
}

func (ed *echoDriver) Echo() {
	tick := time.Tick(5000 * time.Millisecond)

	for {
		select {
		case <-tick:
			if ed.echoString != "" {
				cValue, _ := sdkModels.NewCommandValue(
					"echoString",
					common.ValueTypeString,
					ed.echoString)

				cValueSlice := make([]*sdkModels.CommandValue, 0)
				cValueSlice = append(cValueSlice, cValue)
				d := sdkModels.AsyncValues{
					DeviceName:    "Echo-Device01",
					CommandValues: cValueSlice,
				}
				ed.lc.Infof("sending command values from Echo: %+v", d)
				ed.asyncCh <- &d
			}
		}
	}
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (ed *echoDriver) HandleReadCommands(
	deviceName string,
	protocols map[string]models.ProtocolProperties,
	reqs []sdkModels.CommandRequest) (res []*sdkModels.CommandValue, err error) {

	ed.lc.Debug(fmt.Sprintf("echoDriver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

	if len(reqs) == 1 {
		res = make([]*sdkModels.CommandValue, 1)
		if reqs[0].DeviceResourceName == "echoString" {
			cv, _ := sdkModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, ed.echoString)
			res[0] = cv
		}
	}

	return
}

// HandleWriteCommands triggers a protocol Write operation for the specified device.
func (ed *echoDriver) HandleWriteCommands(
	deviceName string,
	protocols map[string]models.ProtocolProperties,
	reqs []sdkModels.CommandRequest,
	params []*sdkModels.CommandValue) error {

	var err error

	for i, r := range reqs {
		ed.lc.Info(fmt.Sprintf("echoDriver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v", protocols, reqs[i].DeviceResourceName, params[i]))

		switch r.DeviceResourceName {
		case "echoString":
			if ed.echoString, err = params[i].StringValue(); err != nil {
				err := fmt.Errorf("echoDriver.HandleWriteCommands; the data type of parameter should be string, parameter: %ed", params[0].String())
				return err
			}
		}
	}

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (ed *echoDriver) Stop(force bool) error {
	if ed.lc != nil {
		ed.lc.Debugf("echoDriver.Stop called: force=%v", force)
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (ed *echoDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	ed.lc.Debugf("a new Device is added: %ed", deviceName)
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (ed *echoDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	ed.lc.Debugf("Device %ed is updated", deviceName)
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (ed *echoDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	ed.lc.Debugf("Device %ed is removed", deviceName)
	return nil
}
```

This is a fairly straightforward implementation of the device driver interface - the one wrinkle is our Echo method - this is run on a separate goroutine (in the background) by calling with `go ed.Echo()`

### 4.1.3 Add Config files
# res/configuration.toml
```toml
[Writable]
    LogLevel = "INFO"
    # Example InsecureSecrets configuration that simulates SecretStore for when EDGEX_SECURITY_SECRET_STORE=false
    # InsecureSecrets are required for when Redis is used for message bus
    [Writable.InsecureSecrets]
        [Writable.InsecureSecrets.DB]
            path = "redisdb"
        [Writable.InsecureSecrets.DB.Secrets]
            username = ""
            password = ""

[Service]
    HealthCheckInterval = "10s"
    Host = "localhost"
    Port = 59999 # Device serivce are assigned the 599xx range
    ServerBindAddr = ""  # blank value defaults to Service.Host value
    StartupMsg = "device echo started"
    # MaxRequestSize limit the request body size in byte of put command
    MaxRequestSize = 0 # value 0 unlimit the request size.
    RequestTimeout = "20s"
    [Service.CORSConfiguration]
        EnableCORS = false
        CORSAllowCredentials = false
        CORSAllowedOrigin = "https://localhost"
        CORSAllowedMethods = "GET, POST, PUT, PATCH, DELETE"
        CORSAllowedHeaders = "Authorization, Accept, Accept-Language, Content-Language, Content-Type, X-Correlation-ID"
        CORSExposeHeaders = "Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma, X-Correlation-ID"
        CORSMaxAge = 3600

[Registry]
    Host = "localhost"
    Port = 8500
    Type = "consul"

[Clients]
    [Clients.core-data]
        Protocol = "http"
        Host = "localhost"
        Port = 59880

    [Clients.core-metadata]
        Protocol = "http"
        Host = "localhost"
        Port = 59881

[MessageQueue]
    Protocol = "redis"
    Host = "localhost"
    Port = 6379
    Type = "redis"
    AuthMode = "usernamepassword"  # required for redis messagebus (secure or insecure).
    SecretName = "redisdb"
    PublishTopicPrefix = "edgex/events/device" # /<device-profile-name>/<device-name>/<source-name> will be added to this Publish Topic prefix
    [MessageQueue.Optional]
        # Default MQTT Specific options that need to be here to enable environment variable overrides of them
        # Client Identifiers
        ClientId = "device-echo"
        # Connection information
        Qos = "0" # Quality of Sevice values are 0 (At most once), 1 (At least once) or 2 (Exactly once)
        KeepAlive = "10" # Seconds (must be 2 or greater)
        Retained = "false"
        AutoReconnect = "true"
        ConnectTimeout = "5" # Seconds
        SkipCertVerify = "false" # Only used if Cert/Key file or Cert/Key PEMblock are specified

# Example SecretStore configuration.
# Only used when EDGEX_SECURITY_SECRET_STORE=true
# Must also add `ADD_SECRETSTORE_TOKENS: "device-echo"` to vault-worker environment so it generates
# the token and secret store in vault for "device-echo"
[SecretStore]
    Type = "vault"
    Host = "localhost"
    Port = 8200
    Path = "device-echo/"
    Protocol = "http"
    RootCaCertPath = ""
    ServerName = ""
    SecretsFile = ""
    DisableScrubSecretsFile = false
    TokenFile = "/tmp/edgex/secrets/device-echo/secrets-token.json"
    [SecretStore.Authentication]
        AuthType = "X-Vault-Token"

[Device]
    DataTransform = true
    MaxCmdOps = 128
    MaxCmdValueLen = 256
    ProfilesDir = "./res/profiles"
    DevicesDir = "./res/devices"
    UpdateLastConnected = false
    AsyncBufferSize = 1
    EnableAsyncReadings = true
    Labels = []
    UseMessageBus = true
    [Device.Discovery]
    Enabled = false
    Interval = "30s"

```

# res/profiles/Echo-Driver.yaml
```yaml
apiVersion: "v2"
name: "Echo-Device"
manufacturer: "Simple Corp."
model: "ED-01"
labels:
  - "sample"
description: "Example of Echo Device"

deviceResources:
  -
    name: "echoString"
    isHidden: false
    description: "Echo String"
    properties:
      valueType: "String"
      readWrite: "RW"
      defaultValue: ""

deviceCommands:
  -
    name: "echoString"
    isHidden: false
    readWrite: "RW"
    resourceOperations:
      - { deviceResource: "echoString", defaultValue: "" }
```

# res/devices/echo-device.toml
```toml
[[DeviceList]]
  Name = "Echo-Device01"
  ProfileName = "Echo-Device"
  Description = "Example of Echo Device"
  Labels = [ "industrial" ]
  [DeviceList.Protocols]
  [DeviceList.Protocols.other]
    Address = "echo01"
    Port = "300"
```
<br/>

### 4.1.4 Add Dockerfile

Example below:

```dockerfile
FROM golang:1.17-alpine3.14 AS builder

WORKDIR /temp

LABEL license='SPDX-License-Identifier: Apache-2.0'

RUN apk add --update --no-cache make git gcc libc-dev zeromq-dev libsodium-dev

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o ./device-echo

FROM alpine:3.14 as final

RUN apk add --update --no-cache zeromq

WORKDIR /
COPY --from=builder /temp/device-echo /device-echo
COPY ./res/ /res

EXPOSE 59999

ENTRYPOINT ["/device-echo", "-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"]
```
### 4.1.5 Launch and test

Our `device-echo` directory should now look like this:

```shell
$ tree
.
├── Dockerfile
├── go.mod
├── main.go
└── res
    ├── configuration.toml
    ├── devices
    │   └── echo-device.toml
    ├── profiles
    │   └── Echo-Driver.yaml

3 directories, 6 files
```

Once this is done we can test building using `docker build` but to easily join it with our running services from previous steps we can create a file with the following compose snippet in `device-echo`

```yaml
networks:
  edgex-network:
    external: true
services:
  device-echo:
    container_name: edgex-device-echo
    depends_on:
      - consul
      - data
      - metadata
    environment:
      CLIENTS_CORE_COMMAND_HOST: edgex-core-command
      CLIENTS_CORE_DATA_HOST: edgex-core-data
      CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
      CLIENTS_SUPPORT_NOTIFICATIONS_HOST: edgex-support-notifications
      CLIENTS_SUPPORT_SCHEDULER_HOST: edgex-support-scheduler
      DATABASES_PRIMARY_HOST: edgex-redis
      EDGEX_SECURITY_SECRET_STORE: "false"
      MESSAGEQUEUE_HOST: edgex-redis
      REGISTRY_HOST: edgex-core-consul
      SERVICE_HOST: edgex-device-echo
    hostname: edgex-device-echo
    build:
      context: device-echo/.
    networks:
      edgex-network: { }
    ports:
      - 127.0.0.1:59999:59999/tcp
    read_only: true
    restart: always
    security_opt:
      - no-new-privileges:true
    user: 2002:2001
```
From there we can attempt to run the service:

```sh
$ docker-compose up --build
```

Please open a new terminal, login to the RPI, and use **curl** to check the state of the device service:
```sh
# check device registration via metadata service
$ curl localhost:59881/api/v2/device/name/Echo-Device01
{"apiVersion":"v2","statusCode":200,"device":{"created":1642904140586,"modified":1642904140586,"id":"fed088a9-95d4-449c-b651-8681aacdb1ae","name":"Echo-Device01","description":"Example of Echo Device","adminState":"UNLOCKED","operatingState":"UP","labels":["industrial"],"serviceName":"device-echo","profileName":"Echo-Device","protocols":{"other":{"Address":"echo01","Port":"300"}}}}

# check current value of echoString directly via device service
$ curl localhost:59999/api/v2/device/name/Echo-Device01/echoString
{"apiVersion":"v2","statusCode":200,"event":{"apiVersion":"v2","id":"b68134fb-5a64-4f90-92a6-b1c25a2458e2","deviceName":"Echo-Device01","profileName":"Echo-Device","sourceName":"echoString","origin":1642907383831644371,"readings":[{"id":"4aa8810f-be8e-4354-8330-58a60ce560ad","origin":1642907383831644371,"deviceName":"Echo-Device01","resourceName":"echoString","profileName":"Echo-Device","valueType":"String"}]}}

# check current value of echoString via command service
$ curl localhost:59882/api/v2/device/name/Echo-Device01/echoString
{"apiVersion":"v2","statusCode":200,"event":{"apiVersion":"v2","id":"b68134fb-5a64-4f90-92a6-b1c25a2458e2","deviceName":"Echo-Device01","profileName":"Echo-Device","sourceName":"echoString","origin":1642907383831644371,"readings":[{"id":"4aa8810f-be8e-4354-8330-58a60ce560ad","origin":1642907383831644371,"deviceName":"Echo-Device01","resourceName":"echoString","profileName":"Echo-Device","valueType":"String"}]}}

# set value for echoString via command service
$ curl localhost:59882/api/v2/device/name/Echo-Device01/echoString -X PUT -d '{"echoString":"test completed"}'
{"apiVersion":"v2","statusCode":200}

# check value is set via command service
$ curl localhost:59882/api/v2/device/name/Echo-Device01/echoString
{"apiVersion":"v2","statusCode":200,"event":{"apiVersion":"v2","id":"e434eafe-345a-458b-96bc-9ab78a0d6035","deviceName":"Echo-Device01","profileName":"Echo-Device","sourceName":"echoString","origin":1642907593844892877,"readings":[{"id":"e53312e8-4d21-47ad-8614-20e07cca940c","origin":1642907593844892877,"deviceName":"Echo-Device01","resourceName":"echoString","profileName":"Echo-Device","valueType":"String","value":"test completed"}]}}

```

More info on the command service and various EdgeX apis can be found at the below links:
- https://docs.edgexfoundry.org/2.1/api/core/Ch-APICoreCommand/
- https://app.swaggerhub.com/search?type=API&query=%20edgex

<br/>

---

Next: [How to develop custom app services](50_custom_app_services.md)