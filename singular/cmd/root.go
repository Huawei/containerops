package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

// RootCmd is root cmd of dockyard.
var RootCmd = &cobra.Command{
	Use:   "singular",
	Short: "The deployment and operations tools.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("111")

	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// init()
func init() {
	RootCmd.AddCommand(versionCmd)
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	versionCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "http or https listen address.")
	versionCmd.Flags().Int64VarP(&port, "port", "p", 80, "the port of http.")
}

var address string
var port int64

var versionCmd = &cobra.Command{
	Use:   "install",
	Short: "Start to install kubenetes cluster automatically by the configuration file",
	Long:  `Start to install kubenetes cluster automatically by the configuration file`,
	Run:   startInstall,
}

func startInstall(cmd *cobra.Command, args []string) {
	fmt.Println("Start to install kubenetes cluster automatically by the configuration file" + address)
	// fmt.Println("Hugo Static Site Generator v0.9 -- HEAD  " + strings.Join(port, ""))

}
