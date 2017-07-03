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

package cmd

import (
	"fmt"
	"os"

	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/pilotage/module"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "pilotage cli mode",
	Long: `Pilotage cli mode runs orchestration flow in local. It uses the kubectl command interacting
with Kubernetes master which create pod and get logs. The kubectl install and config document:

https://kubernetes.io/docs/user-guide/kubectl-overview

1. The cli mode doesn't have trigger.
2. The cli mode doesn't have database, never save result and log.'`,
}

var runCliCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a orchestration flow.",
	Long:  ``,
	Run:   runCliFlow,
}

// init()
func init() {
	// Add cli sub command.
	RootCmd.AddCommand(cliCmd)

	//Add run sub command to cli.
	cliCmd.AddCommand(runCliCmd)

}

// Run orchestration flow from a flow definition file.
func runCliFlow(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Println(Red("The orchestration flow file is required."))
		os.Exit(1)
	}

	flowFile := args[0]

	if utils.IsFileExist(flowFile) == false {
		cmd.Println(Red("The orchestration flow file is not exist."))
		os.Exit(1)
	}

	flow := new(module.Flow)

	if err := flow.ExecuteFlowFromFile(flowFile, module.CliRun, verbose, timestamp); err != nil {
		cmd.Println(fmt.Sprintf("[red]Execute orchestration flow error: %s", err.Error()))
		os.Exit(1)
	}

}
