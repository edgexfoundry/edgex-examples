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

package transforms

import (
	"errors"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
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
	if data == nil {
		return false, errors.New("No Event Received")
	}
	ctx.LoggingClient().Debug("Transforming to CloudEvent")
	event, ok := data.(models.Event)
	if !ok {
		return false, errors.New("Unexpected type received")
	}
	if len(event.Readings) == 0 {
		return false, errors.New("No event readings to transform")
	}

	var cloudEvents []cloudevents.Event
	for _, reading := range event.Readings {
		cloudevent := cloudevents.NewEvent(cloudevents.VersionV1)
		baseReading := reading.GetBaseReading()
		cloudevent.SetID(baseReading.Id)
		cloudevent.SetType(baseReading.ResourceName)
		cloudevent.SetSource(event.DeviceName)
		unixTime := time.Unix(0, baseReading.Origin) // assuming time is formatted as nanoseconds if seconds time.Unix(reading.Origin, 0)
		timeRFC3339, err := time.Parse(
			time.RFC3339,
			unixTime.Format(time.RFC3339))
		if err != nil {
			return false, fmt.Errorf("Failed to parse reading.Origin as RFC3339 time: %d, %s", baseReading.Origin, err)
		}
		cloudevent.SetTime(timeRFC3339)
		// extension names are always lowercase
		cloudevent.SetExtension("eventid", event.Id)
		cloudevent.SetExtension("valuetype", baseReading.ValueType)

		switch r := reading.(type) {
		case models.SimpleReading:
			cloudevent.SetDataContentType(common.ContentTypeJSON)
			if err := cloudevent.SetData(common.ContentTypeJSON, r.Value); err != nil {
				return false, fmt.Errorf("Error setting data field for cloud event, %s", err)
			}
		case models.BinaryReading:
			cloudevent.SetData("", r.BinaryValue)
		default:
			return false, fmt.Errorf("Unknown reading type: %T", r)
		}

		cloudEvents = append(cloudEvents, cloudevent)
	}
	return true, cloudEvents
}

// TransformFromCloudEvent will transform a Cloud Event to an Edgex models.Event
// It will return an error and stop the pipeline if a non-edgex event is received or if no data is received.
func (f Conversion) TransformFromCloudEvent(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, stringType interface{}) {
	if data == nil {
		return false, errors.New("No Event Received")
	}
	ctx.LoggingClient().Debug("Transforming from CloudEvent to models.Event")
	result, ok := data.([]cloudevents.Event)
	if !ok {
		return false, errors.New("Unexpected type received")
	}
	event := models.Event{}
	for _, cloudevent := range result {
		event.DeviceName = cloudevent.Source()

		baseReading := models.BaseReading{
			Id:           cloudevent.ID(),
			Origin:       cloudevent.Time().Unix(),
			DeviceName:   cloudevent.Source(),
			ResourceName: cloudevent.Type(),
		}

		extensions := cloudevent.Extensions()

		if val, ok := extensions["eventid"]; ok {
			event.Id = val.(string)
		}
		if val, ok := extensions["valuetype"]; ok {
			baseReading.ValueType = val.(string)
		}

		var reading models.Reading

		if cloudevent.DataBase64 {
			reading = models.BinaryReading{BaseReading: baseReading, BinaryValue: cloudevent.Data()}
		} else {
			sr := models.SimpleReading{BaseReading: baseReading, Value: ""}
			tempStr := ""

			if err := cloudevent.DataAs(&tempStr); err != nil {
				panic(err)
			} else {
				sr.Value = tempStr
			}

			reading = sr
		}

		event.Readings = append(event.Readings, reading)
	}
	if len(event.Readings) == 0 {
		return false, errors.New("Event has no readings")
	}
	return true, event
}
