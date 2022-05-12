//
// Copyright (c) 2020 Technotects
// Copyright (c) 2021 One Track Consulting
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/util"
)

const (
	serviceKey = "app-custom-trigger"
)

type stdinTrigger struct {
	tc interfaces.TriggerConfig
}

func (t *stdinTrigger) Initialize(_ *sync.WaitGroup, ctx context.Context, _ <-chan interfaces.BackgroundMessage) (bootstrap.Deferred, error) {
	msgs := make(chan []byte)

	ctx, cancel := context.WithCancel(context.Background())

	receiveMessage := true

	go func() {
		fmt.Print("> ")
		rdr := bufio.NewReader(os.Stdin)
		for receiveMessage {
			s, err := rdr.ReadString('\n')

			if err != nil {
				t.tc.Logger.Error(err.Error())
				continue
			}

			s = strings.TrimRight(s, "\n")

			msgs <- []byte(s)
		}
	}()

	go func() {
		for receiveMessage {
			select {
			case <-ctx.Done():
				receiveMessage = false
			case m := <-msgs:
				go func() {
					spoofTopic := "even"

					if len(m)%2 == 1 {
						spoofTopic = "odd"
					}

					env := types.MessageEnvelope{
						CorrelationID: uuid.NewString(),
						Payload:       m,
						ReceivedTopic: spoofTopic,
					}

					t.tc.Logger.Tracef("sending message to runtime %+v", env)

					err := t.tc.MessageReceived(nil, env, t.responseHandler)

					if err != nil {
						t.tc.Logger.Error(err.Error())
					}
				}()
			}
		}
	}()

	return func() {
		cancel()
	}, nil
}

func (t *stdinTrigger) responseHandler(ctx interfaces.AppFunctionContext, pipeline *interfaces.FunctionPipeline) error {
	t.tc.Logger.Infof("Responding to pipeline %s with '%s'", pipeline.Id, string(ctx.ResponseData()))
	os.Stdout.WriteString(">")
	return nil
}

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// First thing to do is to create an instance of the EdgeX SDK Service, which also runs the bootstrap initialization.
	service, ok := pkg.NewAppServiceWithTargetType(serviceKey, &[]byte{})
	if !ok {
		os.Exit(-1)
	}

	service.RegisterCustomTriggerFactory("custom-stdin", func(config interfaces.TriggerConfig) (interfaces.Trigger, error) {
		return &stdinTrigger{
			tc: config,
		}, nil
	})

	var err error

	//use this to process using default pipeline only
	//err = service.SetDefaultFunctionsPipeline(printLowerToConsole)
	//if err != nil {
	//	service.LoggingClient().Errorf("SetDefaultFunctionsPipeline returned error: %s", err.Error())
	//	os.Exit(-1)
	//}

	//use this to process using varied pipelines by topic (odd/even string length)
	err = service.AddFunctionsPipelineForTopics("odd", []string{"odd"},
		printLowerToConsole,
	)

	if err != nil {
		service.LoggingClient().Errorf("AddFunctionsPipelineForTopic returned error: %s", err.Error())
		os.Exit(-1)
	}

	err = service.AddFunctionsPipelineForTopics("even", []string{"even"},
		printUpperToConsole,
	)

	if err != nil {
		service.LoggingClient().Errorf("AddFunctionsPipelineForTopic returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	err = service.MakeItRun()
	if err != nil {
		service.LoggingClient().Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	service.LoggingClient().Info("Exiting service")
	// Do any required cleanup here
	os.Exit(0)
}

func printLowerToConsole(appContext interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	input, err := util.CoerceType(data)

	if err != nil {
		appContext.LoggingClient().Error(err.Error())
		return false, err
	}

	wait := time.Millisecond * time.Duration(len(input))

	time.Sleep(wait)

	appContext.LoggingClient().Info("PrintToConsole")

	processed := strings.ToLower(string(input))

	appContext.LoggingClient().Infof("'%s' received %s ago\n>", processed, wait.String())

	appContext.SetResponseData([]byte(processed))
	return false, nil
}

func printUpperToConsole(appContext interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	input, err := util.CoerceType(data)

	if err != nil {
		appContext.LoggingClient().Error(err.Error())
		return false, err
	}

	wait := time.Millisecond * time.Duration(len(input))

	time.Sleep(wait)

	appContext.LoggingClient().Info("PrintToConsole")

	processed := strings.ToUpper(string(input))
	appContext.LoggingClient().Infof("'%s' received %s ago\n", processed, wait.String())

	appContext.SetResponseData([]byte(processed))

	return false, nil
}
