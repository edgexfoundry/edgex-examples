package main

import (
	"fmt"
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"

	azureTransforms "azure-export/pkg/transforms"
)

const (
	serviceKey           = "AzureExport"
	appConfigDeviceNames = "DeviceNames"
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

	// 2) shows how to access the application's specific simple configuration settings.
	deviceNames, err := service.GetAppSettingStrings(appConfigDeviceNames)
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}

	lc.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.

	// Load Azure-specific MQTT configuration from App SDK
	// You can also create AzureMQTTConfig struct yourself
	config, err := azureTransforms.LoadAzureMQTTConfig(service)

	if err != nil {
		lc.Error(fmt.Sprintf("Failed to load Azure MQTT configurations: %v\n", err))
		os.Exit(-1)
	}

	if err := service.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		azureTransforms.NewConversion().TransformToAzure,
		azureTransforms.NewAzureMQTTSender(service.LoggingClient, config).MQTTSend,
	); err != nil {
		lc.Error("SetFunctionsPipeline returned error: ", err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}
