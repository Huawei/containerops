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

package module

import (
	"fmt"
	"io"
	"time"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/singular/model"
	"github.com/Huawei/containerops/singular/module/infras"
	"github.com/Huawei/containerops/singular/module/objects"
	"github.com/Huawei/containerops/singular/module/service"
	"github.com/Huawei/containerops/singular/module/tools"
)

const (
	//InfraEtcd is Etcd flag
	InfraEtcd = "etcd"
	//InfraFlannel is Flannel flg
	InfraFlannel = "flannel"
	//InfraDocker is Docekr flag
	InfraDocker = "docker"
	//InfraKubernetes is Kubernetes flag
	InfraKubernetes = "kubernetes"
)

//DeployInfraStacks deploy all infras in the cluster.
//Deploy Sequence:
//   Preparing SSH Key files ->
//   Preparing VM ->
//   Preparing SSL root Key files ->
//   Deploy Etcd ->
//   Deploy flannel ->
//   Deploy k8s Master ->
//   Deploy k8s node ->
//Parameters:
//   db [bool] Save Deploy model and log in the database.
func DeployInfraStacks(d *objects.Deployment, db bool, stdout io.Writer, timestamp bool) error {
	var err error

	//Open Database Connection and migrate table struct.
	if db == true {
		objects.WriteLog("Open database connection and migrate the table struct", stdout, timestamp, d)
		if err := model.OpenDatabase(&common.Database); err != nil {
			return err
		}

		if err := model.Migrate(); err != nil {
			return err
		}
	}

	//Preparing SSH Keys
	if d.Tools.SSH.Private != "" {
		objects.WriteLog(
			fmt.Sprintf("Calculate the public key and fingerprint using the private key file [%s]", d.Tools.SSH.Private),
			stdout, timestamp, d)
		if d.Tools.SSH.Public, d.Tools.SSH.Private, d.Tools.SSH.Fingerprint, err = tools.OpenSSHKeyFiles(d.Tools.SSH.Public, d.Tools.SSH.Private); err != nil {
			return err
		}
	} else {
		if d.Tools.SSH.Public, d.Tools.SSH.Private, d.Tools.SSH.Fingerprint, err = tools.GenerateSSHKeyFiles(d.Config); err != nil {
			return err
		}
		objects.WriteLog(fmt.Sprintf("Generate the SSH keys file [%s] for deploy cloud native stack", d.Tools.SSH.Private),
			stdout, timestamp, d)
	}

	//Preparing VM from service provider like DigitalOcean, GCE, AWS, Azure ...
	if d.Service.Provider != "" {
		switch d.Service.Provider {
		case "digitalocean":
			do := service.DigitalOcean{}
			do.Token = d.Service.Token
			do.Region, do.Size, do.Image = d.Service.Region, d.Service.Size, d.Service.Image
			//TODO parse do.Image get distro
			//Init DigitalOcean API client.
			objects.WriteLog("Create DigitalOcean droplets", stdout, timestamp, d, &do)

			objects.WriteLog("Init DigitalOcean client with token", stdout, timestamp, d, &do)
			do.InitClient()

			objects.WriteLog(
				fmt.Sprintf("Upload ssh public key file [%s] to DigitalOcean", d.Tools.SSH.Public),
				stdout, timestamp, d, &do)
			//Upload ssh public key
			if err := do.UploadSSHKey(d.Tools.SSH.Public); err != nil {
				return err
			}

			//Prepare droplet prefix name and tags
			namespace, repository, name, _ := d.URIs()
			tags := []string{namespace, repository, name, fmt.Sprintf("version-%d", d.Version), d.Tag}
			objects.WriteLog(fmt.Sprintf("Droplets tag is %v", tags), stdout, timestamp, d, &do)

			//Create DigitalOcean Droplets
			if err := do.CreateDroplets(d.Service.Nodes, d.Tools.SSH.Fingerprint, fmt.Sprintf("%s-%s", namespace, repository), tags); err != nil {
				return err
			}

			//Export droplets IP
			for id, droplet := range do.Droplets {
				node := objects.Node{
					ID:      id,
					IP:      droplet["public"],
					Private: droplet["private"],
					User:    service.DORootUser,
					Distro:  tools.DistroUbuntu,
				}

				objects.WriteLog(fmt.Sprintf("Droplet [%d] ip is [%s/%s]", id, droplet["public"], droplet["private"]), stdout, timestamp, d, &do)
				d.Nodes = append(d.Nodes, &node)
			}

			objects.WriteLog("sleep 60 second for preparing Droplets", stdout, timestamp, d, &do)
			time.Sleep(60 * time.Second)

			d.Service.Logs = do.Logs
		default:
			return fmt.Errorf("unsupport service provide: %s", d.Service.Provider)
		}
	}

	if len(d.Nodes) > 0 {
		//Set Out Node Parameters
		i := 0
		for _, node := range d.Nodes {
			d.Output(fmt.Sprintf("NODE_%d", i), node.IP)
			i += 1

			//Export node ip to stdout
			objects.WriteLog(fmt.Sprintf("[COUT] NODE_%d = %s", i, node.IP), stdout, timestamp, d)
		}

		//Initialization node environment
		for _, node := range d.Nodes {
			objects.WriteLog(fmt.Sprintf("Initializate the %s node environment", node.IP), stdout, timestamp, d, node)
			if commands, err := tools.InitializationEnvironment(d.Tools.SSH.Private, node.IP, node.User, node.Distro, stdout); err != nil {
				return err
			} else {
				objects.WriteLog(fmt.Sprintf("execute commands %v in the %d node", commands, node.ID), stdout, timestamp, d, node)
			}
		}

		//Generate CA CARootTemplate files
		if roots, err := tools.GenerateCARootFiles(d.Config); err != nil {
			return err
		} else {
			objects.WriteLog(fmt.Sprintf("Generate cluster root CA files."), stdout, timestamp, d)
			for key, value := range roots {
				d.Output(key, value)
				objects.WriteLog(fmt.Sprintf("[COUT] %s = %s", key, value), stdout, timestamp, d)
			}

			//Upload CA root files to the nodes.
			for _, node := range d.Nodes {
				objects.WriteLog(fmt.Sprintf("Upload root CA files %v to %s node", roots, node.IP), stdout, timestamp, d, node)
				if err := tools.UploadCARootFiles(d.Tools.SSH.Private, roots, node.IP, node.User, stdout); err != nil {
					return err
				}
			}
		}

		//Deploy infras
		for _, infra := range d.Infras {
			switch infra.Name {
			case InfraEtcd:
				objects.WriteLog(fmt.Sprintf("Deploy etcd in cluster"), stdout, timestamp, d, infra)
				if err := infras.DeployEtcdInCluster(d, infra, stdout, timestamp); err != nil {
					return err
				}
			case InfraFlannel:
				objects.WriteLog(fmt.Sprintf("Deploy flannel in cluster"), stdout, timestamp, d, infra)
				if err := infras.DeployFlannelInCluster(d, infra, stdout, timestamp); err != nil {
					return err
				}
			case InfraDocker:
				objects.WriteLog(fmt.Sprintf("Deploy docker in cluster"), stdout, timestamp, d, infra)
				if err := infras.DeployDockerInCluster(d, infra, stdout, timestamp); err != nil {
					return err
				}
			case InfraKubernetes:
				objects.WriteLog(fmt.Sprintf("Deploy k8s in cluster"), stdout, timestamp, d, infra)
				if err := infras.DeployKubernetesInCluster(d, infra, stdout, timestamp); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupport infrastruction software: %s", infra)
			}
		}

	}

	//Save data
	if db == true {
		if err := d.Save(); err != nil {
			return err
		}
	}

	return err
}

//DeleteInfraStacks delete the droplets of stack.
func DeleteInfraStacks(d *objects.Deployment, stdout io.Writer, timestamp bool) error {
	if d.Service.Provider != "" {
		switch d.Service.Provider {
		case "digitalocean":
			do := service.DigitalOcean{}
			do.Token = d.Service.Token
			do.InitClient()

			for _, node := range d.Nodes {
				if err := do.DeleteDroplet(node.ID); err != nil {
					return err
				}
				objects.WriteLog(fmt.Sprintf("Delete %d droplet in DigitalOcaen", node.ID), stdout, timestamp, d)
			}

		default:
			return fmt.Errorf("unsupport service provide: %s", d.Service.Provider)
		}
	}

	return nil
}
