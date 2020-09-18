// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-uart-go"
	"github.com/edgexfoundry/device-uart-go/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-uart"
)

func main() {
	var baud int = 115200
	var devicePath string = "/dev/ttyUSB0"

	d := driver.Driver{DevicePath: devicePath, Baud: baud}
	startup.Bootstrap(serviceName, device.Version, &d)
}
