package cmd

import "github.com/spf13/cobra"

var _ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your nodes of kubernetes cluster",
	Long:  ``,
	Run:   _config,
}

func _config(cmd *cobra.Command, args []string) {

}
