package cmd

import "github.com/spf13/cobra"

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitored running ContainerOps modules, Kubernetes cluster, etcd cluster, flannel network and others.",
	Long: `Monitored the running ContainerOps modules, Kubernetes cluster, etcd cluster and others.
Singular don't monitored the applications in the Kubernetes, only cluster status.'`,
}

// init()
func init() {
	RootCmd.AddCommand(monitorCmd)
}
