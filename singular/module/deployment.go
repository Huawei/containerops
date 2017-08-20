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
	"path"
	"time"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/singular/model"
	"github.com/Huawei/containerops/singular/module/infras"
	"github.com/Huawei/containerops/singular/module/objects"
	"github.com/Huawei/containerops/singular/module/service"
	"github.com/Huawei/containerops/singular/module/tools"
)

const (
	InfraEtcd       = "etcd"
	InfraFlannel    = "flannel"
	InfraDocker     = "docker"
	InfraKubernetes = "kubernetes"
)

// Deploy Sequence: Preparing SSH Key files -> Preparing VM -> Preparing SSL root Key files -> Deploy Etcd
//   -> Deploy flannel -> Deploy k8s Master -> Deploy k8s node -> TODO Deploy other...
// Parameters:
//   db [bool] Save Deploy model and log in the database.
func DeployInfraStacks(d *objects.Deployment, db bool) error {
	var err error

	d.Log("Start deploy Cloud Native stack.")

	// Open Database Connection
	if db == true {
		if err := model.OpenDatabase(&common.Database); err != nil {
			return err
		}

		if err := model.Migrate(); err != nil {
			return err
		}

		d.Log("Open database connection and migrate table struct.")
	}

	// Preparing SSH Keys
	if d.Tools.SSH.Private != "" {
		if d.Tools.SSH.Public, d.Tools.SSH.Private, d.Tools.SSH.Fingerprint, err = tools.OpenSSHKeyFiles(d.Tools.SSH.Public, d.Tools.SSH.Private); err != nil {
			return err
		}
		d.Log(fmt.Sprintf("Open SSH private key and get fingerprint: %s", d.Tools.SSH.Fingerprint))
	} else {
		if d.Tools.SSH.Public, d.Tools.SSH.Private, d.Tools.SSH.Fingerprint, err = tools.GenerateSSHKeyFiles(d.Config); err != nil {
			return err
		}
		d.Log(fmt.Sprintf("Generate SSH files and get fingerprint: %s", d.Tools.SSH.Fingerprint))
	}

	// Preparing VM from service provider like DigitalOcean, GCE, AWS, Azure ...
	if d.Service.Provider != "" {
		switch d.Service.Provider {
		case "digitalocean":
			do := new(service.DigitalOcean)
			do.Token = d.Service.Token
			do.Region, do.Size, do.Image = d.Service.Region, d.Service.Size, d.Service.Image

			// TODO parse do.Image get distro

			// Init DigitalOcean API client.
			do.InitClient()
			d.Log("Init DigitalOcean API client.")

			// Upload ssh public key
			if err := do.UpdateSSHKey(d.Tools.SSH.Public); err != nil {
				return err
			}
			d.Log(fmt.Sprintf("Uploaded SSH public key to the DigitalOcean, and fingerprint is: %s.", d.Tools.SSH.Fingerprint))

			// Create DigitalOcean Droplets
			d.Log("Creating Droplets in DigitalOcean.")
			if err := do.CreateDroplet(d.Service.Nodes, d.Tools.SSH.Fingerprint); err != nil {
				return err
			}
			d.Log(fmt.Sprintf("Created %d Droplets successfully.", d.Service.Nodes))

			// Export droplets IP
			for ip, _ := range do.Droplets {
				node := objects.Node{
					IP:     ip,
					User:   service.DORootUser,
					Distro: tools.DistroUbuntu,
				}

				d.Nodes = append(d.Nodes, node)

				d.Log(fmt.Sprintf("Droplet Node IP: %s", ip))
			}

			d.Log("Waiting 60 seconds for preparing droplets.")
			time.Sleep(60 * time.Second)
		default:
			return fmt.Errorf("Unsupport service provide: %s", d.Service.Provider)
		}
	}

	if len(d.Nodes) > 0 {
		// Set Out Node Parameters
		i := 0
		for _, node := range d.Nodes {
			d.Log(fmt.Sprintf("Node %d ip: %s", i, node.IP))
			d.Output(fmt.Sprintf("NODE_%d", i), node.IP)
			i += 1
		}

		// Initialization node environment
		for _, node := range d.Nodes {
			if err := tools.InitializationEnvironment(d.Tools.SSH.Private, node.IP, node.User, node.Distro); err != nil {
				return err
			}
		}

		// Generate CA CARootTemplate files
		if roots, err := tools.GenerateCARootFiles(d.Config); err != nil {
			return err
		} else {
			d.Log(fmt.Sprintf("CA CARootTemplate files generated successfully in : %s",
				path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder)))

			for key, value := range roots {
				d.Output(key, value)
			}

			for _, node := range d.Nodes {
				d.Log(fmt.Sprintf("Upload SSL CARootTempate files to Droplet[%s] and init environments.", node.IP))
				if err := tools.UploadCARootFiles(d.Tools.SSH.Private, roots, node.IP, node.User); err != nil {
					return err
				}
			}
		}

		// Deploy infras
		for _, infra := range d.Infras {
			switch infra.Name {
			case InfraEtcd:
				if err := infras.DeployEtcdInCluster(d, &infra); err != nil {
					return err
				}
			case InfraFlannel:
				if err := infras.DeployFlannelInCluster(d, &infra); err != nil {
					return err
				}
			case InfraDocker:
				if err := infras.DeployDockerInCluster(d, &infra); err != nil {
					return err
				}
			case InfraKubernetes:
				if err := infras.DeployKubernetesInCluster(d, &infra); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupport infrastruction software: %s", infra)
			}
		}

	}

	return err
}
