package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/macaron.v1"

	"github.com/containerops/configure"
	"github.com/containerops/pilotage/models"
	"github.com/containerops/pilotage/utils"
	"github.com/containerops/pilotage/web"
)

var address string
var port int64

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

// startDeamon() start Pilotage's REST API daemon.
func startDeamon(cmd *cobra.Command, args []string) {
	// first open database conn
	models.OpenDatabase()

	m := macaron.New()

	// Set Macaron Web Middleware And Routers
	web.SetPilotagedMacaron(m)

	listenMode := configure.GetString("listenmode")
	switch listenMode {
	case "http":
		listenaddr := fmt.Sprintf("%s:%d", address, port)
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			fmt.Printf("Start Pilotage http service error: %v\n", err.Error())
		}
		break
	case "https":
		listenaddr := fmt.Sprintf("%s:443", address)
		server := &http.Server{Addr: listenaddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		if err := server.ListenAndServeTLS(configure.GetString("httpscertfile"), configure.GetString("httpskeyfile")); err != nil {
			fmt.Printf("Start Pilotage https service error: %v\n", err.Error())
		}
		break
	case "unix":
		listenaddr := fmt.Sprintf("%s", address)
		if utils.IsFileExist(listenaddr) {
			os.Remove(listenaddr)
		}

		if listener, err := net.Listen("unix", listenaddr); err != nil {
			fmt.Printf("Start Pilotage unix socket error: %v\n", err.Error())

		} else {
			server := &http.Server{Handler: m}
			if err := server.Serve(listener); err != nil {
				fmt.Printf("Start Pilotage unix socket error: %v\n", err.Error())
			}
		}
		break
	default:
		break
	}
}

// stopDaemon() stop Pilotage's REST API daemon.
func stopDaemon(cmd *cobra.Command, args []string) {

}

// monitordAemon() monitor Pilotage's REST API deamon.
func monitorDaemon(cmd *cobra.Command, args []string) {

}
