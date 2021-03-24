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

package functions

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"
)

var precision = 4

func ConvertToReadableFloatValues(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("Convert to Readable Float Values")

	if data == nil {
		return false, errors.New("ConvertToReadableFloatValues: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("ConvertToReadableFloatValues: didn't receive expect Event type")
	}

	for index := range event.Readings {
		eventReading := &event.Readings[index]
		lc.Debugf("Event Reading for %s: %s is '%s'", event.DeviceName, eventReading.ResourceName, eventReading.Value)

		// TODO: Change for E-Notation Only
		data, err := base64.StdEncoding.DecodeString(eventReading.Value)
		if err != nil {
			return false, fmt.Errorf("unable to Base 64 decode float32/64 value ('%s'): %s", eventReading.Value, err.Error())
		}

		switch eventReading.ResourceName {
		case "Float32":
			var value float32
			err = binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			if err != nil {
				return false, fmt.Errorf("unable to decode float32 value bytes: %s", err.Error())
			}

			eventReading.Value = strconv.FormatFloat(float64(value), 'f', precision, 32)

		case "Float64":
			var value float64
			err := binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			if err != nil {
				return false, fmt.Errorf("unable to decode float64 value bytes: %s", err.Error())
			}

			eventReading.Value = strconv.FormatFloat(value, 'f', precision, 64)
		}

		lc.Debugf("Converted Event Reading for %s: %s is '%s'", event.DeviceName, eventReading.ResourceName, eventReading.Value)
	}

	return true, event
}
