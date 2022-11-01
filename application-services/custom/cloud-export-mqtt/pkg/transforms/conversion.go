package transforms

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

// Conversion Struct
type Conversion struct {
}

// NewConversion returns a conversion struct
func NewConversion() Conversion {
	return Conversion{}
}

// TransformToCloudFormat converts the event into AWS readable format
func (f Conversion) TransformToCloudFormat(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, stringType interface{}) {
	if data == nil {
		return false, errors.New("no Event received")
	}

	ctx.LoggingClient().Debug("Transforming to AWS format")

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("unexpected type received")
	}

	readings := map[string]interface{}{}

	for _, reading := range event.Readings {
		switch reading.ValueType {
		case common.ValueTypeBinary:
			readings[reading.ResourceName] = reading.BinaryValue
		case common.ValueTypeObject:
			readings[reading.ResourceName] = reading.ObjectValue
		default:
			readings[reading.ResourceName] = reading.Value
		}
	}

	msg, err := json.Marshal(readings)
	if err != nil {
		return false, fmt.Errorf("failed to transform Event to cloud format: %s", err)
	}

	return true, string(msg)
}
