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
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"
)

const (
	serviceKey = "json-logic-example"
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
	deviceNames, err := service.GetAppSettingStrings("DeviceNames")
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}

	lc.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))
	listOfDevices := strings.Join(deviceNames, "\",\"")

	// Rule to look for a specific set of devices provided (similar to filterbydevicename)
	jsonlogicrule := "{ \"in\" : [{ \"var\" : \"device\" }, [\"" + listOfDevices + "\"] ] }"

	// Rule to filter out a single device and allow all others
	//jsonlogicrule := "{ \"!\" : {\"in\" : [{ \"var\" : \"device\" }, [\"" + listOfDevices + "\"] ] }}"

	// Rule to look for values > 0 in float readings -- MUST CONVERT TO READABLE FLOATS FIRST! Uncomment "ConvertToReadableFloatValues"
	//jsonlogicrule := "{ \"all\" : [ { \"var\" : \"readings\" } , {  \">\" : [ {\"var\":\"value\"}, 0 ] } ] }"

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		// ConvertToReadableFloatValues, // Used for when looking at float values
		transforms.NewJSONLogic(jsonlogicrule).Evaluate,
		transforms.NewConversion().TransformToXML,
		printXMLToConsole,
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

func printXMLToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Leverage the built in logging service in EdgeX
	lc := ctx.LoggingClient()

	if data == nil {
		return false, errors.New("printXMLToConsole: No data received")
	}

	xml, ok := data.(string)
	if !ok {
		return false, errors.New("printXMLToConsole: Data received is not the expected 'string' type")
	}

	lc.Debug(xml)
	ctx.SetResponseData([]byte(xml))
	return true, xml
}

var precision = 4

func ConvertToReadableFloatValues(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()

	lc.Debug("Convert to Readable Float Values")

	if data == nil {
		return false, errors.New("ConvertToReadableFloatValues: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("ConvertToReadableFloatValues: Data received is not the expected 'Event' type")
	}

	for index := range event.Readings {
		eventReading := &event.Readings[index]
		lc.Debug(fmt.Sprintf("Event Reading for %s: %s is '%s'", event.DeviceName, eventReading.ResourceName, eventReading.Value))

		data, err := base64.StdEncoding.DecodeString(eventReading.Value)
		if err != nil {
			return false, fmt.Errorf("unable to Base 64 decode float32/64 value ('%s'): %s", eventReading.Value, err.Error())
		}

		switch eventReading.ResourceName {
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

		lc.Debug(fmt.Sprintf("Converted Event Reading for %s: %s is '%s'", event.DeviceName, eventReading.ResourceName, eventReading.Value))
	}

	return true, event
}
