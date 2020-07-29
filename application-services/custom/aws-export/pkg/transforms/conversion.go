package transforms

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// Conversion Struct
type Conversion struct {
}

// NewConversion returns a conversion struct
func NewConversion() Conversion {
	return Conversion{}
}

// TransformToAWS converts the event into AWS readable format
func (f Conversion) TransformToAWS(edgexcontext *appcontext.Context, params ...interface{}) (continuePipeline bool, stringType interface{}) {
	if len(params) < 1 {
		return false, errors.New("No Event Received")
	}

	edgexcontext.LoggingClient.Debug("Transforming to AWS format")

	if event, ok := params[0].(models.Event); ok {
		readings := map[string]interface{}{}

		for _, reading := range event.Readings {
			readings[reading.Name] = reading.Value
		}

		msg, err := json.Marshal(readings)
		if err != nil {
			return false, errors.New(fmt.Sprintf("Failed to transform AWS data: %s", err))
		}

		return true, string(msg)
	}

	return false, errors.New("Unexpected type received")
}
