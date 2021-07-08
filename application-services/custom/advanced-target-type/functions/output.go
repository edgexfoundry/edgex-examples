//
// Copyright (c) 2019 Intel Corporation
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

package functions

import (
	"fmt"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
)

func PrintXmlToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	ctx.LoggingClient().Debug("PrintXmlToConsole")

	if data == nil {
		// We didn't receive a result
		return false, nil
	}

	fmt.Println(data.(string))

	return true, data
}
