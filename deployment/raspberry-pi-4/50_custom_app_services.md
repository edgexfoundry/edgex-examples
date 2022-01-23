[To README](README.md)

## 5 How to use EdgeX app functions SDK

Our first custom device service works good with the EdgeX services. Now, it is the time to create our own custom app service, which gets messages from the device service via the core data.

To make our own app service, readers should:
- Clone EdgeX app functions SDK
- Edit the configuration file
- Compile, launch, and test

EdgeX foundry offers plenty of documents as well:
- https://docs.edgexfoundry.org/2.1/getting-started/ApplicationFunctionsSDK/
- https://docs.edgexfoundry.org/2.1/microservices/application/ApplicationServices/
- https://docs.edgexfoundry.org/2.1/examples/AppServiceExamples/
- https://docs.edgexfoundry.org/2.1/getting-started/ApplicationFunctionsSDK/
- https://github.com/edgexfoundry/app-functions-sdk-go
- https://github.com/edgexfoundry/edgex-examples/blob/master/application-services/custom/simple-filter-xml/main.go

<br/>

### 5.1 Build app functions SDK example

For exercise purposes we will work in the folder `repo` with a new subfolder `app-echo` to hold our service.  Other examples in this repository have better project structures to follow but for convenience and ease of discussion we will implement the entire service in main.go.  But first some setup is needed.  We are using the `simple-filter-xml` example as the basis for our service.

# Initialize Module and Install App Function SDK

This can be done with a few simple commands run in our service directory `app-echo`
```shell
$ go mod init main
go: creating new go.mod: module main
$ go get github.com/edgexfoundry/app-functions-sdk-go/v2@v2.1.0
go get: added github.com/edgexfoundry/app-functions-sdk-go/v2 v2.1.0
```

# Add Service

Create a new file `main.go` with the following content:

```go
package main

import (
	"errors"
	"fmt"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"os"
)

const (
	serviceKey = "app-echo"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an new instance of an EdgeX Application Service.
	service, ok := pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(-1)
	}

	// Leverage the built in logging service in EdgeX
	lc := service.LoggingClient()

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := service.GetAppSettingStrings("DeviceNames")
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}
	lc.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		transforms.NewConversion().TransformToXML,
		printXMLToConsole,
	); err != nil {
		lc.Errorf("SetFunctionsPipeline returned error: %s", err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func printXMLToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Leverage the built in logging service in EdgeX
	lc := ctx.LoggingClient()

	if data == nil {
		return false, errors.New("printXMLToConsole: No data received")
	}

	xml, ok := data.(string)
	if !ok {
		return false, errors.New("printXMLToConsole: Data received is not the expected 'string' type")
	}

	lc.Debug(xml)
	ctx.SetResponseData([]byte(xml))
	return true, xml
}

```



### 5.2 Update Configuration

The app functions SDK offers handlers and filters for the message stream of EdgeX core data service. Examples of app functions SDK show various use cases but we will start with a very simple one - echoing events coming through the message bus for our new device `Device-Echo01` only.

The **res/configuration.toml** is the configuration file for this app function. The sub section **ApplicationSettings** should have device names as target message sources. Since our device service has the name **Echo-Device01** registered, we need to write the same name for DeviceNames as below.

```toml
[Writable]
    LogLevel = "INFO"
    [Writable.StoreAndForward]
        Enabled = false
        RetryInterval = "5m"
        MaxRetryCount = 10
    [Writable.InsecureSecrets]
        [Writable.InsecureSecrets.DB]
            path = "redisdb"
        [Writable.InsecureSecrets.DB.Secrets]
            username = ""
            password = ""

[Service]
    HealthCheckInterval = "10s"
    Host = "localhost"
    Port = 59998 # Adjust if running multiple examples at the same time to avoid duplicate port conflicts
    ServerBindAddr = "" # if blank, uses default Go behavior https://golang.org/pkg/net/#Listen
    StartupMsg = "This is a sample Filter/XML Transform Application Service"
    RequestTimeout = "30s"
    MaxRequestSize = 0
    MaxResultCount = 0

[Registry]
    Host = "localhost"
    Port = 8500
    Type = "consul"

# Database is require when Store and Forward is enabled
[Database]
    Type = "redisdb"
    Host = "localhost"
    Port = 6379
    Timeout = "30s"

# SecretStore is required when Store and Forward is enabled and running with security
# so Database credentials can be pulled from Vault.
# Note when running in docker from compose file set the following environment variables:
#   - SecretStore_Host: edgex-vault
[SecretStore]
    Type = "vault"
    Host = "localhost"
    Port = 8200
    Path = "app-simple-filter-xml/"
    Protocol = "http"
    TokenFile = "/tmp/edgex/secrets/app-simple-filter-xml/secrets-token.json"
    RootCaCertPath = ""
    ServerName = ""
    [SecretStore.Authentication]
        AuthType = "X-Vault-Token"

[Clients]
    [Clients.core-metadata]
    Protocol = "http"
    Host = "localhost"
    Port = 59881

# Choose either the http trigger or edgex-messagebus trigger

#[Trigger]
#Type="http"

[Trigger]
    Type="edgex-messagebus"
    [Trigger.EdgexMessageBus]
        Type = "redis"
    [Trigger.EdgexMessageBus.SubscribeHost]
        Host = "localhost"
        Port = 6379
        Protocol = "redis"
        SubscribeTopics="edgex/events/#"
    [Trigger.EdgexMessageBus.PublishHost]
        Host = "localhost"
        Port = 6379
        Protocol = "redis"
        PublishTopic="example"

# App Service specifc simple settings
# Great for single string settings. For more complex structured custom configuration
# See https://docs.edgexfoundry.org/2.0/microservices/application/AdvancedTopics/#custom-configuration
[ApplicationSettings]
    DeviceNames = "Echo-Device01"
```

