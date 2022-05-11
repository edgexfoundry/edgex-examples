//
// Copyright (c) 2019 Intel Corporation
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
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"os"

	"advanced-target-type/functions"
)

const (
	serviceKey = "app-target-type"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	// 1) First thing to do is to create an instance of an edgex service with your TargetType set
	//    and initialize it. Note that the TargetType is a pointer to an instance of the type.
	appService, ok := pkg.NewAppServiceWithTargetType(serviceKey, &functions.Person{})
	if !ok {
		appService.LoggingClient().Errorf("App Service initialization failed for %s", serviceKey)
		os.Exit(-1)
	}

	// 2) This is our functions pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	err := appService.SetFunctionsPipeline(
		functions.FormatPhoneDisplay,                 // Expects a Person as set by TargetType
		functions.ConvertToXML,                       // Expects a Person
		functions.PrintXmlToConsole,                  // Expects XML string
		transforms.NewResponseData().SetResponseData, // Expects string or []byte. Returns XML formatted Person with PhoneDisplay set sent as the trigger response
	)

	if err != nil {
		appService.LoggingClient().Errorf("Setting Functions Pipeline failed: %s" + err.Error())
		os.Exit(-1)
	}

	// 3) Lastly, we'll go ahead and tell the service to "start" and begin listening for Persons
	// to trigger the pipeline.
	err = appService.MakeItRun()
	if err != nil {
		appService.LoggingClient().Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}
