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
	"io"
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
func DeployDockerInCluster(d *objects.Deployment, infra *objects.Infra, stdout io.Writer, timestamp bool) error {
	nodes := []objects.Node{}

	for i := 0; i < infra.Master; i++ {
		nodes = append(nodes, d.Nodes[i])
	}

	objects.WriteLog(fmt.Sprintf("Docker nodes is %v", nodes), stdout, timestamp, d, infra)

	//Generate Docker systemd file
	if files, err := generateDockerFiles(d.Config, nodes, infra.Version); err != nil {
		return err
	} else {
		objects.WriteLog(fmt.Sprintf("Docker CA/Systemd/config files: [%v]", files), stdout, timestamp, d, infra)
		objects.WriteLog(fmt.Sprintf("Upload Docker CA/Systemd/config files: [%v]", files), stdout, timestamp, d, infra)

		//Upload Docker Systemd file
		if err := uploadDockerFiles(files, d.Tools.SSH.Private, nodes, stdout, timestamp); err != nil {
			return err
		}
	}

	//Download Docker files
	for _, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, nodes, stdout, timestamp); err != nil {
			return err
		}

		//Run Docker before scripts
		if c.Before != "" {
			for _, node := range nodes {
				if err := beforeDockerExecute(d.Tools.SSH.Private, node.IP, c.Before, node.User); err != nil {
					return err
				}
			}
		}
	}

	//Start Docker daemon
	for _, node := range nodes {
		if err := startDockerDaemon(d.Tools.SSH.Private, node.IP, node.User); err != nil {
			return err
		}
	}

	//Run after script
	for _, c := range infra.Components {
		if c.After != "" {
			for _, node := range nodes {
				if err := afterDockerExecute(d.Tools.SSH.Private, node.IP, c.After, node.User); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//generateDockerFiles generate Docker systemd service file.
func generateDockerFiles(src string, nodes []objects.Node, version string) (map[string]map[string]string, error) {
	result := map[string]map[string]string{}

	//Preparing the SSL folder
	sslBase := path.Join(src, tools.CAFilesFolder, tools.CADockerFolder)
	if utils.IsDirExist(sslBase) == true {
		os.RemoveAll(sslBase)
	}
	os.MkdirAll(sslBase, os.ModePerm)

	//Preparing the Systemd folder
	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceDockerFolder)
	if utils.IsDirExist(serviceBase) == true {
		os.RemoveAll(serviceBase)
	}
	os.MkdirAll(serviceBase, os.ModePerm)

	//Loop the nodes, generate Docker systemd service files
	for _, node := range nodes {
		if utils.IsDirExist(path.Join(serviceBase, node.IP)) == false {
			os.MkdirAll(path.Join(serviceBase, node.IP), os.ModePerm)
		}

		if files, err := generateDockerServiceFile(version, path.Join(serviceBase, node.IP)); err != nil {
			return result, err
		} else {
			for k, v := range files {
				result[node.IP][k] = v
			}
		}
	}

	return result, nil
}

//generateDockerServiceFile generate the Docker systemd file
func generateDockerServiceFile(version, base string) (map[string]string, error) {
	var serviceTpl bytes.Buffer
	var err error

	files := map[string]string{}
	files[tools.ServiceDockerFile] = path.Join(base, tools.ServiceDockerFile)

	serviceTp := template.New("docker-systemd")
	serviceTp, err = serviceTp.Parse(t.DockerSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, nil)
	serviceTpFileBytes := serviceTpl.Bytes()

	err = ioutil.WriteFile(files[tools.ServiceDockerFile], serviceTpFileBytes, 0700)

	return files, err
}

//Upload docker systemd file
//TODO apt-get install -y bridge-utils aufs-tools cgroupfs-mount libltdl7
func uploadDockerFiles(files map[string]map[string]string, key string, nodes []objects.Node, stdout io.Writer, timestamp bool) error {
	for _, node := range nodes {
		if cmd, err := tools.DownloadComponent(files[node.IP][tools.ServiceDockerFile], path.Join(tools.SystemdServerPath, tools.ServiceDockerFile), node.IP, key, node.User, stdout); err != nil {
			return err
		} else {
			objects.WriteLog(
				fmt.Sprintf("upload %s to %s@%s with %s", files[node.IP][tools.CAEtcdCSRConfigFile], node.IP, path.Join(EtcdServerConfig, EtcdServerSSL, tools.CAEtcdCSRConfigFile), cmd),
				stdout, timestamp, &node)
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
