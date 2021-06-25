//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"math"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
)

var device = newRandomDevice()

func TestValue_Int8(t *testing.T) {
	valueType := common.ValueTypeInt8

	val, err := device.value(valueType)

	if err != nil {
		t.Fatalf("Failed to generate random %v value", valueType)
	}
	if val <= math.MinInt8 || val >= math.MaxInt8 {
		t.Fatalf("Unexpected test result. %v is not in %v value range", val, valueType)
	}
}

func TestValue_Int16(t *testing.T) {
	valueType := common.ValueTypeInt16

	val, err := device.value(valueType)

	if err != nil {
		t.Fatalf("Failed to generate random %v value", valueType)
	}
	if val <= math.MinInt16 || val >= math.MaxInt16 {
		t.Fatalf("Unexpected test result. %v is not in %v value range", val, valueType)
	}
}

func TestValue_Int32(t *testing.T) {
	valueType := common.ValueTypeInt32

	val, err := device.value(valueType)

	if err != nil {
		t.Fatalf("Failed to generate random %v value", valueType)
	}
	if val <= math.MinInt32 || val >= math.MaxInt32 {
		t.Fatalf("Unexpected test result. %v is not in %v value range", val, valueType)
	}
}

func TestValue_Int64(t *testing.T) {
	valueType := common.ValueTypeInt64

	_, err := device.value(valueType)

	if err == nil {
		t.Fatalf("RandomDevice only support Int8, Int16, Int32. Use %v should return error.", valueType)
	}
}
