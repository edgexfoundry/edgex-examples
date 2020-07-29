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

package transforms

import (
	"errors"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
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
func (f Conversion) TransformToCloudEvent(edgexcontext *appcontext.Context, params ...interface{}) (continuePipeline bool, stringType interface{}) {
	if len(params) < 1 {
		return false, errors.New("No Event Received")
	}
	edgexcontext.LoggingClient.Debug("Transforming to CloudEvent")
	event, ok := params[0].(models.Event)
	if !ok {
		return false, errors.New("Unexpected type received")
	}
	if len(event.Readings) == 0 {
		return false, errors.New("No event readings to transform")
	}

	var cloudeventReadings []cloudevents.Event
	for _, reading := range event.Readings {
		cloudevent := cloudevents.NewEvent(cloudevents.VersionV1)
		cloudevent.SetID(reading.Id)
		cloudevent.SetType(reading.Name)
		cloudevent.SetSource(event.Device)
		unixTime := time.Unix(0, reading.Origin) // assuming time is formatted as nanoseconds if seconds time.Unix(reading.Origin, 0)
		timeRFC3339, err := time.Parse(
			time.RFC3339,
			unixTime.Format(time.RFC3339))
		if err != nil {
			return false, fmt.Errorf("Failed to parse reading.Origin as RFC3339 time: %d, %s", reading.Origin, err)
		}
		cloudevent.SetTime(timeRFC3339)
		// extension names are always lowercase
		cloudevent.SetExtension("eventid", event.ID)
		cloudevent.SetExtension("valuetype", reading.ValueType)
		cloudevent.SetExtension("floatencoding", reading.FloatEncoding)
		if len(reading.BinaryValue) > 0 {
			// if reading.BinaryValue is set that becomes data and reading.Value is ignored
			if err := cloudevent.SetData(reading.BinaryValue); err != nil {
				return false, fmt.Errorf("Error setting data field for cloud event, %s", err)
			}
		} else {
			cloudevent.SetDataContentType("application/json")
			if err := cloudevent.SetData(reading.Value); err != nil {
				return false, fmt.Errorf("Error setting data field for cloud event, %s", err)
			}
		}
		cloudeventReadings = append(cloudeventReadings, cloudevent)
	}
	return true, cloudeventReadings
}

// TransformFromCloudEvent will transform a Cloud Event to an Edgex models.Event
// It will return an error and stop the pipeline if a non-edgex event is received or if no data is received.
func (f Conversion) TransformFromCloudEvent(edgexcontext *appcontext.Context, params ...interface{}) (continuePipeline bool, stringType interface{}) {
	if len(params) < 1 {
		return false, errors.New("No Event Received")
	}
	edgexcontext.LoggingClient.Debug("Transforming from CloudEvent to models.Event")
	result, ok := params[0].([]cloudevents.Event)
	if !ok {
		return false, errors.New("Unexpected type received")
	}
	event := models.Event{}
	for _, cloudevent := range result {
		event.Device = cloudevent.Source()
		reading := models.Reading{}
		reading.Id = cloudevent.ID()
		reading.Name = cloudevent.Type()
		reading.Origin = cloudevent.Time().Unix()
		extensions := cloudevent.Extensions()
		if val, ok := extensions["eventid"]; ok {
			event.ID = val.(string)
		}
		if val, ok := extensions["valuetype"]; ok {
			reading.ValueType = val.(string)
		}
		if val, ok := extensions["floatencoding"]; ok {
			reading.FloatEncoding = val.(string)
		}
		dataBytes, dataBinary := (cloudevent.Data).([]byte)
		if cloudevent.DataBinary && dataBinary {
			reading.BinaryValue = dataBytes
		} else {
			if err := cloudevent.DataAs(&reading.Value); err != nil {
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
