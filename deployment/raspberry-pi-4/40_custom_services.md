[To README](README.md)

# 4. How to develop custom device and app services 

Launching the core services is the first step of EdgeX based development. However, how can we connect our own device and app services to the core services? How can our services communicate with the core services? EdgeX team offers SDKs for custom device and app services so that we can start from the SDKs.

<br/>

## 4.1 How to use EdgeX device service SDK

This chapter shows how to make an example device service based on EdgeX device service SDK, which offers basic capability to communicate with EdgeX core services. We will implement an "echo" feature so that it returns strings as same as it receives. With this example, readers can find API entries and handlers of EdgeX core services. 

Below resources can be used to learn more about the device service SDK:
- https://docs.edgexfoundry.org/1.2/microservices/device/sdk/Ch-DeviceSDK/
- https://docs.edgexfoundry.org/1.2/getting-started/Ch-GettingStartedSDK-Go/
- https://github.com/edgexfoundry/device-sdk-go

To make our own device service, readers should:
- Clone EdgeX device service SDK
- Relocate files
- Edit configuration and Go files
- Compile, launch, and test

<br/>

### 4.1.1 Clone and config SDK

The device service SDK can be cloned and configured so that we can build/test the service as follow:
```sh
# Clone SDK
$ mkdir ~/go/src/github.com/edgexfoundry -p
$ cd ~/go/src/github.com/edgexfoundry
$ git clone --depth 1 \
    --branch v1.2.2 https://github.com/edgexfoundry/device-sdk-go.git

# Relocate files
$ mkdir ~/repo/device-simple -p
$ cp -rf device-sdk-go/example/* ~/repo/device-simple/
$ cp device-sdk-go/Makefile ~/repo/device-simple/
$ cp device-sdk-go/version.go ~/repo/device-simple/
$ cd ~/repo/device-simple

# Check contents
$ tree
.
├── cmd
│   └── device-simple
│       ├── Attribution.txt
│       ├── Dockerfile
│       ├── main.go
│       └── res
│           ├── configuration.toml
│           ├── off.jpg
│           ├── on.png
│           ├── provisionwatcher.json
│           └── Simple-Driver.yaml
├── driver
│   └── simpledriver.go
├── Makefile
├── README.md
└── version.go

4 directories, 12 files

# Edit main.go: this command removes "example/" 
$ sed -i '/\"github.com\/edgexfoundry\/device-sdk-go\/example\/driver\"/c\\t\"main\/driver\"' ./cmd/device-simple/main.go

# Edit Makefile: this command replaces "device-sdk-go" to "device-simple"
$ sed -i '/GOFLAGS=-ldflags \"-X github.com\/edgexfoundry\/device-sdk-go.Version=$(VERSION)\"/c\GOFLAGS=-ldflags \"-X github.com\/edgexfoundry\/device-simple.Version=$(VERSION)\"' ./Makefile

# Edit Makefile: these commands remove "example/"
$ sed -i '/MICROSERVICES=example\/cmd\/device-simple\/device-simple/c\MICROSERVICES=cmd\/device-simple\/device-simple' ./Makefile
$ sed -i '/example\/cmd\/device-simple\/device-simple:/c\cmd\/device-simple\/device-simple:' ./Makefile
$ sed -i '/$(GO) build $(GOFLAGS) -o $@ .\/example\/cmd\/device-simple/c\\t$(GO) build $(GOFLAGS) -o $@ .\/cmd\/device-simple' ./Makefile

# Enable Go module  
$ go mod init main
$ echo "require (
    github.com/edgexfoundry/device-sdk-go v1.2.2
    github.com/edgexfoundry/go-mod-core-contracts v0.1.58
)" >> go.mod

# Test build
$ make build
$ file cmd/device-simple/device-simple
cmd/device-simple/device-simple: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, Go BuildID=azNP-8ouXJ3WjHlZkFSb/XRFJnCq3Yoz8hQqzObnB/Q1HqdTQy74eFvb6rvbrv/d3hOwwcqL_RurDo7Aesj, not stripped
```

<br/>

### 4.1.2 Update handlers

There is a file, which has all the handlers for us so that we can edit the code to interact with the core services:
```sh
# Update go file for handlers
$ vi driver/simpledriver.go
```

