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
	"github.com/spf13/cobra"
)

// auth sub command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Login and logout a Dockyard server.",
	Long:  ``,
}

// login sub command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login the Dockyard server.",
	Long:  ``,
	Run:   loginServer,
}

// logout sub command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Dockyard server.",
	Long:  ``,
	Run:   logoutServer,
}

// init()
func init() {
	RootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
}

// loginServer is
func loginServer(cmd *cobra.Command, args []string) {

}

// logoutServer is
func logoutServer(cmd *cobra.Command, args []string) {

}
