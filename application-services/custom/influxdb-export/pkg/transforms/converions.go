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

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type Conversion struct {
}

func NewConversion() Conversion {
	return Conversion{}
}

func (f Conversion) TransformToInflux(edgexcontext *appcontext.Context, params ...interface{}) (continuePipeline bool, stringType interface{}) {
	if len(params) < 1 {
		return false, errors.New("No Event Received")
	}

	edgexcontext.LoggingClient.Debug("Transforming to InfluxDB Line Protocol format")

	if event, ok := params[0].(models.Event); ok {

		var buffer strings.Builder

		// write device name as measurement
		buffer.WriteString(event.Device)
		// write tags if any, comma separated [,<tag_key>=<tag_value>[,<tag_key>=<tag_value>]]
		for i, tag := range event.Tags {
			// write comma
			buffer.WriteString(",")
			buffer.WriteString("tag")
			buffer.WriteString(i)
			buffer.WriteString("=")
			buffer.WriteString(tag)
		}
		// write space
		buffer.WriteString(" ")
		// write fields (readings) comma separated <field_key>=<field_value>[,<field_key>=<field_value>]
		for j, reading := range event.Readings {
			if j > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(reading.Name)
			buffer.WriteString("=")
			buffer.WriteString(reading.Value)
		}
		// write space
		buffer.WriteString(" ")
		// write timestamp in nanosecond form
		buffer.WriteString(strconv.Itoa(int(event.Origin)))
		msg := buffer.String()
		edgexcontext.LoggingClient.Debug(fmt.Sprintf("InfluxDB Payload: %s", msg))
		return true, string(msg)
	}

	return false, errors.New("Unexpected type received")
}