In the go file, we need to change:
- SimpleDriver struct
- Initialize method
- HandleReadCommands method
- HandleWriteCommands method
- HandleEcho method (new!)
- Other parts shouldn't be changed
 
```go
type SimpleDriver struct {
        lc           logger.LoggingClient
        asyncCh      chan<- *dsModels.AsyncValues
        deviceCh     chan<- []dsModels.DiscoveredDevice
        switchButton bool
        xRotation    int32
        yRotation    int32
        zRotation    int32
        echoString   string // Added
}

func (s *SimpleDriver) Initialize(
    lc logger.LoggingClient, 
    asyncCh chan<- *dsModels.AsyncValues, 
    deviceCh chan<- []dsModels.DiscoveredDevice) error {

        s.lc = lc
        s.asyncCh = asyncCh
        s.deviceCh = deviceCh

        s.echoString = "" // Added
        go s.Echo() // Added

        return nil
}

// A new function
func (s *SimpleDriver) Echo(){
    tick := time.Tick(5000 * time.Millisecond)
     
    for { 
        select{
            case <-tick:
                if s.echoString != "" {
                    cValue := dsModels.NewStringValue(
                        "echoString",
                        int64(time.Now().Unix()), s.echoString)

                    cValueSlice := make([]*dsModels.CommandValue, 0)
                    cValueSlice = append(cValueSlice, cValue)
                    d := dsModels.AsyncValues{
                        DeviceName:    "Simple-Device02",
                        CommandValues: cValueSlice,
                    }
                    s.asyncCh <- &d
                }
        }
    }
}

func (s *SimpleDriver) HandleReadCommands(
    deviceName string, 
    protocols map[string]contract.ProtocolProperties, 
    reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

        s.lc.Debug(fmt.Sprintf("SimpleDriver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

        // Replaced the contents of this if block
        if len(reqs) == 1 {
                res = make([]*dsModels.CommandValue, 1)
                now := time.Now().UnixNano()
                if reqs[0].DeviceResourceName == "echoString" {
                        cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, s.echoString)
                        res[0] = cv
                } 
        }

        return
}

func (s *SimpleDriver) HandleWriteCommands(
    deviceName string, 
    protocols map[string]contract.ProtocolProperties, 
    reqs []dsModels.CommandRequest,
    params []*dsModels.CommandValue) error {

        var err error

        // Replaced the contents of this for block
        for i, r := range reqs {
                s.lc.Info(fmt.Sprintf("SimpleDriver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v", protocols, reqs[i].DeviceResourceName, params[i]))

                switch r.DeviceResourceName {
                case "echoString":
                        if s.echoString, err = params[i].StringValue(); err != nil {
                                err := fmt.Errorf("SimpleDriver.HandleWriteCommands; the data type of parameter should be Boolean, parameter: %s", params[0].String())
                                return err
                        }
                }
        }

        return nil
}
```

Code formatting and test build:
```sh
$ gofmt -s -w .
$ make build
```

<br/>

### 4.1.3 Update config files

