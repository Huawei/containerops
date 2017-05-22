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

func Cerkey(cmd *cobra.Command, args []string) {
	//  124  ssh-keygen -t rsa
	//ssh-keygen -b 2048 -t rsa
	//ssh-keygen -b 2048 -t rsa -f ~/.ssh/lidian
	LocalExecCMDparams("ssh-keygen", []string{"-b", "2048", "-t", "rsa", "-f", "~/.ssh/lidian111"})

	LocalExecCMDparams("ssh-keygen", []string{"-b", "2048", "-t", "rsa", "-f", args[1]})

	// singular_cmd.ExecCPparams("/tmp/flanneld", "/usr/bin/flanneld")
	// singular_cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/libexec/flannel/"})
	fmt.Println("[singular] Generated Certificate Authority key and certificate.")
	fmt.Println("[singular] Created keys and certificates in \"/usr/singular/\"")
	fmt.Println(args)
}
