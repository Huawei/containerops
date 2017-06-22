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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Huawei/containerops/common"
)

var repoType string

// repository sub command
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "repository sub command which create/delete and other manage repository.",
	Long: `When using Dockyard as binary repository, should create a repository with
binary type before uploading file.`,
}

// create repository command
var createRepositoryCmd = &cobra.Command{
	Use:   "create",
	Short: "Namespace and repository as the args like `containerops/cncf-demo`.",
	Long:  `There are two repository types support just now, it's docker and binary`,
	Run:   createRepository,
}

// init()
func init() {
	RootCmd.AddCommand(repositoryCmd)

	//Add create repository sub command.
	repositoryCmd.AddCommand(createRepositoryCmd)
	createRepositoryCmd.Flags().StringVarP(&repoType, "type", "t", "", "Repository type")
}

// createRepository is
func createRepository(cmd *cobra.Command, args []string) {
	if domain == "" {
		domain = common.Warship.Domain
	}

	fmt.Println(domain)
}
