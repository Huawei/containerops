/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	"fmt"
	"os"

	. "github.com/logrusorgru/aurora"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/pilotage/module"
)

// SetRunDaemonMiddlewares is setting function when the pilotage run a HTTP daemon
// with flow file.
//   1. The mode doesn't export log to file or database, it prints logs in Terminal.
//   2. The mode doesn't have database support, the flow definition from file.
//   3. The mode
func SetRunDaemonMiddlewares(m *macaron.Macaron, cfgFile, flowFile string) {
	flow := new(module.Flow)

	if err := flow.ParseFlowFromFile(flowFile, module.DaemonRun, true, true); err != nil {
		fmt.Println(Red("Parse flow file error: "), err.Error())
		os.Exit(1)
	}

	go func() {
		flow.LocalRun(true, true)
	}()

	// Init the flow and set into context.
	m.Use(func(ctx *macaron.Context) {
		ctx.Data["flow"] = flow
		ctx.Data["mode"] = module.DaemonRun
	})

	//
	m.Use(macaron.Logger())

	// Set recovery handler to returns a middleware that recovers from any panics
	m.Use(macaron.Recovery())
}

func SetStartDaemonMiddlewares(m *macaron.Macaron, cfgFile string) {
	// Nil flow and set into context.
	m.Use(func(ctx *macaron.Context) {
		ctx.Data["flow"] = nil
		ctx.Data["mode"] = module.DaemonStart
	})
}
