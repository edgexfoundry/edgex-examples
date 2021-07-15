//
// Copyright (c) 2020 Intel Corporation
// Copyright (c) 2020 One Track Consulting
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
	"strconv"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

const (
	serviceKey = "app-json-logic"
)

var counter int

func main() {
	// turn off secure mode for examples. Not recommended for production
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an instance of the app service and initialize it.
	appService, ok := pkg.NewAppService(serviceKey)
	if !ok {
		appService.LoggingClient().Error("SDK initialization failed")
		os.Exit(-1)
	}

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := appService.GetAppSettingStrings("DeviceNames")
	if err != nil {
		appService.LoggingClient().Error(err.Error())
		os.Exit(-1)
	}
	appService.LoggingClient().Infof("Filtering for devices %v", deviceNames)
	listOfDevices := strings.Join(deviceNames, "\",\"")

	// Rule to look for a specific set of devices provided (similar to filterbydevicename)
	jsonlogicrule := "{ \"in\" : [{ \"var\" : \"deviceName\" }, [\"" + listOfDevices + "\"] ] }"

	// Rule to filter out a single device and allow all others
	//jsonlogicrule := "{ \"!\" : {\"in\" : [{ \"var\" : \"deviceName\" }, [\"" + listOfDevices + "\"] ] }}"

	// Rule to look for values > 0 in float readings -- MUST CONVERT TO READABLE FLOATS FIRST! Uncomment "ConvertToReadableFloatValues"
	//jsonlogicrule := "{ \"all\" : [ { \"var\" : \"readings\" } , {  \">\" : [ {\"var\":\"value\"}, 0 ] } ] }"

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	err = appService.SetFunctionsPipeline(
		// ConvertToReadableFloatValues, // Used for when looking at float values
		transforms.NewJSONLogic(jsonlogicrule).Evaluate,
		transforms.NewConversion().TransformToXML,
		printXMLToConsole,
	)

	if err != nil {
		appService.LoggingClient().Errorf("SetFunctionsPipeline returned error: %s", err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = appService.MakeItRun()
	if err != nil {
		appService.LoggingClient().Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func printXMLToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	if data == nil {
		// We didn't receive a result
		return false, nil
	}

	ds, ok := data.(string)

	if !ok {
		return false, fmt.Errorf("expected a string, got %T", data)
	}

	fmt.Println(ds)

	// Leverage the built in logging service in EdgeX
	ctx.LoggingClient().Debug("XML printed to console")

	ctx.SetResponseData([]byte(ds))
	return false, nil
}

var precision = 4

//ConvertToReadableFloatValues is used to facilitate writing jsonlogic rules
func ConvertToReadableFloatValues(ctx interfaces.AppFunctionContext, param interface{}) (bool, interface{}) {

	ctx.LoggingClient().Debug("Convert to Readable Float Values")

	if param == nil {
		// We didn't receive a result
		return false, nil
	}

	event := param.(models.Event)
	for index := range event.Readings {
		if eventReading, ok := event.Readings[index].(models.SimpleReading); ok {
			ctx.LoggingClient().Debugf("Event Reading for %s: %s is '%s'", event.DeviceName, eventReading.ResourceName, eventReading.Value)

			switch eventReading.ResourceName {
			case "Float32":
				f, err := strconv.ParseFloat(eventReading.Value, 32)
				if err != nil {
					return false, fmt.Errorf("unable to decode float32 value bytes: %s", err.Error())
				}

				eventReading.Value = strconv.FormatFloat(f, 'f', precision, 32)

			case "Float64":
				f, err := strconv.ParseFloat(eventReading.Value, 64)
				if err != nil {
					return false, fmt.Errorf("unable to decode float64 value bytes: %s", err.Error())
				}

				eventReading.Value = strconv.FormatFloat(f, 'f', precision, 64)
			}

			ctx.LoggingClient().Debugf("Converted Event Reading for %s: %s is '%s'", event.DeviceName, eventReading.ResourceName, eventReading.Value)
		}
	}

	return true, event
}
