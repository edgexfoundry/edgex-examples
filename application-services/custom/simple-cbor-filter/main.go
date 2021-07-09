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

package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

const (
	serviceKey = "app-sample-cbor-filter"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an new instance of an EdgeX Application Service.
	service, ok := pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(-1)
	}

	// Leverage the built in logging service in EdgeX
	lc := service.LoggingClient()

	// 2) shows how to access the application's specific configuration settings.
	resourceNames, err := service.GetAppSettingStrings("ResourceNames")
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}
	lc.Info(fmt.Sprintf("Filtering for ResourceNames %v", resourceNames))

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		transforms.NewFilterFor(resourceNames).FilterByResourceName,
		processImages,
	); err != nil {
		lc.Error("SetFunctionsPipeline returned error: ", err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func processImages(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Leverage the built in logging service in EdgeX
	lc := ctx.LoggingClient()

	if data == nil {
		return false, errors.New("processImages: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("processImages: didn't receive expect Event type")
	}

	for _, reading := range event.Readings {
		// For this to work the image/jpeg & image/png packages must be imported to register their decoder
		imageData, imageType, err := image.Decode(bytes.NewReader(reading.BinaryValue))

		if err != nil {
			return false, errors.New("processImages: unable to decode image: " + err.Error())
		}

		// Since this is a example, we will just print put some stats from the images received
		lc.Infof("Received Image from Device: %s, ResourceName: %s, Image Type: %s, Image Size: %s, Color in middle: %v\n",
			reading.DeviceName, reading.ResourceName, imageType, imageData.Bounds().Size().String(),
			imageData.At(imageData.Bounds().Size().X/2, imageData.Bounds().Size().Y/2))
	}

	return false, nil
}
