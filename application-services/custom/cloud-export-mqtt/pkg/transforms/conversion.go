package transforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
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
		return false, errors.New("No Event Received")
	}

	ctx.LoggingClient().Debug("Transforming to AWS format")

	if event, ok := data.(models.Event); ok {
		readings := map[string]interface{}{}

		for _, reading := range event.Readings {
			br, ok := reading.(models.SimpleReading)

			if !ok {
				continue
			}

			readings[br.ResourceName] = br.Value
		}

		msg, err := json.Marshal(readings)
		if err != nil {
			return false, errors.New(fmt.Sprintf("Failed to transform AWS data: %s", err))
		}

		return true, string(msg)
	}

	return false, errors.New("Unexpected type received")
}
