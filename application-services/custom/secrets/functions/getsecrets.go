//
// Copyright (c) 2020 Intel Corporation
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

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
)

type secretsInfo struct {
	path string
	keys []string
}

var expectedSecrets = []secretsInfo{
	secretsInfo{
		// Get secrets for valid sub-path, empty keys list should return all secrets
		path: "/mqtt",
		keys: []string{},
	},
	secretsInfo{
		// Get secrets for valid sub-path with single key
		path: "/mqtt",
		keys: []string{"password"},
	},
}

// GetSecretsToConsole ...
func GetSecretsToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	for _, secretInfo := range expectedSecrets {
		// this is just an example. NEVER log your secrets to console
		edgexcontext.LoggingClient.Info(fmt.Sprintf("--- Get secrets at location %v, keys: %v  ---", secretInfo.path, secretInfo.keys))
		secrets, err := edgexcontext.GetSecrets(secretInfo.path, secretInfo.keys...)
		if err != nil {
			edgexcontext.LoggingClient.Error(err.Error())
			return false, nil
		}
		for k, v := range secrets {
			edgexcontext.LoggingClient.Info(fmt.Sprintf("key:%v, value:%v", k, v))
		}
	}

	edgexcontext.Complete([]byte(params[0].(string)))
	return false, nil
}
