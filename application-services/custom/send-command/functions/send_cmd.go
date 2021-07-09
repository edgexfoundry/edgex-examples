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
	"context"
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
)

type ActionRequest struct {
	Action       string `json:"action"`
	DeviceName   string `json:"deviceName"`
	CommandName  string `json:"commandName"`
	ResourceName string `json:"resourceName"`
	Value        string `json:"value"`
}

type SendCommand struct {
}

func NewSendCommand() *SendCommand {
	return &SendCommand{}
}

func (s *SendCommand) SendCommand(funcCtx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := funcCtx.LoggingClient()

	lc.Debug("Sending Command")

	if data == nil {
		return false, errors.New("SendCommand: No data received")
	}

	if funcCtx.CommandClient == nil {
		return false, errors.New("SendCommand: Command client is not available")
	}

	actionRequest, ok := data.(ActionRequest)
	if !ok {
		return false, errors.New("SendCommand: Data received is not the expected 'ActionRequest' type")
	}
	action := actionRequest.Action
	device := actionRequest.DeviceName
	command := actionRequest.CommandName

	var response interface{}
	var err error

	switch action {
	case "set":
		lc.Infof("executing %s action", action)
		lc.Infof("Sending command '%s' for device '%s'", command, device)

		settings := make(map[string]string)
		settings[actionRequest.ResourceName] = actionRequest.Value
		response, err = funcCtx.CommandClient().IssueSetCommandByName(context.Background(), device, command, settings)
		if err != nil {
			return false, fmt.Errorf("failed to send '%s' set command to '%s' device: %s", command, device, err.Error())
		}

	case "get":
		lc.Infof("executing %s action", action)
		lc.Infof("Sending command '%s' for device '%s'", command, device)
		response, err = funcCtx.CommandClient().IssueGetCommandByName(context.Background(), device, command, "no", "yes")
		if err != nil {
			return false, fmt.Errorf("failed to send '%s' get command to '%s' device: %s", command, device, err.Error())
		}

	default:
		lc.Errorf("Invalid action requested: %s", action)
		return false, nil
	}

	return true, response
}
