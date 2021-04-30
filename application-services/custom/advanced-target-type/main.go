//
// Copyright (c) 2021 Intel Corporation
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

	"advanced-target-type/functions"
)

const (
	serviceKey = "advancedTargetType"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an new instance of an EdgeX Application Service using
	//    custom Trigger Type
	service, ok := pkg.NewAppServiceWithTargetType(serviceKey, &functions.Person{})
	if !ok {
		os.Exit(-1)
	}

	// Leverage the built in logging service in EdgeX
	lc := service.LoggingClient()

	// 2) This is our functions pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		functions.FormatPhoneDisplay,                 // Expects a Person as set by TargetType
		functions.ConvertToXML,                       // Expects a Person
		functions.PrintXMLToConsole,                  // Expects XML string
		transforms.NewResponseData().SetResponseData, // Expects string or []byte. Returns XML formatted Person with PhoneDisplay set sent as the trigger response
	); err != nil {
		lc.Errorf("Setting Functions Pipeline failed: %s", err.Error())
		os.Exit(-1)
	}

	// 3) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for Persons
	// to trigger the pipeline.
	if err := service.MakeItRun(); err != nil {
		lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}
