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
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/util"
)

func PrintToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	ctx.LoggingClient().Debug("PrintToConsole")

	if data == nil {
		return false, errors.New("PrintToConsole: No data received")
	}

	// Data is expected to be a JSON response to previous command executed.
	bytes, err := util.CoerceType(data)
	if err != nil {
		return false, fmt.Errorf("PrintToConsole: CoerceType failed: %s", err.Error())
	}

	strData := string(bytes)
	ctx.LoggingClient().Info(strData)

	ctx.SetResponseContentType("application/json")

	return true, data
}
