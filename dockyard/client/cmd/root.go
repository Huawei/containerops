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

var cfgFile, domain string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "warship",
	Short: "Warship is a dockyard client.",
	Long: `Dockyard is a container and artifact repository storing and distributing 
container image, software artifact and virtual images of KVM or XEN. 
Warship is the Dockyard client which push/pull binary file with Dockyard server.
And we also working on push/pull Docker image, rkt ACI and OCI Image.`,
}

// init()
func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file path")
	RootCmd.PersistentFlags().StringVar(&domain, "domain", "", "Dockyard service domain")

	viper.BindPFlag("domain", RootCmd.Flags().Lookup("domain"))
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
