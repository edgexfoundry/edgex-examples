// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Systems LTD
//
// SPDX-License-Identifier: Apache-2.0

package transforms

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

type Conversion struct {
}

func NewConversion() Conversion {
	return Conversion{}
}

func (f Conversion) TransformToInflux(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("Transforming to InfluxDB Line Protocol format")

	if data == nil {
		return false, errors.New("TransformToInflux: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("TransformToInflux: didn't receive expect Event type")
	}

	var buffer strings.Builder

	// write device name as measurement
	buffer.WriteString(event.DeviceName)
	// write tags if any, comma separated
	// see Influx docs for syntax and example
	// https://docs.influxdata.com/influxdb/v2.0/reference/syntax/line-protocol/
	for key, val := range event.Tags {
		// write comma
		buffer.WriteString(",")
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(fmt.Sprintf("%v", val))
	}
	// write space
	buffer.WriteString(" ")
	// write fields (readings) comma separated
	// see Influx docs for syntax and example
	// https://docs.influxdata.com/influxdb/v2.0/reference/syntax/line-protocol/
	for j, reading := range event.Readings {
		if j > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(reading.ResourceName)
		buffer.WriteString("=")
		buffer.WriteString(reading.Value)
	}
	// write space
	buffer.WriteString(" ")
	// write timestamp in nanosecond form
	buffer.WriteString(strconv.Itoa(int(event.Origin)))
	msg := buffer.String()
	lc.Debugf("InfluxDB Payload: %s", msg)
	return true, msg
}
