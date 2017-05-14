/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/dockyard/utils"
	"github.com/Huawei/dockyard/web"
	"github.com/containerops/configure"
)

var address string
var port int64

// webCmd is subcommand which start/stop/monitor Dockyard's REST API daemon.
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Web subcommand start/stop/monitor Dockyard's REST API daemon.",
	Long:  ``,
}

// start Dockyard deamon subcommand
var startDaemonCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Dockyard's REST API daemon.",
	Long:  ``,
	Run:   startDeamon,
}

// stop Dockyard deamon subcommand
var stopDaemonCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop Dockyard's REST API daemon.",
	Long:  ``,
	Run:   stopDaemon,
}

// monitor Dockyard deamon subcommand
var monitorDeamonCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor Dockyard's REST API daemon.",
	Long:  ``,
	Run:   monitorDaemon,
}

// init()
func init() {
	RootCmd.AddCommand(daemonCmd)

	// Add start subcommand
	daemonCmd.AddCommand(startDaemonCmd)
	startDaemonCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "http or https listen address.")
	startDaemonCmd.Flags().Int64VarP(&port, "port", "p", 80, "the port of http.")

	// Add stop subcommand
	daemonCmd.AddCommand(stopDaemonCmd)
	// Add daemon subcommand
	daemonCmd.AddCommand(monitorDeamonCmd)
}

// startDeamon() start Dockyard's REST API daemon.
func startDeamon(cmd *cobra.Command, args []string) {
	m := macaron.New()

	// Set Macaron Web Middleware And Routers
	web.SetDockyardMacaron(m)

	listenMode := configure.GetString("listenmode")
	switch listenMode {
	case "http":
		listenaddr := fmt.Sprintf("%s:%d", address, port)
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			fmt.Printf("Start Dockyard http service error: %v\n", err.Error())
		}
		break
	case "https":
		listenaddr := fmt.Sprintf("%s:443", address)
		server := &http.Server{Addr: listenaddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		if err := server.ListenAndServeTLS(configure.GetString("httpscertfile"), configure.GetString("httpskeyfile")); err != nil {
			fmt.Printf("Start Dockyard https service error: %v\n", err.Error())
		}
		break
	case "unix":
		listenaddr := fmt.Sprintf("%s", address)
		if utils.IsFileExist(listenaddr) {
			os.Remove(listenaddr)
		}

		if listener, err := net.Listen("unix", listenaddr); err != nil {
			fmt.Printf("Start Dockyard unix socket error: %v\n", err.Error())

		} else {
			server := &http.Server{Handler: m}
			if err := server.Serve(listener); err != nil {
				fmt.Printf("Start Dockyard unix socket error: %v\n", err.Error())
			}
		}
		break
	default:
		break
	}
}

// stopDaemon() stop Dockyard's REST API daemon.
func stopDaemon(cmd *cobra.Command, args []string) {

}

// monitordAemon() monitor Dockyard's REST API deamon.
func monitorDaemon(cmd *cobra.Command, args []string) {

}
