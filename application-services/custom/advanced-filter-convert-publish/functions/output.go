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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

func PrintFloatValuesToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("Convert to Readable Float Values")

	if data == nil {
		return false, errors.New("PrintFloatValuesToConsole: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("PrintFloatValuesToConsole: didn't receive expect Event type")
	}

	for _, eventReading := range event.Readings {
		fmt.Printf("%s readable value from %s is %s\n", eventReading.ResourceName, event.DeviceName, eventReading.Value)
	}

	return true, event

}

func Publish(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("Publish")

	if data == nil {
		return false, errors.New("Publish: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("Publish: didn't receive expect Event type")
	}

	payload, _ := json.Marshal(event)

	// By calling SetResponseData, the filtered and converted events will be posted back to the message bus on the new topic defined in the configuration.
	ctx.SetResponseData(payload)
	ctx.SetResponseContentType("application/json")
	return false, nil
}