As the code got updated, this device service can handle read and write requests from core services. The device service also needs to register itself to the core services and these are the files to be used for the registration:
- ~/repo/device-simple/cmd/device-simple/res/**configuration.toml**
- ~/repo/device-simple/cmd/device-simple/res/**Simple-Driver.yaml**

For configuration.toml:
```toml
[Writable]
LogLevel = 'INFO'

[Service]
BootTimeout = 30000
CheckInterval = '10s'
ClientMonitor = 15000
Host = '172.17.0.1' # If the core services run as Docker containers
Port = 49980 # Don't use a port being used!
Protocol = 'http'
StartupMsg = 'device simple started'
Timeout = 20000
ConnectRetries = 20
Labels = []
EnableAsyncReadings = true
AsyncBufferSize = 16

[Registry]
Host = 'localhost'
Port = 8500
Type = 'consul'

[Clients]
  [Clients.Data]
  Protocol = 'http'
  Host = 'localhost'
  Port = 48080

  [Clients.Metadata]
  Protocol = 'http'
  Host = 'localhost'
  Port = 48081

  [Clients.Logging]
  Protocol = 'http'
  Host = 'localhost'
  Port = 48061

[Device]
  DataTransform = true
  InitCmd = ''
  InitCmdArgs = ''
  MaxCmdOps = 128
  MaxCmdValueLen = 256
  RemoveCmd = ''
  RemoveCmdArgs = ''
  ProfilesDir = './res'
  UpdateLastConnected = false
  [Device.Discovery]
    Enabled = false
    Interval = '30s'

# Remote and file logging disabled so only stdout logging is used
[Logging]
EnableRemote = false
File = ''

# Pre-define Devices
[[DeviceList]]
  Name = 'Simple-Device02'
  Profile = 'Simple-Device'
  Description = 'Example of Simple Device'
  Labels = [ 'industrial' ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]
      Address = 'simple01'
      Port = '300'

# Auto events are removed for the purpose of this tutorial
```

For Simple-Driver.yaml (be careful about the indentation!):
```yaml
name: "Simple-Device"
manufacturer: "HP Corp."
model: "ED-01"
description: "Example of Simple Echo Device"

deviceResources:
  -
    name: "echoString"
    description: "Echo String"
    properties:
      value:
        { type: "String", readWrite: "RW", defaultValue: "" }
      units:
        { type: "String", readWrite: "R", defaultValue: "" }

deviceCommands:
  -
    name: "echoString"
    get:
      - { operation: "get", deviceResource: "echoString" }
    set:
      - { operation: "set", deviceResource: "echoString", parameter: "false" }

coreCommands:
  -
    name: "echoString"
    get:
      path: "/api/v1/device/{deviceId}/echoString"
      responses:
        -
          code: "200"
          description: ""
          expectedValues: ["echoString"]
        -
          code: "503"
          description: "echo string unavailable"
          expectedValues: []
    put:
      path: "/api/v1/device/{deviceId}/echoString"
      parameterNames: ["echoString"]
      responses:
        -
          code: "200"
          description: ""
        -
          code: "503"
          description: "echo string unavailable"
          expectedValues: []
```

<br/>

### 4.1.4 Launch and test

The code was compiled well and the files for registration are ready. When the binary of this device service is executed, it does bootstrapping for the communication with the core services and uses the files to tell what it is and available commands. To excute: 
```sh
$ cd cmd/device-simple
$ ./device-simple
...
level=INFO ts=2020-09-15T10:48:51.813731783Z app=device-simple source=init.go:42 msg="Service clients initialize successful."
level=INFO ts=2020-09-15T10:48:51.814815574Z app=device-simple source=service.go:83 msg="Device Service device-simple doesn't exist, creating a new one"
...
```

Please open a new terminal, login to the RPI, and use **curl** to check the state of the device service:
```sh
# Basic info of the device service.
$ curl http://localhost:48081/api/v1/addressable -X GET -s | jq '.[] | {name,address,port}'
{
  "name": "device-simple",
  "address": "172.17.0.1",
  "port": 49980
}

# Returned IDs of our device service may vary.
$ curl http://localhost:48081/api/v1/device/name/Simple-Device02 -X GET -s -S | jq '.name,.id,.service.id'
"Simple-Device02"
"9be2790a-dab9-447f-ac59-74505527252f"
"20a1fec0-de09-4af8-bb20-000b33573f9e"

# Returned ID of the value descriptor may vary.
$ curl http://localhost:48080/api/v1/valuedescriptor/name/echoString -X GET -s -S | jq '.id'
"0242e148-a5d3-40a5-a850-c5af5dd8456f"

# To find out the available device name and id. The retured URL may vary.
$ curl http://localhost:48082/api/v1/device -s | jq '.[] | select(.name == "Simple-Device02").id,.name'
"9be2790a-dab9-447f-ac59-74505527252f"
"Simple-Device02"

# To find out the available commands id. The returned URL may vary.
$ curl http://localhost:48082/api/v1/device -s | jq '.[] | select(.name == "Simple-Device02").commands | .[] | .id'
"80f2ba5e-73d1-4130-b039-c158a2f43388"

# With the commands above, we could find the ID of the device/command
# Device ID: "9be2790a-dab9-447f-ac59-74505527252f"
# Command ID: "80f2ba5e-73d1-4130-b039-c158a2f43388"

# To change echoString value, let's store the IDs first:
$ DEVICE_ID=$(curl http://localhost:48082/api/v1/device -s | jq -r '.[] | select(.name == "Simple-Device02").id')
$ COMMAND_ID=$(curl http://localhost:48082/api/v1/device -s | jq -r '.[] | select(.name == "Simple-Device02").commands | .[] | .id')

# Query with the gathered IDs:
$ curl http://localhost:48082/api/v1/device/$DEVICE_ID/command/$COMMAND_ID \
    -s \
    -X PUT \
    -H 'Content-Type: application/json' \
    -H 'cache-control: no-cache' \
    -d '{"echoString": "HELLO"}'

# To check the value changed in 3 different ways

# Asking to the core data service
$ curl http://localhost:48080/api/v1/reading/device/Simple-Device02/1 -X GET -s -S | json_pp

# Asking to the core command service
$ curl http://localhost:48082/api/v1/device/name/Simple-Device02/command/echoString -X GET | json_pp

# Asking to the device service itself
$ curl http://localhost:49980/api/v1/device/name/Simple-Device02/echoString -X GET -s -S | json_pp

...
"name" : "echoString",
"value" : "HELLO",
...

# To check the latest 10 async events/readings of the device service via the core data service
$ curl -s http://localhost:48080/api/v1/event/device/Simple-Device02/10 | json_pp
```

More APIs can be found from:
- https://docs.edgexfoundry.org/1.2/api/core/Ch-APICoreCommand/
- https://app.swaggerhub.com/search?type=API&query=%20edgex

<br/>

## 4.2 How to use EdgeX app functions SDK

Our first custom device service works good with the EdgeX services. Now, it is the time to create our own custom app service, which gets messages from the device service via the core data.

To make our own app service, readers should:
- Clone EdgeX app functions SDK
- Edit the configuration file
- Compile, launch, and test

EdgeX foundry offers plenty of documents as well:
- https://docs.edgexfoundry.org/1.2/getting-started/ApplicationFunctionsSDK/
- https://docs.edgexfoundry.org/1.2/microservices/application/ApplicationServices/
- https://docs.edgexfoundry.org/1.2/examples/AppServiceExamples/
- https://docs.edgexfoundry.org/1.2/getting-started/ApplicationFunctionsSDK/
- https://github.com/edgexfoundry/app-functions-sdk-go
- https://github.com/edgexfoundry/edgex-examples/blob/master/application-services/custom/simple-filter-xml/main.go

<br/>

### 4.2.1 Build app functions SDK example

Now, we can clone and build one of the app functions SDK examples:
```sh
# Originally, the app functions SDK requires the libzmq library and sometimes we need to build it from the source code if we use OS doesn't deliver the library package out of the box. However, Ubuntu server 20.10 comes with the libzmq3-dev package and it is already installed in the previous chapter.

$ cd ~/repo
$ git clone https://github.com/edgexfoundry/edgex-examples
$ cp -rf edgex-examples/application-services/custom/simple-filter-xml .
$ cd simple-filter-xml
$ tree
.
├── Dockerfile
├── EdgeX Application Function SDK Device Name.postman_collection.json
├── EdgeX Applications Function SDK.postman_collection.json
├── go.mod
├── main.go
├── Makefile
└── res
    └── configuration.toml

1 directory, 7 files
```

Change go.mod of simple-filter-xml to be:
```go
go 1.15

require (
	github.com/edgexfoundry/app-functions-sdk-go v1.2.0
)
```

Test build:
```sh
$ make build
CGO_ENABLED=1 GO111MODULE=on go build -o app-service
```

If test build fails, try this and build again:
```sh
$ go get github.com/rjeczalik/pkgconfig/cmd/pkg-config
```

<br/>

### 4.2.2 Edit the configuration file

The app functions SDK offers handlers and filters for the message stream of EdgeX core data service. Examples of app functions SDK show various use cases but we can start from the simplest one. In the previous step, the example is already compiled but we need to take a look into the main.go and res/configuration.toml files.

The **res/configuration.toml** is the configuration file for this app function. Target message source can be specified as long as other settings. The sub section **ApplicationSettings** should have device names as target message sources. Since our device service has the name as **Simple-Device02** in its configuration.toml, we need to write the same name for DeviceNames as below.

```toml
[ApplicationSettings]
DeviceNames = "Simple-Device02"
```

The **main.go** is the place where the actual handlers and filters can be written. The structure and flow are straightforward. In the main function, it initializes the app SDK with a secret key. Then it reads the configuration.toml file and DeviceNames variable. Pipeline is configured with chained functions for message handling and filtering. The printXMLToConsole function is specified at the end of the chained functions so that we can write some code there to use the data filtered from the pipeline so that the incoming messages can be passed to other go routines as we normally write Go code. 

```go
func main() {
        // turn off secure mode for examples. Not recommended for production
        os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

        // 1) First thing to do is to create an instance of the EdgeX SDK and initialize it.
        edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
        if err := edgexSdk.Initialize(); err != nil {
                edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
                os.Exit(-1)
        }

        // 2) shows how to access the application's specific configuration settings.
        deviceNames, err := edgexSdk.GetAppSettingStrings("DeviceNames")
        if err != nil {
                edgexSdk.LoggingClient.Error(err.Error())
                os.Exit(-1)
        }
        edgexSdk.LoggingClient.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

        // 3) This is our pipeline configuration, the collection of functions to
        // execute every time an event is triggered.
        // Also, "TransformToXML" can be edited as "TransformToJSON" for JSON format.
        edgexSdk.SetFunctionsPipeline(
                transforms.NewFilter(deviceNames).FilterByDeviceName,
                transforms.NewConversion().TransformToXML,
                printXMLToConsole,
        )

        // 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
        // to trigger the pipeline.
        err = edgexSdk.MakeItRun()
        if err != nil {
                edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
                os.Exit(-1)
        }

        // Do any required cleanup here

        os.Exit(0)
}
```

### 4.2.3 Customize app functions SDK example

Since we already compiled this example (and main.go is not changed), we can just launch it:
```sh
$ ./app-service
...
level=INFO ts=2020-09-18T10:10:22.624535012Z app=sampleFilterXml source=server.go:350 msg="Starting HTTP Web Server on port :48095"
...

<Event><ID>b53ae300-6bcc-42a4-bf3b-f58165d890f3</ID><Pushed>0</Pushed><Device>Simple-Device02</Device><Created>1600422912988</Created><Modified>0</Modified><Origin>1600422912986695320</Origin><Readings><Id>0305fccf-e5bd-45ba-a0f5-5fe900244751</Id><Pushed>0</Pushed><Created>0</Created><Origin>1600422912</Origin><Modified>0</Modified><Device>Simple-Device02</Device><Name>echoString</Name><Value>HELLO</Value><ValueType>String</ValueType><FloatEncoding></FloatEncoding><BinaryValue></BinaryValue><MediaType></MediaType></Readings></Event>
```



As our device service keeps sending events with the value "HELLO" once every 5 seconds, we can see the XML messages. To test how the message changes, we can send commands to the device service via the core command service. 

Let's open a new terminal and:
```sh
# Let's store the IDs as same as earlier
$ DEVICE_ID=$(curl http://localhost:48082/api/v1/device -s | jq -r '.[] | select(.name == "Simple-Device02").id')
$ COMMAND_ID=$(curl http://localhost:48082/api/v1/device -s | jq -r '.[] | select(.name == "Simple-Device02").commands | .[] | .id')

# Query with the gathered IDs (the DEVICE_ID and COMMAND_ID variables are defined earlier):
$ curl http://localhost:48082/api/v1/device/$DEVICE_ID/command/$COMMAND_ID \
    -s \
    -X PUT \
    -H 'Content-Type: application/json' \
    -H 'cache-control: no-cache' \
    -d '{"echoString": "HELLO, WORLD"}'
```

Then, the app service will print out XML message with the new string "HELLO, WORLD" and this shows the new string went throught the services well.

<br/>

## Conclusion

In this tutorial, we prepared Ubuntu server 20.10 on RPI, launched EdgeX services, and created the custom device and app services. Although there are many unwritten details to keep it simple, now we know about the flow of EdgeX service development - where the important files are and how to build/test. To me, the benefit of using EdgeX is that all the queries get stored without concern of DB management and that is a huge plus if we deploy tons of edge devices everywhere. As we could see, running EdgeX on RPI is not difficult at all. Everything is ready there for us and the next exciting IoT projects!

<br/>

---

[To README](README.md)
