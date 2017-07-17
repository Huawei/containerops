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

	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/pilotage/middleware"
	"github.com/Huawei/containerops/pilotage/router"
)

var addressOption, webMode string
var keyFile, certFile string
var portOption int

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
	// Add cli daemon sub command.
	RootCmd.AddCommand(daemonCmd)

	daemonCmd.Flags().StringVarP(&addressOption, "address", "a", "localhost", "The daemon listen address.")
	daemonCmd.Flags().StringVarP(&certFile, "cert", "c", "", "The cert file for HTTPS mode")
	daemonCmd.Flags().StringVarP(&keyFile, "key", "k", "", "The key file for HTTPS mode")
	daemonCmd.Flags().StringVarP(&webMode, "mode", "m", "http", "The http mode")
	daemonCmd.Flags().IntVarP(&portOption, "port", "p", 8080, "The port of http.")

	viper.BindPFlag("address", daemonCmd.Flags().Lookup("address"))
	viper.BindPFlag("cert", daemonCmd.Flags().Lookup("cert"))
	viper.BindPFlag("key", daemonCmd.Flags().Lookup("key"))
	viper.BindPFlag("mode", daemonCmd.Flags().Lookup("mode"))
	viper.BindPFlag("port", daemonCmd.Flags().Lookup("port"))

	daemonCmd.AddCommand(runDaemonCmd)

	daemonCmd.AddCommand(startDaemonCmd)
}

// runDaemonFlow is
func runDaemonFlow(cmd *cobra.Command, args []string) {
	if len(args) <= 0 || utils.IsFileExist(args[0]) == false {
		cmd.Println(Red("The orchestration flow file is required."))
		os.Exit(1)
	}

	m := macaron.New()
	middleware.SetRunDaemonMiddlewares(m, cfgFile, args[0])
	router.SetRunDaemonRouters(m)

	var server *http.Server

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		switch webMode {
		case "http":
			listenAddr := fmt.Sprintf("%s:%d", addressOption, portOption)
			if err := http.ListenAndServe(listenAddr, m); err != nil {
				cmd.Println(Red("Start pilotage http daemon error:"), Red(err.Error()))
			}

			break
		case "https":
			if certFile != "" || keyFile != "" {
				cmd.Println(Red("HTTPS mode need TLS cert and key files"))
			}

			listenAddr := fmt.Sprintf("%s:%d", addressOption, portOption)
			server = &http.Server{Addr: listenAddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
			if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
				cmd.Println(Red("Start pilotage https daemon error: "), Red(err.Error()))
			}

			break
		case "unix":
			listenAddr := fmt.Sprintf("%s", addressOption)
			if utils.IsFileExist(listenAddr) {
				os.Remove(listenAddr)
			}

			if listener, err := net.Listen("unix", listenAddr); err != nil {
				cmd.Println(Red("Start pilotage Unix Socket daemon error: "), Red(err.Error()))
			} else {
				server = &http.Server{Handler: m}
				if err := server.Serve(listener); err != nil {
					cmd.Println(Red("Start pilotage Unix Socket error: "), Red(err.Error()))
				}
			}
			break
		default:
			cmd.Println(Red("Invalid listen mode: "), Red(webMode))
			os.Exit(1)
			break
		}
	}()

	// Graceful shutdown
	<-stopChan // wait for SIGINT
	cmd.Println(Green("Shutting down pilotage daemon server..."))

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}

	cmd.Println(Green(("pilotage daemon gracefully stopped")))

}

// startDaemonFlow is
func startDaemonFlow(cmd *cobra.Command, args []string) {
	if len(args) <= 0 || utils.IsFileExist(args[0]) == false {
		cmd.Println(Red("The orchestration flow file is required."))
		os.Exit(1)
	}

	m := macaron.New()
	middleware.SetStartDaemonMiddlewares(m, cfgFile)
	router.SetStartDaemonRouters(m)
}
