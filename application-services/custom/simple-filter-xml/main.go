//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
)

const (
	serviceKey = "sampleFilterXml"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an instance of the EdgeX SDK and initialize it.
	service, ok := pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(-1)
	}

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := service.GetAppSettingStrings("DeviceNames")
	if err != nil {
		service.LoggingClient().Error(err.Error())
		os.Exit(-1)
	}
	service.LoggingClient().Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	service.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		transforms.NewConversion().TransformToXML,
		printXMLToConsole,
	)

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		service.LoggingClient().Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func printXMLToConsole(appContext interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	if data == nil {
		// We didn't receive a result
		return false, nil
	}

	fmt.Println(data.(string))

	// Leverage the built in logging service in EdgeX
	appContext.LoggingClient().Debug("XML printed to console")

	appContext.SetResponseData([]byte(data.(string)))
	return false, nil
}
