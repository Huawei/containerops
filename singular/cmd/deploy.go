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
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module"
	"github.com/Huawei/containerops/singular/module/objects"
)

var privateKey, publicKey, output string
var db, del bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy ContainerOps, Kubernetes, etcd, flannel and others.",
	Long: `Deploy ContainerOps modules, Kubernetes masters and nodes, etcd cluster, flannel plugin,
CoreDNS, Prometheus and others in Cloud Native stack.`,
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "deploy stack with a template file",
	Long:  `Deploy ContainerOps modules or others with deploy template file, the template.`,
	Run:   templateDeploy,
}

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "deploy template files in folder",
	Long:  "Deploy all template files in the folder",
	Run:   folderDeploy,
}

//init()
func init() {
	RootCmd.AddCommand(deployCmd)

	//Add sub command.
	deployCmd.AddCommand(templateCmd)
	deployCmd.AddCommand(folderCmd)

	templateCmd.Flags().StringVarP(&privateKey, "private-key", "i", "", "ssh identity file")
	templateCmd.Flags().StringVarP(&publicKey, "public-key", "p", "", "ssh public identity file")
	templateCmd.Flags().StringVarP(&output, "output", "o", "", "output data folder")
	templateCmd.Flags().BoolVarP(&db, "db", "d", false, "save deploy data in database.")
	templateCmd.Flags().BoolVarP(&del, "delete", "r", false, "del the nodes when deploy and test done.")

	viper.BindPFlag("private-key", templateCmd.Flags().Lookup("private-key"))
	viper.BindPFlag("public-key", templateCmd.Flags().Lookup("public-key"))
	viper.BindPFlag("output", templateCmd.Flags().Lookup("output"))
	viper.BindPFlag("db", templateCmd.Flags().Lookup("db"))
	viper.BindPFlag("delete", templateCmd.Flags().Lookup("del"))
}

//Deploy the Cloud Native stack with a template file.
func templateDeploy(cmd *cobra.Command, args []string) {
	//Check deploy template file.
	if len(args) <= 0 || utils.IsFileExist(args[0]) == false {
		fmt.Fprintf(os.Stderr, "the deploy template file is required, %s\n", "see https://github.com/Huawei/containerops/singular for more detail.")
		os.Exit(1)
	}

	template := args[0]
	d := new(objects.Deployment)

	//Read template file and parse.
	if err := d.ParseFromFile(template, output); err != nil {
		fmt.Fprintf(os.Stderr, "parse deploy template error: %s\n", err.Error())
		os.Exit(1)
	}

	//Set private key file path.
	if privateKey != "" {
		d.Tools.SSH.Private, d.Tools.SSH.Public = privateKey, publicKey
	}

	//The integrity checking of deploy template.
	if err := d.Check(); err != nil {
		fmt.Fprintf(os.Stderr, "parse deploy template error: %s\n", err.Error())
		os.Exit(1)
	}

	//Set log and error io.Writer
	var logWriters io.Writer

	//Generate stdout/stderr io.Writer
	stdoutFile, _ := os.Create(path.Join(d.Config, "deploy.log"))
	defer stdoutFile.Close()

	//Using MultiWriter log and error.
	if verbose == true {
		logWriters = io.MultiWriter(stdoutFile, os.Stdout)
	} else {
		logWriters = io.MultiWriter(stdoutFile)
	}

	//Deploy cloud native stack
	if err := module.DeployInfraStacks(d, db, logWriters, timestamp); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	//Delete droplets
	if del == true {
		if err := module.DeleteInfraStacks(d, logWriters, timestamp); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}
}

//Deploy the Cloud Native stack with a template file.
func folderDeploy(cmd *cobra.Command, args []string) {
	//Check deploy template file.
	if len(args) <= 0 || utils.IsFileExist(args[0]) == false {
		fmt.Fprintf(os.Stderr, "the template files path is required, %s\n", "see https://github.com/Huawei/containerops/singular for more detail.")
		os.Exit(1)
	}

	if utils.IsDirExist(args[0]) == false {
		fmt.Fprintf(os.Stderr, "the %s folder is not exist\n", args[0])
		os.Exit(1)
	}

	if files, err := utils.WalkDir(args[0], "*.yml"); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	} else {
		for _, file := range files {
			if utils.IsFileExist(file) == false {
				fmt.Fprintf(os.Stderr, "the %s file is not exist\n", file)
			}

			d := new(objects.Deployment)

			//Read template file and parse.
			if err := d.ParseFromFile(file, output); err != nil {
				fmt.Fprintf(os.Stderr, "parse deploy template error: %s\n", err.Error())
				os.Exit(1)
			}

			//Set private key file path.
			if privateKey != "" {
				d.Tools.SSH.Private, d.Tools.SSH.Public = privateKey, publicKey
			}

			//The integrity checking of deploy template.
			if err := d.Check(); err != nil {
				fmt.Fprintf(os.Stderr, "parse deploy template error: %s\n", err.Error())
			}

			//Set log and error io.Writer
			var logWriters io.Writer

			//Generate stdout/stderr io.Writer
			stdoutFile, _ := os.Create(path.Join(d.Config, "deploy.log"))
			defer stdoutFile.Close()

			//Using MultiWriter log and error.
			if verbose == true {
				logWriters = io.MultiWriter(stdoutFile, os.Stdout)
			} else {
				logWriters = io.MultiWriter(stdoutFile)
			}

			//Deploy cloud native stack
			if err := module.DeployInfraStacks(d, db, logWriters, timestamp); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			//Delete droplets
			if del == true {
				if err := module.DeleteInfraStacks(d, logWriters, timestamp); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					os.Exit(1)
				}
			}
		}
	}
}
