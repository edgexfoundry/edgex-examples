//
// Copyright (c) 2020 Intel Corporation
// Copyright (c) 2021 One Track Consulting
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
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"

	"secrets/functions"
)

const (
	serviceKey = "app-secrets"
)

var counter int

func main() {

	// 1) First thing to do is to create an instance of the app service and initialize it.
	appService, ok := pkg.NewAppService(serviceKey)
	if !ok {
		appService.LoggingClient().Error("SDK initialization failed.")
		os.Exit(-1)
	}

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := appService.GetAppSettingStrings("DeviceNames")
	if err != nil {
		appService.LoggingClient().Error(err.Error())
		os.Exit(-1)
	}
	appService.LoggingClient().Infof("Filtering for devices %v", deviceNames)

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	err = appService.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		transforms.NewConversion().TransformToXML,
		functions.GetSecretsToConsole,
	)

	if err != nil {
		appService.LoggingClient().Errorf("SetFunctionsPipeline returned error: %s", err.Error())
	}

	// 4) Lastly, we'll go ahead and tell the service to "start" and begin listening for events
	// to trigger the pipeline.
	err = appService.MakeItRun()
	if err != nil {
		appService.LoggingClient().Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}
