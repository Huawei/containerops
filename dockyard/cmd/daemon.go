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
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/dockyard/model"
	"github.com/Huawei/containerops/dockyard/web"
)

var addressOption string
var portOption int

// webCmd is sub command which start/stop/monitor Dockyard's REST API daemon.
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Web sub command start/stop/monitor Dockyard's REST API daemon.",
	Long:  ``,
}

// start Dockyard daemon sub command
var startDaemonCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Dockyard's REST API daemon.",
	Long:  ``,
	Run:   startDeamon,
}

// stop Dockyard daemon sub command
var stopDaemonCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop Dockyard's REST API daemon.",
	Long:  ``,
	Run:   stopDaemon,
}

// monitor Dockyard daemon sub command
var monitorDaemonCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor Dockyard's REST API daemon.",
	Long:  ``,
	Run:   monitorDaemon,
}

// init()
func init() {
	RootCmd.AddCommand(daemonCmd)

	// Add start sub command
	daemonCmd.AddCommand(startDaemonCmd)
	startDaemonCmd.Flags().StringVarP(&addressOption, "address", "a", "", "http or https listen address.")
	startDaemonCmd.Flags().IntVarP(&portOption, "port", "p", 0, "the port of http.")

	// Add stop sub command
	daemonCmd.AddCommand(stopDaemonCmd)
	// Add daemon sub command
	daemonCmd.AddCommand(monitorDaemonCmd)
}

// startDeamon() start Dockyard's REST API daemon.
func startDeamon(cmd *cobra.Command, args []string) {
	model.OpenDatabase(&common.Database)
	m := macaron.New()

	// Set Macaron Web Middleware And Routers
	web.SetDockyardMacaron(m, cfgFile)

	var server *http.Server
	stopChan := make(chan os.Signal)

	signal.Notify(stopChan, os.Interrupt)

	address := common.Web.Address
	if addressOption != "" {
		address = addressOption
	}
	port := common.Web.Port
	if portOption != 0 {
		port = portOption
	}

	go func() {
		switch common.Web.Mode {
		case "https":
			listenAddr := fmt.Sprintf("%s:%d", address, port)
			server = &http.Server{Addr: listenAddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
			if err := server.ListenAndServeTLS(common.Web.Cert, common.Web.Key); err != nil {
				log.Errorf("Start Dockyard https service error: %s\n", err.Error())
			}

			break
		case "unix":
			listenAddr := fmt.Sprintf("%s", address)
			if utils.IsFileExist(listenAddr) {
				os.Remove(listenAddr)
			}

			if listener, err := net.Listen("unix", listenAddr); err != nil {
				log.Errorf("Start Dockyard unix socket error: %s\n", err.Error())
			} else {
				server = &http.Server{Handler: m}
				if err := server.Serve(listener); err != nil {
					log.Errorf("Start Dockyard unix socket error: %s\n", err.Error())
				}
			}
			break
		default:
			log.Fatalf("Invalid listen mode: %s\n", common.Web.Mode)
			os.Exit(1)
			break
		}
	}()

	// Graceful shutdown
	<-stopChan // wait for SIGINT
	log.Errorln("Shutting down server...")

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}

	log.Errorln("Server gracefully stopped")
}

// stopDaemon() stop Dockyard's REST API daemon.
func stopDaemon(cmd *cobra.Command, args []string) {

}

// monitorDaemon() monitor Dockyard's REST API deamon.
func monitorDaemon(cmd *cobra.Command, args []string) {

}
