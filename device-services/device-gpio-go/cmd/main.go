// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides device service  of a GPIO device service.
package main

import (
	"github.com/edgexfoundry/device-gpio-go"
	"github.com/edgexfoundry/device-gpio-go/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-gpio-go"
)

func main() {

	d := driver.Driver{}
	startup.Bootstrap(serviceName, device.Version, &d)
}
