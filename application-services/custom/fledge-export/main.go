package main

import (
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"

	fledgeTransforms "fledge-export/pkg/transforms"
)

const serviceKey = "app-fledge-export"

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

	// 2) Shows how to access the application's specific configuration settings.
	fledgeEndpoint, err := service.GetAppSetting("FledgeSouthHTTPEndpoint")
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		fledgeTransforms.NewConversion().TransformToFledge,
		transforms.NewHTTPSender(fledgeEndpoint, "application/json", false).HTTPPost,
	); err != nil {
		lc.Errorf("SDK SetPipeline failed: %s\n", err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}
