/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

	"github.com/Huawei/containerops/singular/vm"
	cobra "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "singular",
	Short: "The kubernetes deployment and operations tools",
	Long:  ``,
	Run:   rootion,
}

func rootion(cmd *cobra.Command, args []string) {
	fmt.Println("Singular version 0.1, build 5604cbe")
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var mSize int
var region string
var count int

func init() {
	//RootCmd.AddCommand(versionCmd)
	RootCmd.Flags().StringVarP(&address, "version", "v", "", "Show the Singular version information")

	//ConfigCmd
	RootCmd.AddCommand(ConfigCmd)

	ConfigCmd.Flags().IntVarP(&mSize, "mSize", "m", 512, "Node memory Size")
	ConfigCmd.Flags().StringVarP(&region, "region", "r", "sfo", "Cluster's localization of the region")
	ConfigCmd.Flags().StringVarP(&region, "system", "s", "ubuntu-17-04-x64", "Virtual machine system version")

	//createCmd
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().IntVarP(&count, "count", "c", 3, "Number of nodes in cluster")
	createCmd.AddCommand(pullCmd)

	RootCmd.AddCommand(cerkeyCmd)

	RootCmd.AddCommand(deployCmd)

	// versionCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "http or https listen address.")
	// versionCmd.Flags().Int64VarP(&port, "port", "p", 80, "the port of http.")
	// versionCmd.Flags().StringVarP(&address, "a", "v", "", "Source directory to read from")

}

var address string
var port int64

// var versionCmd = &cobra.Command{
// 	Use:   "version",
// 	Short: "",
// 	Long:  ``,
// 	Run:   version,
// }

// func version(cmd *cobra.Command, args []string) {
// 	fmt.Println("Singular version 0.1, build 5604cbe  " + address)
// }

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your nodes of kubernetes cluster",
	Long:  `Configure your nodes of kubernetes cluster`,
	Run:   config,
}

func config(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)
	// fmt.Println("Hugo Static Site Generator v0.9 -- HEAD  " + strings.Join(port, ""))

}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,
	Run:   Install,
}

func Install(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run:   create,
}

func create(cmd *cobra.Command, args []string) {
	//vm.CreateNewVM("lidian-unbantu-wk-master")

	vm.CreateVMs()
}

var pullCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run:   pull,
}

func pull(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "",
	Long:  ``,
	Run:   deploy,
}

func deploy(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)

	//get iplist while
	// download.Download_main()
	// deploy.DeployNodes()

}

var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "",
	Long:  ``,
	Run:   apikey,
}

func apikey(cmd *cobra.Command, args []string) {

	fmt.Println("[singular] API Server key is ready.")

}

var cerpathCmd = &cobra.Command{
	Use:   "cerpath",
	Short: "",
	Long:  ``,
	Run:   cerpath,
}

func cerpath(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)
}

var cerkeyCmd = &cobra.Command{
	Use:   "cerkey",
	Short: "",
	Long:  ``,
	Run:   cerkey,
}

func cerkey(cmd *cobra.Command, args []string) {
	fmt.Println("[singular] Generated Certificate Authority key and certificate.")
	fmt.Println("[singular] Created keys and certificates in \"/usr/singular/\"")
	fmt.Println(args)
}
