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

package functions

import (
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/util"
)

type secretsInfo struct {
	path string
	keys []string
}

var expectedSecrets = []secretsInfo{
	secretsInfo{
		// Get secrets for valid sub-path, empty keys list should return all secrets
		path: "mqtt",
		keys: []string{},
	},
	secretsInfo{
		// Get secrets for valid sub-path with single key
		path: "mqtt",
		keys: []string{"password"},
	},
}

// GetSecretsToConsole ...
func GetSecretsToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	if data == nil {
		// We didn't receive a result
		return false, nil
	}

	for _, secretInfo := range expectedSecrets {
		// this is just an example. NEVER log your secrets to console
		ctx.LoggingClient().Infof("--- Get secrets at location %v, keys: %v  ---", secretInfo.path, secretInfo.keys)

		secrets, err := ctx.GetSecret(secretInfo.path, secretInfo.keys...)
		if err != nil {
			return false, err
		}
		for k, v := range secrets {
			ctx.LoggingClient().Infof("key:%v, value:%v", k, v)
		}
	}

	response, err := util.CoerceType(data)

	if err != nil {
		return false, err
	}

	ctx.SetResponseData(response)
	return false, nil
}
