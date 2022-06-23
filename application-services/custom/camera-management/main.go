//
// Copyright (c) 2022 Intel Corporation
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

package main

import (
	"fmt"
	appsdk "github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/edgex-examples/application-services/custom/camera-management/appcamera"
	"os"
)

const (
	serviceKey = "app-camera-management"
)

func main() {
	service, ok := appsdk.NewAppService(serviceKey)
	if !ok {
		fmt.Printf("error: unable to create new app service %s!\n", serviceKey)
		os.Exit(-1)
	}

	app := appcamera.NewCameraManagementApp(service)
	if err := app.Run(); err != nil {
		service.LoggingClient().Error(err.Error())
		os.Exit(-1)
	}

	os.Exit(0)
}