### 5.3 Add Dockerfile

We can use the same dockerfile as with device services (both require ZeroMQ present to build/run).  Just need to specify a different port.

```dockerfile
FROM golang:1.17-alpine3.14 AS builder

WORKDIR /temp

LABEL license='SPDX-License-Identifier: Apache-2.0'

RUN apk add --update --no-cache make git gcc libc-dev zeromq-dev libsodium-dev

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o ./app-echo

FROM alpine:3.14 as final

RUN apk add --update --no-cache zeromq

WORKDIR /
COPY --from=builder /temp/app-echo /app-echo
COPY ./res/ /res

EXPOSE 59798

ENTRYPOINT ["/app-echo", "-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"]
```

### 5.4 Launch and test

Our `app-echo` directory should now look like this:

```shell
$ tree
.
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
└── res
    └── configuration.toml

1 directory, 5 files

```

Once this is done we can test building using `docker build` but to easily join it with our running services from previous steps we can create a file with the following compose snippet in `device-echo`

```yaml
networks:
  edgex-network:
    external: true
services:
  app-echo:
    container_name: edgex-app-echo
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
      TRIGGER_EDGEXMESSAGEBUS_SUBSCRIBEHOST_HOST: edgex-redis
      TRIGGER_EDGEXMESSAGEBUS_PUBLISHHOST_HOST: edgex-redis
      REGISTRY_HOST: edgex-core-consul
      SERVICE_HOST: edgex-app-echo
    hostname: edgex-app-echo
    build:
      context: app-echo/.
    networks:
      edgex-network: { }
    ports:
      - 127.0.0.1:59798:59798/tcp
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

Once services start the logs should be fairly quiet - we will need to set an echo value before the events start sending:

```shell
$ curl localhost:59882/api/v2/device/name/Echo-Device01/echoString -X PUT -d '{"echoString":"test completed"}'
{"apiVersion":"v2","statusCode":200}
```

After this is set we should start seeing pairs of log messages like this:

```
edgex-device-echo    | level=INFO ts=2022-01-23T13:25:02.457868491Z app=device-echo source=main.go:66 msg="sending command values from Echo: {DeviceName:Echo-Device01 SourceName: CommandValues:[DeviceResource: echoString, String: test completed]}"
edgex-app-echo       | level=INFO ts=2022-01-23T13:25:02.458789589Z app=app-echo source=main.go:74 msg="<Event><ApiVersion>v2</ApiVersion><Id>fb9ecfe1-57bb-4f1b-8d9b-408526e41951</Id><DeviceName>Echo-Device01</DeviceName><ProfileName>Echo-Device</ProfileName><SourceName>echoString</SourceName><Origin>1642944302457960155</Origin><Readings><Id>d7365a0c-9c91-40a5-9333-a4ab7fe5307d</Id><Origin>1642944302457960155</Origin><DeviceName>Echo-Device01</DeviceName><ResourceName>echoString</ResourceName><ProfileName>Echo-Device</ProfileName><ValueType>String</ValueType><BinaryValue></BinaryValue><MediaType></MediaType><Value>test completed</Value></Readings></Event>"
```
## Conclusion

I hope you have enjoyed this introduction.  For a recap, we prepared Ubuntu server 21.10 on RPI, launched EdgeX services, and created custom device and app services. While these are simplistic examples they should demonstrate how custom services can be connected to the EdgeX ecosystem and the way data flows between them on the message bus.  Feel free to explore the other examples in this repository for further inspiration.

<br/>

---

[To README](README.md)
