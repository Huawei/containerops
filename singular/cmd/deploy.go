package cmd

import "github.com/spf13/cobra"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy ContainerOps, Kubernetes, etcd, flannel and others.",
	Long: `Deploy ContainerOps modules, Kubernetes masters and nodes, etcd cluster, flannel plugin,
CoreDNS, Prometheus and others in Cloud Native stack.`,
}

// init()
func init() {
	RootCmd.AddCommand(deployCmd)
}
