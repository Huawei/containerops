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

package infras

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module/objects"
	t "github.com/Huawei/containerops/singular/module/template"
	"github.com/Huawei/containerops/singular/module/tools"
)

//DeployDockerInCluster deploy Docker in Cluster
func DeployDockerInCluster(d *objects.Deployment, infra *objects.Infra) error {
	dockerNodes := map[string]string{}
	for i := 0; i < infra.Master; i++ {
		dockerNodes[fmt.Sprintf("docker-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
	}

	//Generate Docker systemd file
	if err := generateDockerFiles(d.Config, dockerNodes, infra.Version); err != nil {
		return err
	}

	//Upload Docker Systemd file
	if err := uploadDockerFiles(d.Config, d.Tools.SSH.Private, dockerNodes, tools.DefaultSSHUser); err != nil {
		return err
	}

	//Download Docker files
	for _, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, dockerNodes); err != nil {
			return err
		}

		//Run Docker before scripts
		if c.Before != "" {
			for _, ip := range dockerNodes {
				if err := beforeDockerExecute(d.Tools.SSH.Private, ip, c.Before, tools.DefaultSSHUser); err != nil {
					return err
				}
			}
		}
	}

	//Start Docker daemon
	for _, ip := range dockerNodes {
		if err := startDockerDaemon(d.Tools.SSH.Private, ip, tools.DefaultSSHUser); err != nil {
			return err
		}
	}

	//Run after script
	for _, c := range infra.Components {
		if c.After != "" {
			for _, ip := range dockerNodes {
				if err := afterDockerExecute(d.Tools.SSH.Private, ip, c.After, tools.DefaultSSHUser); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//generateDockerFiles generate Docker systemd service file.
func generateDockerFiles(src string, nodes map[string]string, version string) error {
	sslBase := path.Join(src, tools.CAFilesFolder, tools.CADockerFolder)
	if utils.IsDirExist(sslBase) == true {
		os.RemoveAll(sslBase)
	}

	os.MkdirAll(sslBase, os.ModePerm)

	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceDockerFolder)
	if utils.IsDirExist(serviceBase) == true {
		os.RemoveAll(serviceBase)
	}

	os.MkdirAll(serviceBase, os.ModePerm)

	for _, ip := range nodes {
		if utils.IsDirExist(path.Join(serviceBase, ip)) == false {
			os.MkdirAll(path.Join(serviceBase, ip), os.ModePerm)
		}

		var serviceTpl bytes.Buffer
		var err error

		serviceTp := template.New("docker-systemd")
		serviceTp, err = serviceTp.Parse(t.DockerSystemdTemplate[version])
		serviceTp.Execute(&serviceTpl, nil)
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(path.Join(sslBase, ip, tools.ServiceDockerFile), serviceTpFileBytes, 0700)

		return err
	}

	return nil
}

// Upload docker systemd file
// TODO apt-get install -y bridge-utils aufs-tools cgroupfs-mount libltdl7
func uploadDockerFiles(src, key string, nodes map[string]string, user string) error {
	sslBase := path.Join(src, tools.CAFilesFolder, tools.CADockerFolder)
	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceDockerFolder)

	if utils.IsDirExist(sslBase) == false || utils.IsDirExist(serviceBase) {
		return fmt.Errorf("locate docker folders %s or %s error", sslBase, serviceBase)
	}

	for _, ip := range nodes {
		if err := tools.DownloadComponent(path.Join(serviceBase, ip, tools.ServiceDockerFile), path.Join(tools.SytemdServerPath, tools.ServiceDockerFile), ip, key, user); err != nil {
			return err
		}
	}

	return nil
}

// Execute before script
func beforeDockerExecute(key, ip, cmd, user string) error {
	if err := utils.SSHCommand(user, key, ip, tools.DefaultSSHPort, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

// Start docker daemon
func startDockerDaemon(key, ip, user string) error {
	cmd := "systemctl daemon-reload && systemctl enable docker && systemctl start --no-block docker"

	if err := utils.SSHCommand(user, key, ip, tools.DefaultSSHPort, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

// Execute after daemon start
func afterDockerExecute(key, ip, cmd, user string) error {
	if err := utils.SSHCommand(user, key, ip, tools.DefaultSSHPort, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}
