package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,
	Run:   _Install,
}

func _Install(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)
}
