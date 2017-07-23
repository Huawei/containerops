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

	web "github.com/Huawei/containerops/assembling/web"
	log "github.com/Sirupsen/logrus"
	macaron "gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(daemonCmd)

	// Add start sub command
	daemonCmd.AddCommand(startDaemonCmd)

	// Add stop sub command
	daemonCmd.AddCommand(stopDaemonCmd)
}

// webCmd is sub command which start/stop/monitor Assembling's REST API daemon.
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Web sub command start/stop/monitor Assembling's REST API daemon.",
	Long:  ``,
}

// start Assembling daemon sub command
var startDaemonCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Assembling's REST API daemon.",
	Long:  ``,
	Run:   startDaemon,
}

// stop Assembling daemon sub command
var stopDaemonCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop Assembling's REST API daemon.",
	Long:  ``,
	Run:   stopDaemon,
}

var addressOption string
var portOption int

func startDaemon(cmd *cobra.Command, args []string) {
	m := macaron.New()

	web.SetAssemblingMacaron(m, "CONFIG_FILE_PATH")

	var server *http.Server
	stopChan := make(chan os.Signal)

	signal.Notify(stopChan, os.Interrupt)

	address := common.Assembling.Address
	port := common.Assembling.Port

	common.Assembling.Mode = "unix"

	go func() {
		switch common.Assembling.Mode {
		case "https":
			listenAddr := fmt.Sprintf("%s:%d", address, port)
			server = &http.Server{Addr: listenAddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
			if err := server.ListenAndServeTLS(common.Assembling.Cert, common.Assembling.Key); err != nil {
				log.Errorf("Start Assembling https service error: %s\n", err.Error())
			}

			break
		case "unix":
			listenAddr := fmt.Sprintf("%s", address)
			if utils.IsFileExist(listenAddr) {
				os.Remove(listenAddr)
			}

			if listener, err := net.Listen("unix", listenAddr); err != nil {
				log.Errorf("Start Assembling unix socket error: %s\n", err.Error())
			} else {
				server = &http.Server{Handler: m}
				if err := server.Serve(listener); err != nil {
					log.Errorf("Start Assembling unix socket error: %s\n", err.Error())
				}
			}
			break
		default:
			log.Fatalf("Invalid listen mode: %s\n", common.Assembling.Mode)
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

func stopDaemon(cmd *cobra.Command, args []string) {
	fmt.Println("stop assembling daemon, not implemented yet")
}
