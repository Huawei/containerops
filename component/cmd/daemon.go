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

	logs "github.com/Huawei/containerops/component/log"
	"github.com/Huawei/containerops/component/models"
	"github.com/Huawei/containerops/component/utils"
	"github.com/Huawei/containerops/component/web"
	"github.com/containerops/configure"
)

// address is address that server will listen
var address string

// port is port that server will listen
var port int64

// log is
var log *logs.Logger

// webCmd is subcommand which start/stop/monitor Pilotage's REST API daemon.
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Web subcommand start/stop/monitor Pilotage's REST API daemon.",
	Long:  ``,
}

// start Pilotage deamon subcommand
var startDaemonCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Pilotage's REST API daemon.",
	Long:  ``,
	Run:   startDeamon,
}

// stop Pilotage deamon subcommand
var stopDaemonCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop Pilotage's REST API daemon.",
	Long:  ``,
	Run:   stopDaemon,
}

// monitor Pilotage deamon subcommand
var monitorDeamonCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor Pilotage's REST API daemon.",
	Long:  ``,
	Run:   monitorDaemon,
}

// init is
func init() {
	log = logs.New()
	RootCmd.AddCommand(daemonCmd)

	defaultPort := int64(80)
	if configure.GetString("http.listenmode") == "https" {
		defaultPort = int64(443)
	}

	// Add start subcommand
	daemonCmd.AddCommand(startDaemonCmd)
	startDaemonCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "http or https listen address.")
	startDaemonCmd.Flags().Int64VarP(&port, "port", "p", defaultPort, "the port of http.")

	// Add stop subcommand
	daemonCmd.AddCommand(stopDaemonCmd)
	// Add daemon subcommand
	daemonCmd.AddCommand(monitorDeamonCmd)
}

// startDeamon is start Component's REST API daemon.
func startDeamon(cmd *cobra.Command, args []string) {
	models.OpenDatabase()

	m := macaron.New()

	// Set Macaron Web Middleware And Routers
	web.SetPilotagedMacaron(m)

	listenMode := configure.GetString("listenmode")
	switch listenMode {
	case "http":
		listenaddr := fmt.Sprintf("%s:%d", address, port)
		log.Debugln("component is listening:", listenaddr)
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			log.Errorf("Start Pilotage http service error: %v\n", err.Error())
			return
		}
	case "https":
		listenaddr := fmt.Sprintf("%s:%d", address, port)
		server := &http.Server{Addr: listenaddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		log.Debugln("component is listening:", listenaddr)
		if err := server.ListenAndServeTLS(configure.GetString("httpscertfile"), configure.GetString("httpskeyfile")); err != nil {
			log.Errorf("Start Pilotage https service error: %v\n", err.Error())
			return
		}
	case "unix":
		listenaddr := fmt.Sprintf("%s", address)
		if utils.IsFileExist(listenaddr) {
			os.Remove(listenaddr)
		}

		if listener, err := net.Listen("unix", listenaddr); err != nil {
			log.Errorf("Start Pilotage unix socket error: %v\n", err.Error())
		} else {
			server := &http.Server{Handler: m}
			if err := server.Serve(listener); err != nil {
				log.Errorf("Start Pilotage unix socket error: %v\n", err.Error())
				return
			}
		}
	default:
		break
	}
}

// stopDaemon is stop Pilotage's REST API daemon.
func stopDaemon(cmd *cobra.Command, args []string) {

}

// monitordAemon() monitor Pilotage's REST API deamon.
func monitorDaemon(cmd *cobra.Command, args []string) {

}
