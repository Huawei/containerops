package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "",
	Long:  ``,
	Run:   _deploy,
}

func _deploy(cmd *cobra.Command, args []string) {
	fmt.Println("" + address)

	//get iplist while
	// download.Download_main()
	// deploy.DeployNodes()

}
