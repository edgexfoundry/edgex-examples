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
	"encoding/base64"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/coredata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/stretchr/testify/assert"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/urlclient"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

var context *appcontext.Context
var lc logger.LoggingClient

const (
	devID1 = "id1"
	devID2 = "id2"
)

func init() {
	lc := logger.NewClient("app_functions_sdk_go", false, "./test.log", "DEBUG")
	eventClient := coredata.NewEventClient(
		urlclient.New(nil, nil, nil, "", "", 0, "http://test"+clients.ApiEventRoute),
	)

	context = &appcontext.Context{
		LoggingClient: lc,
		EventClient:   eventClient,
	}
}

func TestTransformToCloudEvent(t *testing.T) {
	type dataTest struct {
		valueType string
		want      string
		event     models.Event
	}

	valueTypesToValues := map[string]string{
		models.ValueTypeBool:    "true",
		models.ValueTypeString:  "This is a string",
		models.ValueTypeUint8:   "123",
		models.ValueTypeUint16:  "123",
		models.ValueTypeUint32:  "123",
		models.ValueTypeUint64:  "123",
		models.ValueTypeInt8:    "-123",
		models.ValueTypeInt16:   "-123",
		models.ValueTypeInt32:   "-123",
		models.ValueTypeInt64:   "-123",
		models.ValueTypeFloat32: "1.23",
		models.ValueTypeFloat64: "1.23",
		models.ValueTypeBinary:  "will be converted to bytes in test",
	}

	testCases := map[string]dataTest{}
	for valueType, value := range valueTypesToValues {

		var dataBytes []byte
		var want string
		if valueType == models.ValueTypeBinary {
			dataBytes = []byte(value)
			want = `{"data_base64":"` + base64.StdEncoding.EncodeToString(dataBytes) + `","eventid":"event-` + devID1 + `","floatencoding":"","id":"123-abc","source":"id1","specversion":"1.0","time":"1970-01-01T00:00:00Z","type":"test-reading","valuetype":"` + valueType + `"}`
		} else {
			want = `{"data":"` + value + `","datacontenttype":"application/json","eventid":"event-` + devID1 + `","floatencoding":"","id":"123-abc","source":"id1","specversion":"1.0","time":"1970-01-01T00:00:00Z","type":"test-reading","valuetype":"` + valueType + `"}`
		}

		testCases["test to cloudevent: "+valueType] = dataTest{
			event: models.Event{
				ID:     "event-" + devID1,
				Device: devID1,
				Readings: []models.Reading{models.Reading{
					Id:          "123-abc",
					Name:        "test-reading",
					Value:       value,
					ValueType:   valueType,
					BinaryValue: dataBytes}},
			},
			valueType: valueType,
			want:      want,
		}
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			conv := NewConversion()
			continuePipeline, result := conv.TransformToCloudEvent(context, tc.event)
			resultBytes := result.([][]byte)
			resultString := string(resultBytes[0])
			assert.NotNil(t, result)
			assert.True(t, continuePipeline)
			assert.Equal(t, tc.want, resultString)
		})
	}
}

func TestTransformToCloudEventWrongEvent(t *testing.T) {
	eventIn := "Not a models.Event, a string"
	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, eventIn)
	assert.Error(t, result.(error))
	assert.False(t, continuePipeline)
}

func TestTransformToCloudEventMultipleEvents(t *testing.T) {
	value1 := "1234"
	value2 := "4321"
	valueType := models.ValueTypeInt64
	want := []string{
		`{"data":"` + value1 + `","datacontenttype":"application/json","eventid":"event-` + devID1 + `","floatencoding":"","id":"123-abc","source":"id1","specversion":"1.0","time":"1970-01-01T00:00:00Z","type":"test-reading","valuetype":"` + valueType + `"}`,
		`{"data":"` + value2 + `","datacontenttype":"application/json","eventid":"event-` + devID1 + `","floatencoding":"","id":"123-abc","source":"id1","specversion":"1.0","time":"1970-01-01T00:00:00Z","type":"test-reading","valuetype":"` + valueType + `"}`}
	event := models.Event{
		ID:     "event-" + devID1,
		Device: devID1,
		Readings: []models.Reading{
			models.Reading{
				Id:        "123-abc",
				Name:      "test-reading",
				Value:     value1,
				ValueType: valueType},
			models.Reading{
				Id:        "123-abc",
				Name:      "test-reading",
				Value:     value2,
				ValueType: valueType}}}

	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, event)
	assert.NotNil(t, result)
	assert.True(t, continuePipeline)
	resultBytes := result.([][]byte)
	for i, cloudEvent := range resultBytes {
		resultString := string(cloudEvent)
		assert.Equal(t, want[i], resultString)
	}
}

func TestTransformToCloudEventNoReadings(t *testing.T) {
	eventIn := models.Event{
		ID:     "event-" + devID1,
		Device: devID1,
	}
	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, eventIn)
	assert.Error(t, result.(error))
	assert.False(t, continuePipeline)
}

func TestTransformToCloudEventNoEvent(t *testing.T) {
	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context)
	assert.Error(t, result.(error))
	assert.False(t, continuePipeline)
}

func TestTransformFromCloudEvent(t *testing.T) {
	value := "4321"
	valueType := models.ValueTypeInt64
	singleCloudEvent := []byte(`{"data":"` + value + `","datacontenttype":"application/json","floatencoding":"", "eventid":"` + "event-" + devID1 + `","id":"123-abc","source":"id1","specversion":"1.0","time":"1970-01-01T00:00:00Z","type":"test-reading","valuetype":"` + valueType + `"}`)
	var cloudEvents [][]byte
	cloudEvents = append(cloudEvents, singleCloudEvent)
	expectedEvent := models.Event{
		ID:     "event-" + devID1,
		Device: devID1,
		Readings: []models.Reading{
			models.Reading{
				Id:        "123-abc",
				Name:      "test-reading",
				Value:     value,
				ValueType: valueType}}}

	conv := NewConversion()
	continuePipeline, result := conv.TransformFromCloudEvent(context, cloudEvents)
	edgexEvent, ok := result.(models.Event)
	assert.NotNil(t, result)
	assert.True(t, ok)
	assert.True(t, continuePipeline)
	assert.Equal(t, expectedEvent.ID, edgexEvent.ID)
	assert.Equal(t, expectedEvent.Readings[0].Id, edgexEvent.Readings[0].Id)
	assert.Equal(t, expectedEvent.Readings[0].Name, edgexEvent.Readings[0].Name)
	assert.Equal(t, expectedEvent.Readings[0].Value, edgexEvent.Readings[0].Value)
}

func TestTransformFromCloudEventEmptyCloudEvent(t *testing.T) {
	var cloudEvent cloudevents.Event
	conv := NewConversion()
	continuePipeline, result := conv.TransformFromCloudEvent(context, cloudEvent)
	_, ok := result.(error)

	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.False(t, continuePipeline)
}

func TestTransformFromCloudEventEmptyReadins(t *testing.T) {
	singleCloudEvent := []byte(`{}`)
	var cloudEvents [][]byte
	cloudEvents = append(cloudEvents, singleCloudEvent)
	conv := NewConversion()
	continuePipeline, result := conv.TransformFromCloudEvent(context, cloudEvents)
	_, ok := result.(error)
	assert.NotNil(t, result)
	assert.True(t, ok)
	assert.False(t, continuePipeline)
}
