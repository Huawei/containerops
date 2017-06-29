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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Huawei/containerops/common"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pilotage",
	Short: "pilotage is a DevOps orchestration engine.",
	Long: `Pilotage is the core module of ContainerOps. And it is the engine DevOps Orchestration.
It has two modes of running: cli and daemon. The cli mode reads a file of orchestration
flow and executes it. It uses the kubectl of local connecting the Kubernetes cluster and
collecting logs.The daemon mode is HTTPS server and exposes APIs interacting with Web UI.
And UI has editor and monitor for orchestration flow.
`,
}

// init()
func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file path")
	viper.BindPFlag("config", RootCmd.Flags().Lookup("config"))
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if err := common.SetConfig(cfgFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
