//
// Copyright (c) 2020 Intel Corporation
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
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	serviceKey = "json-logic-example"
)

var counter int

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
	listOfDevices := strings.Join(deviceNames, "\",\"")

	// Rule to look for a specific set of devices provided (similar to filterbydevicename)
	jsonlogicrule := "{ \"in\" : [{ \"var\" : \"device\" }, [\"" + listOfDevices + "\"] ] }"

	// Rule to filter out a single device and allow all others
	//jsonlogicrule := "{ \"!\" : {\"in\" : [{ \"var\" : \"device\" }, [\"" + listOfDevices + "\"] ] }}"

	// Rule to look for values > 0 in float readings -- MUST CONVERT TO READABLE FLOATS FIRST! Uncomment "ConvertToReadableFloatValues"
	//jsonlogicrule := "{ \"all\" : [ { \"var\" : \"readings\" } , {  \">\" : [ {\"var\":\"value\"}, 0 ] } ] }"

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	edgexSdk.SetFunctionsPipeline(
		// ConvertToReadableFloatValues, // Used for when looking at float values
		transforms.NewJSONLogic(jsonlogicrule).Evaluate,
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

func printXMLToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	fmt.Println(params[0].(string))

	// Leverage the built in logging service in EdgeX
	edgexcontext.LoggingClient.Debug("XML printed to console")

	edgexcontext.Complete([]byte(params[0].(string)))
	return false, nil
}

var precision = 4

func ConvertToReadableFloatValues(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	edgexcontext.LoggingClient.Debug("Convert to Readable Float Values")

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	event := params[0].(models.Event)
	for index := range event.Readings {
		eventReading := &event.Readings[index]
		edgexcontext.LoggingClient.Debug(fmt.Sprintf("Event Reading for %s: %s is '%s'", event.Device, eventReading.Name, eventReading.Value))

		data, err := base64.StdEncoding.DecodeString(eventReading.Value)
		if err != nil {
			return false, fmt.Errorf("unable to Base 64 decode float32/64 value ('%s'): %s", eventReading.Value, err.Error())
		}

		switch eventReading.Name {
		case "RandomValue_Float32":
			var value float32
			err = binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			if err != nil {
				return false, fmt.Errorf("unable to decode float32 value bytes: %s", err.Error())
			}

			eventReading.Value = strconv.FormatFloat(float64(value), 'f', precision, 32)

		case "RandomValue_Float64":
			var value float64
			err := binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			if err != nil {
				return false, fmt.Errorf("unable to decode float64 value bytes: %s", err.Error())
			}

			eventReading.Value = strconv.FormatFloat(value, 'f', precision, 64)
		}

		edgexcontext.LoggingClient.Debug(fmt.Sprintf("Converted Event Reading for %s: %s is '%s'", event.Device, eventReading.Name, eventReading.Value))
	}

	return true, event
}
