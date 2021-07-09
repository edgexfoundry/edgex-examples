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
	"encoding/base64"
	"encoding/json"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/stretchr/testify/assert"
)

var context interfaces.AppFunctionContext
var lc logger.LoggingClient

const (
	devID1 = "id1"
	devID2 = "id2"
)

func init() {
	lc := logger.NewMockClient()

	context = pkg.NewAppFuncContextForTest(uuid.NewString(), lc)
}

func TestTransformToCloudEvent(t *testing.T) {
	type dataTest struct {
		valueType string
		want      string
		event     models.Event
	}

	valueTypesToValues := map[string]string{
		common.ValueTypeBool:    "true",
		common.ValueTypeString:  "This is a string",
		common.ValueTypeUint8:   "123",
		common.ValueTypeUint16:  "123",
		common.ValueTypeUint32:  "123",
		common.ValueTypeUint64:  "123",
		common.ValueTypeInt8:    "-123",
		common.ValueTypeInt16:   "-123",
		common.ValueTypeInt32:   "-123",
		common.ValueTypeInt64:   "-123",
		common.ValueTypeFloat32: "1.23",
		common.ValueTypeFloat64: "1.23",
		common.ValueTypeBinary:  "will be converted to bytes in test",
	}

	testCases := map[string]dataTest{}
	for valueType, value := range valueTypesToValues {

		var dataBytes []byte
		var want string

		var reading models.Reading

		if valueType == common.ValueTypeBinary {
			dataBytes = []byte(value)

			reading = models.BinaryReading{
				BaseReading: models.BaseReading{
					Id:           "123-abc",
					ResourceName: "test-reading",
					ValueType:    valueType,
				},
				BinaryValue: dataBytes,
			}

			want = `{"specversion":"1.0","id":"123-abc","source":"id1","type":"test-reading","time":"1970-01-01T00:00:00Z","data_base64":"` + base64.StdEncoding.EncodeToString(dataBytes) // + `","valuetype":"` + valueType  +  `","eventid":"event-` + devID1  + `"}`
		} else {
			reading = models.SimpleReading{
				BaseReading: models.BaseReading{
					Id:           "123-abc",
					ResourceName: "test-reading",
					ValueType:    valueType,
				},
				Value: value,
			}
			want = `{"specversion":"1.0","id":"123-abc","source":"id1","type":"test-reading","datacontenttype":"application/json","time":"1970-01-01T00:00:00Z","data":"` + value // + `","eventid":"event-` + devID1 + `","valuetype":"` + valueType + `"}`
		}

		testCases["test to cloudevent: "+valueType] = dataTest{
			event: models.Event{
				Id:         "event-" + devID1,
				DeviceName: devID1,
				Readings: []models.Reading{
					reading,
				},
			},
			valueType: valueType,
			want:      want,
		}
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			conv := NewConversion()
			continuePipeline, result := conv.TransformToCloudEvent(context, tc.event)
			resultEvent, ok := result.([]cloudevents.Event)

			require.True(t, ok)

			assert.NotNil(t, resultEvent)
			assert.True(t, continuePipeline)

			have, err := resultEvent[0].MarshalJSON()

			require.NoError(t, err)

			// this is kind of wild - cloudevents marshaling seems to juxtapose valuetype and eventid at times here, so only checking up to the correct value being added
			assert.True(t, strings.HasPrefix(string(have), tc.want))
		})
	}
}

func TestTransformToCloudEventWrongEvent(t *testing.T) {
	eventIn := "Not a common.Event, a string"
	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, eventIn)
	assert.Error(t, result.(error))
	assert.False(t, continuePipeline)
}

func TestTransformToCloudEventMultipleEvents(t *testing.T) {
	value1 := "1234"
	value2 := "4321"
	valueType := common.ValueTypeInt64
	want := []string{
		`{"specversion":"1.0","id":"123-abc","source":"id1","type":"test-reading","datacontenttype":"application/json","time":"1970-01-01T00:00:00Z","data":"` + value1 + `","eventid":"event-` + devID1 + `","valuetype":"` + valueType + `"}`,
		`{"specversion":"1.0","id":"123-abc","source":"id1","type":"test-reading","datacontenttype":"application/json","time":"1970-01-01T00:00:00Z","data":"` + value2 + `","eventid":"event-` + devID1 + `","valuetype":"` + valueType + `"}`,
	}

	event := models.Event{
		Id:         "event-" + devID1,
		DeviceName: devID1,
		Readings: []models.Reading{
			models.SimpleReading{
				BaseReading: models.BaseReading{
					Id:           "123-abc",
					ResourceName: "test-reading",
					ValueType:    valueType,
				},
				Value: value1,
			},
			models.SimpleReading{
				BaseReading: models.BaseReading{
					Id:           "123-abc",
					ResourceName: "test-reading",
					ValueType:    valueType,
				},
				Value: value2,
			},
		},
	}

	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, event)
	assert.NotNil(t, result)
	assert.True(t, continuePipeline)
	cloudEvents := result.([]cloudevents.Event)
	for i, cloudEvent := range cloudEvents {
		resultBytes, err := json.Marshal(cloudEvent)
		require.NoError(t, err)
		assert.Equal(t, want[i], string(resultBytes))
	}
}

func TestTransformToCloudEventNoReadings(t *testing.T) {
	eventIn := models.Event{
		Id:         "event-" + devID1,
		DeviceName: devID1,
	}
	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, eventIn)
	assert.Error(t, result.(error))
	assert.False(t, continuePipeline)
}

func TestTransformToCloudEventNoEvent(t *testing.T) {
	conv := NewConversion()
	continuePipeline, result := conv.TransformToCloudEvent(context, nil)
	assert.Error(t, result.(error))
	assert.False(t, continuePipeline)
}

func TestTransformFromCloudEvent(t *testing.T) {
	value := "4321"
	valueType := common.ValueTypeInt64
	cloudeventBytes := []byte(`{"data":"` + value + `","datacontenttype":"application/json","floatencoding":"", "eventid":"` + "event-" + devID1 + `","id":"123-abc","source":"id1","specversion":"1.0","time":"1970-01-01T00:00:00Z","type":"test-reading","valuetype":"` + valueType + `"}`)
	cloudevent := cloudevents.Event{}
	cloudevent.UnmarshalJSON(cloudeventBytes)

	expectedEvent := models.Event{
		Id:         "event-" + devID1,
		DeviceName: devID1,
		Readings: []models.Reading{
			models.SimpleReading{
				BaseReading: models.BaseReading{
					Id:           "123-abc",
					ResourceName: "test-reading",
					ValueType:    valueType},
				Value: value,
			},
		},
	}

	conv := NewConversion()
	continuePipeline, result := conv.TransformFromCloudEvent(context, []cloudevents.Event{cloudevent})

	edgexEvent, ok := result.(models.Event)
	assert.NotNil(t, result)
	assert.True(t, ok)
	assert.True(t, continuePipeline)
	assert.Equal(t, expectedEvent.Id, edgexEvent.Id)
	assert.Equal(t, expectedEvent.Readings[0].GetBaseReading().ValueType, edgexEvent.Readings[0].GetBaseReading().ValueType)
	assert.Equal(t, expectedEvent.Readings[0].GetBaseReading().ResourceName, edgexEvent.Readings[0].GetBaseReading().ResourceName)
	assert.Equal(t, expectedEvent.Readings[0].(models.SimpleReading).Value, edgexEvent.Readings[0].(models.SimpleReading).Value)
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
