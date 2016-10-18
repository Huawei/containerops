package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is root cmd of pilotage.
var RootCmd = &cobra.Command{
	Use:   "pilotage",
	Short: "pilotage is a DevOps workflow engine, both for daemon and client.",
	Long: `Pilotage is a DevOps workflow engine with customizable DevOps component repository with container, and it's core project of ContainerOps.  
ContainerOps is whole new concept of DevOps with DevOps workflow engine and components.`,
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

}
