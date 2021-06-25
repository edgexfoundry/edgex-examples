// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"

	"github.com/edgexfoundry/device-random"
	"github.com/edgexfoundry/device-random/driver"
)

const (
	serviceName string = "device-random"
)

func main() {
	d := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_random.Version, d)
}
