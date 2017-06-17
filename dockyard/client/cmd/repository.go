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

var domain string
var namespace, repository, repoType string

// repository sub command
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "repository sub command create/delete and other manage function.",
	Long:  ``,
}

// create repository command
var createRepositoryCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a repository in Dockyard",
	Long:  ``,
	Run:   createRepository,
}

// init()
func init() {
	RootCmd.AddCommand(repositoryCmd)

	repositoryCmd.AddCommand(repositoryCmd)
	repositoryCmd.Flags().StringVar(&namespace, "user", "u", "Username or Organization for repository")
	repositoryCmd.Flags().StringVar(&repository, "repo", "r", "Repository name")
	repositoryCmd.Flags().StringVar(&repoType, "type", "t", "Repository type")
}

// createRepository is
func createRepository(cmd *cobra.Command, args []string) {

}
