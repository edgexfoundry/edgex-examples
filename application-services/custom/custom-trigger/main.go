//
// Copyright (c) 2020 Technotects
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
			s = strings.TrimRight(s, "\n")

			if err != nil {
				t.tc.Logger.Error(err.Error())
				continue
			}

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
					env := types.MessageEnvelope{
						Payload: m,
					}

					ctx := t.tc.ContextBuilder(env)

					err := t.tc.MessageProcessor(ctx, env)

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

	service.SetFunctionsPipeline(
		printToConsole,
	)

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	err := service.MakeItRun()
	if err != nil {
		service.LoggingClient().Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	os.Exit(0)
}

func printToConsole(appContext interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	input, err := util.CoerceType(data)

	if err != nil {
		appContext.LoggingClient().Error(err.Error())
		return false, err
	}

	wait := time.Millisecond * time.Duration(len(input))

	time.Sleep(wait)

	appContext.LoggingClient().Info("PrintToConsole")

	os.Stdout.WriteString(fmt.Sprintf("'%s' received %s ago\n>", string(input), wait.String()))

	return false, nil
}
