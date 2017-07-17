/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

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
