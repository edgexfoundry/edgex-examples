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

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
)

func PrintXMLToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Leverage the built in logging service in EdgeX
	lc := ctx.LoggingClient()

	if data == nil {
		return false, errors.New("PrintXMLToConsole: No data received")
	}

	xml, ok := data.(string)
	if !ok {
		return false, errors.New("PrintXMLToConsole: Data received is not the expected 'string' type")

	}

	lc.Debug(xml)
	ctx.SetResponseData([]byte(xml))
	return true, xml
}
