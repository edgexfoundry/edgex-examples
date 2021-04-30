package transforms

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"
)

// Conversion Struct
type Conversion struct {
}

// NewConversion returns a conversion struct
func NewConversion() Conversion {
	return Conversion{}
}

// TransformToAWS converts the event into AWS readable format
func (f Conversion) TransformToAWS(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, stringType interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("Transforming to AWS format")

	if data == nil {
		return false, errors.New("TransformToAWS: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("TransformToAWS: didn't receive expect Event type")
	}

	readings := map[string]interface{}{}

	for _, reading := range event.Readings {
		readings[reading.ResourceName] = reading.Value
	}

	msg, err := json.Marshal(readings)
	if err != nil {
		return false, fmt.Errorf("Failed to transform AWS data: %s", err.Error())
	}

	return true, string(msg)
}
