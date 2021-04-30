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

package transforms

import (
	"errors"
	"fmt"
	"time"

	cloud "github.com/cloudevents/sdk-go"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"
)

// Conversion houses various conversion
type Conversion struct {
}

// NewConversion creates, initializes and returns a new instance of Conversion
func NewConversion() Conversion {
	return Conversion{}
}

// TransformToCloudEvent will transform a models.Event to a Cloud Event
// It will return an error and stop the pipeline if a non-edgex event is received or if no data is received.
func (f Conversion) TransformToCloudEvent(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, stringType interface{}) {
	lc := ctx.LoggingClient()

	lc.Debug("Transforming to CloudEvent")

	if data == nil {
		return false, errors.New("TransformToCloudEvent: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("TransformToCloudEvent: didn't receive expect Event type")
	}

	if len(event.Readings) == 0 {
		return false, errors.New("TransformToCloudEvent: No event readings to transform")
	}

	var cloudEventReadings []cloud.Event
	for _, reading := range event.Readings {
		cloudEvent := cloud.NewEvent(cloud.VersionV1)
		cloudEvent.SetID(reading.Id)
		cloudEvent.SetType(reading.ResourceName)
		cloudEvent.SetSource(event.DeviceName)
		unixTime := time.Unix(0, reading.Origin) // assuming time is formatted as nanoseconds if seconds time.Unix(reading.Origin, 0)
		timeRFC3339, err := time.Parse(
			time.RFC3339,
			unixTime.Format(time.RFC3339))
		if err != nil {
			return false, fmt.Errorf("Failed to parse reading.Origin as RFC3339 time: %d, %s", reading.Origin, err)
		}
		cloudEvent.SetTime(timeRFC3339)
		// extension names are always lowercase
		cloudEvent.SetExtension("eventid", event.Id)
		cloudEvent.SetExtension("eventsource", event.SourceName)
		cloudEvent.SetExtension("valuetype", reading.ValueType)
		cloudEvent.SetExtension("profilename", reading.ProfileName)

		if len(reading.BinaryValue) > 0 {
			// if reading.BinaryValue is set that becomes data and reading.Value is ignored
			if err := cloudEvent.SetData(reading.BinaryValue); err != nil {
				return false, fmt.Errorf("Error setting data field for cloud event, %s", err)
			}
		} else {
			cloudEvent.SetDataContentType("application/json")
			if err := cloudEvent.SetData(reading.Value); err != nil {
				return false, fmt.Errorf("Error setting data field for cloud event, %s", err)
			}
		}
		cloudEventReadings = append(cloudEventReadings, cloudEvent)
	}
	return true, cloudEventReadings
}

// TransformFromCloudEvent will transform a Cloud Event to an Edgex models.Event
// It will return an error and stop the pipeline if a non-edgex event is received or if no data is received.
func (f Conversion) TransformFromCloudEvent(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, stringType interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("Transforming from Cloud Event to Edgex Event")

	if data == nil {
		return false, errors.New("TransformFromCloudEvent: No data received")
	}

	cloudEvents, ok := data.([]cloud.Event)
	if !ok {
		return false, errors.New("TransformFromCloudEvent: didn't receive expect '[]cloud.Event' type")
	}

	event := dtos.NewEvent("TBD", "TBD", "TBD")
	for _, cloudEvent := range cloudEvents {
		event.DeviceName = cloudEvent.Source()
		reading := dtos.BaseReading{}
		reading.Id = cloudEvent.ID()
		reading.ResourceName = cloudEvent.Type()
		reading.Origin = cloudEvent.Time().Unix()
		extensions := cloudEvent.Extensions()
		if val, ok := extensions["eventid"]; ok {
			event.Id = val.(string)
		}
		if val, ok := extensions["eventsource"]; ok {
			event.SourceName = val.(string)
		}
		if val, ok := extensions["valuetype"]; ok {
			reading.ValueType = val.(string)
		}
		if val, ok := extensions["profilename"]; ok {
			reading.ProfileName = val.(string)
			event.ProfileName = val.(string)
		}

		dataBytes, dataBinary := (cloudEvent.Data).([]byte)
		if cloudEvent.DataBinary && dataBinary {
			reading.BinaryValue = dataBytes
		} else {
			if err := cloudEvent.DataAs(&reading.Value); err != nil {
				return false, fmt.Errorf("Can't unmarshal cloud event data, %s", err)
			}
		}
		event.Readings = append(event.Readings, reading)
	}
	if len(event.Readings) == 0 {
		return false, errors.New("Event has no readings")
	}
	return true, event
}
