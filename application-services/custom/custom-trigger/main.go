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
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/util"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap"
	"github.com/edgexfoundry/go-mod-messaging/pkg/types"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	serviceKey = "customTrigger"
)

type stdinTrigger struct {
	tc appsdk.TriggerConfig
}

func (t *stdinTrigger) Initialize(wg *sync.WaitGroup, ctx context.Context, background <-chan types.MessageEnvelope) (bootstrap.Deferred, error) {
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
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	// First thing to do is to create an instance of the EdgeX SDK and initialize it.
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey, TargetType: &[]byte{}}

	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		os.Exit(-1)
	}

	edgexSdk.RegisterCustomTriggerFactory("custom-stdin", func(config appsdk.TriggerConfig) (appsdk.Trigger, error) {
		return &stdinTrigger{
			tc: config,
		}, nil
	})

	edgexSdk.SetFunctionsPipeline(
		printToConsole,
	)

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	err := edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	os.Exit(0)
}

func printToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	input, err := util.CoerceType(params[0])

	if err != nil {
		edgexcontext.LoggingClient.Error(err.Error())
		return false, err
	}

	wait := time.Millisecond * time.Duration(len(input))

	time.Sleep(wait)

	edgexcontext.LoggingClient.Info("PrintToConsole")

	os.Stdout.WriteString(fmt.Sprintf("'%s' received %s ago\n>", string(input), wait.String()))

	return false, nil
}
