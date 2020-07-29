package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/util"

	awsTransforms "aws-export/pkg/transforms"
)

const (
	serviceKey = "AWSExport"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an instance of the EdgeX SDK and initialize it.
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		os.Exit(-1)
	}

	// 2) Load AWS-specific MQTT configuration from App SDK
	// You can also create AWSMQTTConfig struct yourself
	config, err := awsTransforms.LoadAWSMQTTConfig(edgexSdk)
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("Failed to load AWS MQTT configurations: %v\n", err))
		os.Exit(-1)
	}

	// 3) Get DeviceNameFilter from config
	deviceNamesCleaned := util.DeleteEmptyAndTrim(strings.FieldsFunc(config.DeviceNames, util.SplitComma))
	edgexSdk.LoggingClient.Debug(fmt.Sprintf("Device names read %s\n", deviceNamesCleaned))

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(deviceNamesCleaned).FilterByDeviceName,
		awsTransforms.NewConversion().TransformToAWS,
		printAWSDataToConsole,
		awsTransforms.NewAWSMQTTSender(edgexSdk.LoggingClient, config).MQTTSend,
	)

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func printAWSDataToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	fmt.Println(params[0].(string))

	// Leverage the built in logging service in EdgeX
	edgexcontext.LoggingClient.Debug("Printed to console")

	edgexcontext.Complete([]byte(params[0].(string)))
	return false, nil

}
