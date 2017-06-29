package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your nodes of kubernetes cluster",
	Long:  ``,
	Run:   _config,
}

func _config(cmd *cobra.Command, args []string) {

}

func _Cerkey(cmd *cobra.Command, args []string) {

	//ssh-keygen -t rsa
	//ssh-keygen -b 2048 -t rsa
	//ssh-keygen -b 2048 -t rsa -f ~/.ssh/name
	if len(args) > 0 {
		fmt.Println(args)

		LocalExecCMDparams("ssh-keygen", []string{"-b", "2048", "-t", "rsa", "-f", args[1]})
	} else {
		LocalExecCMDparams("ssh-keygen", []string{"-b", "2048", "-t", "rsa", "-f", "~/.ssh/id_rsa"})
	}
	// singular_cmd.ExecCPparams("/tmp/flanneld", "/usr/bin/flanneld")
	// singular_cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/libexec/flannel/"})
	fmt.Println("[singular] Generated Certificate Authority key and certificate.")
	fmt.Println("[singular] Created keys and certificates in \"/usr/singular/\"")
}
