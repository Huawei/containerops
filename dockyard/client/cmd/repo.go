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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/dockyard/client/repo"
)

var repoType string

// repo sub command
var repositoryCmd = &cobra.Command{
	Use:   "repo",
	Short: "Create/delete and other manage repository.",
	Long: `When using Dockyard as binary repository, should create a repository with
binary type before uploading file.
    
    warship repo create --type binary containerops/cncf-demo
    
`,
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

	if repoType == "" {
		fmt.Println("The repository must be `docker` or `binary`, not be a null value.")
		os.Exit(1)
	}

	if args[0] == "" {
		fmt.Println("The repository name and namespace is required.")
		os.Exit(1)
	}

	namespace := strings.Split(args[0], "/")[0]
	repository := strings.Split(args[0], "/")[1]

	if err := repo.CreateRepository(domain, namespace, repository, repoType); err != nil {
		fmt.Println(fmt.Sprintf("Create repository %s error: %s", args[0], err.Error()))
		os.Exit(1)
	}

	fmt.Println("Create repository successfully.")
}
