//
// Copyright (c) 2021 Intel Corporation
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

type secretsInfo struct {
	path string
	keys []string
}

var expectedSecrets = []secretsInfo{
	{
		// Get secrets for valid sub-path, empty keys list should return all secrets
		path: "/mqtt",
		keys: []string{},
	},
	{
		// Get secrets for valid sub-path with single key
		path: "/mqtt",
		keys: []string{"password"},
	},
}

// GetSecretsToConsole ...
func GetSecretsToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	if data == nil {
		return false, errors.New("GetSecretsToConsole: No data received")
	}

	lc := ctx.LoggingClient()

	for _, secretInfo := range expectedSecrets {
		// this is just an example. NEVER log your secrets to console
		lc.Infof("--- Get secrets at location %v, keys: %v  ---", secretInfo.path, secretInfo.keys)
		secrets, err := ctx.GetSecret(secretInfo.path, secretInfo.keys...)
		if err != nil {
			lc.Error(err.Error())
			return false, nil
		}
		for k, v := range secrets {
			lc.Infof("key:%v, value:%v", k, v)
		}
	}

	ctx.SetResponseData([]byte(data.(string)))
	return false, nil
}
