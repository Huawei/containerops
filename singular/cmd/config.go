package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage Singular configurations",
	Long: `Singular manage deploy templates, systemd template files and
access tokens and others in the config file and config folders.`,
}

// init()
func init() {
	RootCmd.AddCommand(configCmd)
}
