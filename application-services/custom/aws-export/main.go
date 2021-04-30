package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/util"

	awsTransforms "aws-export/pkg/transforms"
)

const (
	serviceKey = "AWSExport"
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

	// 2) Load AWS-specific MQTT configuration from App SDK
	// You can also create AWSMQTTConfig struct yourself
	config, err := awsTransforms.LoadAWSMQTTConfig(service)
	if err != nil {
		lc.Error(fmt.Sprintf("Failed to load AWS MQTT configurations: %v\n", err))
		os.Exit(-1)
	}

	// 3) Get DeviceNameFilter from config
	deviceNamesCleaned := util.DeleteEmptyAndTrim(strings.FieldsFunc(config.DeviceNames, util.SplitComma))
	lc.Debug(fmt.Sprintf("Device names read %s\n", deviceNamesCleaned))

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	service.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNamesCleaned).FilterByDeviceName,
		awsTransforms.NewConversion().TransformToAWS,
		printAWSDataToConsole,
		awsTransforms.NewAWSMQTTSender(lc, config).MQTTSend,
	)

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func printAWSDataToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Leverage the built in logging service in EdgeX
	lc := ctx.LoggingClient()

	if data == nil {
		return false, errors.New("printAWSDataToConsole: No data received")
	}

	fmt.Println(data.(string))

	// Leverage the built in logging service in EdgeX
	lc.Debug("Printed to console")

	ctx.SetResponseData([]byte(data.(string)))
	return false, nil
}
