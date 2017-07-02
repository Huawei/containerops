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

import "github.com/spf13/cobra"

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "pilotage daemon mode",
	Long: `Pilotage daemon will runs orchestration flow in a daemon mode which has Web GUI and REST API.
It will provide endpoint for flow. Pilotage provides a simple mode running with flow YAML file.

pilotage daemon run cncf-demo.yaml

Pilotage provides full daemon run mode with HTTPS support.

pilotage daemon start --listen https

Pilotage daemon run with database supported, and create, modify or run a flow.'`,
}

var runDaemonCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a flow with a YAML file",
	Long:  ``,
	Run:   runDaemonFlow,
}

var startDaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run a flow with HTTPS/Unix mode",
	Long:  ``,
	Run:   startDaemonFlow,
}

// init()
func init() {
	// Add cli sub command.
	RootCmd.AddCommand(daemonCmd)

	daemonCmd.AddCommand(runDaemonCmd)
	daemonCmd.AddCommand(startDaemonCmd)
}

// runDaemonFlow is
func runDaemonFlow(cmd *cobra.Command, args []string) {

}

// startDaemonFlow is
func startDaemonFlow(cmd *cobra.Command, args []string) {

}
