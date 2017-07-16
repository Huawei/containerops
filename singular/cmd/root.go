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

// RootCmd is
var RootCmd = &cobra.Command{
	Use:   "singular",
	Short: "Singular is deploy and operation tools for ContainerOps platform.",
	Long: `Singular design for deploy and operation ContainerOps platform,
mostly focus on Kubernetes, Prometheus, others in the Cloud Native technology stack.
We are trying to deploy stack cross cloud, OpenStack, bare metals. `,
}

// init()
func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file path")

	viper.BindPFlag("config", RootCmd.Flags().Lookup("config"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if err := common.SetConfig(cfgFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
