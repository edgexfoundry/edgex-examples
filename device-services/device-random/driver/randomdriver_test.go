//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"testing"

	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

var d *RandomDriver

func init() {
	d = new(RandomDriver)
	d.lc = logger.NewMockClient()
}

func TestHandleReadCommands(t *testing.T) {
	deviceName := "testDevice"
	protocols := map[string]models.ProtocolProperties{
		"other": {
			"Address": "simple01",
			"Port":    "300",
		},
	}

	requests := []dsModels.CommandRequest{
		{
			DeviceResourceName: "RandomValue_Int8",
			Type:               common.ValueTypeInt8,
		},
		{
			DeviceResourceName: "RandomValue_Int16",
			Type:               common.ValueTypeInt16,
		},
		{
			DeviceResourceName: "RandomValue_Int32",
			Type:               common.ValueTypeInt32,
		},
	}

	res, err := d.HandleReadCommands(deviceName, protocols, requests)

	if err != nil {
		t.Fatalf("Failed to read command, %v", err)
	}
	if len(res) != len(requests) {
		t.Fatalf("Number of results fetched '%v' should match '%v'", len(res), len(requests))
	}
	if res[0].DeviceResourceName != "RandomValue_Int8" || res[1].DeviceResourceName != "RandomValue_Int16" || res[2].DeviceResourceName != "RandomValue_Int32" {
		t.Fatalf("Unexpected test result. Wrong resource object.")
	}
	if res[0].Type != common.ValueTypeInt8 || res[1].Type != common.ValueTypeInt16 || res[2].Type != common.ValueTypeInt32 {
		t.Fatalf("Unexpected test result. Wrong value type.")
	}
}

func TestHandleWriteCommands(t *testing.T) {
	deviceName := "testDevice"
	protocols := map[string]models.ProtocolProperties{
		"other": {
			"Address": "simple01",
			"Port":    "300",
		},
	}
	var requests []dsModels.CommandRequest

	cv, err := dsModels.NewCommandValue("Max_Int8", common.ValueTypeInt8, int8(127))
	if err != nil {
		t.Fatalf("Failed to create command value, %v", err)
	}
	params := []*dsModels.CommandValue{cv}

	err = d.HandleWriteCommands(deviceName, protocols, requests, params)

	if err != nil {
		t.Fatalf("Failed to write command, %v", err)
	}
}
