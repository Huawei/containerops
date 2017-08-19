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
	////"io/ioutil"
	////"net/url"
	////"os"
	////"os/exec"
	"path"
	"time"

	////"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/singular/model"
	"github.com/Huawei/containerops/singular/module/objects"
	"github.com/Huawei/containerops/singular/module/service"
	"github.com/Huawei/containerops/singular/module/tools"
	//"github.com/Huawei/containerops/singular/module/service"
	////"github.com/Huawei/containerops/singular/module/template"
	//"github.com/Huawei/containerops/singular/module/tools"
	"github.com/Huawei/containerops/singular/module/infras"
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
				// Deploy Etcd Cluster
				if err := infras.DeployEtcdCluster(d, &infra); err != nil {
					return err
				}
			//case InfraFlannel:
			//	if err := d.DeployFlannel(infra); err != nil {
			//		return err
			//	}
			//case InfraDocker:
			//	if err := d.DeployDocker(infra); err != nil {
			//		return err
			//	}
			//case InfraKubernetes:
			//	if err := d.DeployKubernetes(infra); err != nil {
			//		return err
			//	}
			default:
				return fmt.Errorf("Unsupport infrastruction software: %s", infra)
			}
		}

	}

	return err
}

//func (d *objects.Deployment) Deploy(db bool) error {

// Deploy infra
//for _, infra := range d.Infras {
//	switch infra.Name {
//	case InfraEtcd:
//		if err := infras.DeployEtcd(&d, infra); err != nil {
//			return err
//		}
//	case InfraFlannel:
//		if err := d.DeployFlannel(infra); err != nil {
//			return err
//		}
//	case InfraDocker:
//		if err := d.DeployDocker(infra); err != nil {
//			return err
//		}
//	case InfraKubernetes:
//		if err := d.DeployKubernetes(infra); err != nil {
//			return err
//		}
//	default:
//		return fmt.Errorf("Unsupport infrastruction software: %s", infra)
//	}
//}
//	}
//
//	return nil
//}

//
//func (d *Deployment) DeployFlannel(infra Infra) error {
//	flanneldNodes := map[string]string{}
//	for i := 0; i < infra.Master; i++ {
//		flanneldNodes[fmt.Sprintf("flanneld-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
//	}
//
//	d.Log(fmt.Sprintf("Generating SSL files and systemd service file for Flanneld."))
//	if err := GenerateFlanneldFiles(d.Config, flanneldNodes, d.Outputs["EtcdEndpoints"].(string), infra.Version); err != nil {
//		return err
//	} else {
//		d.Log(fmt.Sprintf("Uploading SSL files and systemd service to nodes of Flanneld."))
//		if err := UploadFlanneldCAFiles(d.Config, d.Tools.SSH.Private, flanneldNodes); err != nil {
//			return err
//		}
//
//		for i, c := range infra.Components {
//			d.Log(fmt.Sprintf("Downloading Flanneld binary files to Nodes."))
//			if err := d.DownloadBinaryFile(c.Binary, c.URL, flanneldNodes); err != nil {
//				return err
//			}
//
//			if c.Before != "" && i == 0 {
//				d.Log(fmt.Sprintf("Execute Flanneld before scripts: %s", c.Before))
//				if err := BeforeFlanneldExecute(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", i)].(string), c.Before, d.Outputs["EtcdEndpoints"].(string)); err != nil {
//					return err
//				}
//			}
//		}
//
//		d.Log(fmt.Sprintf("Staring Flanneld Service."))
//		if err := StartFlanneldCluster(d.Tools.SSH.Private, flanneldNodes); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (d *Deployment) DeployDocker(infra Infra) error {
//	dockerNodes := map[string]string{}
//	for i := 0; i < infra.Master; i++ {
//		dockerNodes[fmt.Sprintf("docker-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
//	}
//
//	d.Log(fmt.Sprintf("Generating SSL files and systemd service file for Docker."))
//	if err := GenerateDockerFiles(d.Config, dockerNodes, infra.Version); err != nil {
//		return err
//	} else {
//		d.Log(fmt.Sprintf("Uploading SSL files and systemd service to nodes of Docker."))
//		if err := UploadDockerCAFiles(d.Config, d.Tools.SSH.Private, dockerNodes); err != nil {
//			return err
//		}
//
//		for _, c := range infra.Components {
//			d.Log(fmt.Sprintf("Downloading Docker binary files to Nodes."))
//			if err := d.DownloadBinaryFile(c.Binary, c.URL, dockerNodes); err != nil {
//				return err
//			}
//
//			if c.Before != "" {
//				for _, ip := range dockerNodes {
//					d.Log(fmt.Sprintf("Execute Docker before scripts: %s in %s", ip, c.Before))
//					if err := BeforeDockerExecute(d.Tools.SSH.Private, ip, c.Before); err != nil {
//						return err
//					}
//				}
//			}
//		}
//
//		for _, ip := range dockerNodes {
//			d.Log(fmt.Sprintf("Start Docker in %s", ip))
//			if err := StartDockerDaemon(d.Tools.SSH.Private, ip); err != nil {
//				return err
//			}
//		}
//
//		for _, c := range infra.Components {
//			if c.After != "" {
//				for _, ip := range dockerNodes {
//					d.Log(fmt.Sprintf("Execute Docker After scripts: %s in %s", c.After, ip))
//					if err := AfterDockerExecute(d.Tools.SSH.Private, ip, c.After); err != nil {
//						return err
//					}
//				}
//			}
//		}
//
//	}
//
//	return nil
//}
//
//// DeployKubernetes is function deployment Kubernetes cluster include master and nodes.
//// Notes:
////   1. Kubernetes master cluster IP.
////   2. Set kubectl config files.
////   3. Deploy Kubernetes master.
////   4. Deploy Kubernetes nodes.
//func (d *Deployment) DeployKubernetes(infra Infra) error {
//	// TODO Now singular only support one master and multiple nodes architect.
//	// TODO So we decide the Kubernetes master IP is NODE_0 .
//	masterIp := d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)
//	etcdEndpoints := d.Outputs["EtcdEndpoints"].(string)
//
//	d.Output("MASTER_IP", masterIp)
//	d.Output("KUBE_APISERVER", fmt.Sprintf("https://%s:6443", masterIp))
//
//	kubeMasterNodes := map[string]string{}
//	for i := 0; i < infra.Master; i++ {
//		kubeMasterNodes[fmt.Sprintf("kube-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
//	}
//
//	kubeSlaveNodes := map[string]string{}
//	for i := 0; i < infra.Minion; i++ {
//		kubeSlaveNodes[fmt.Sprintf("kube-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
//	}
//
//	for _, c := range infra.Components {
//		d.Log(fmt.Sprintf("Download %s binary files", c.Binary))
//		if err := d.DownloadBinaryFile(c.Binary, c.URL, kubeSlaveNodes); err != nil {
//			return err
//		}
//	}
//
//	for _, c := range infra.Components {
//		if c.Binary == "kubectl" {
//			if utils.IsDirExist(path.Join(d.Config, "kubectl")) == true {
//				os.RemoveAll(path.Join(d.Config, "kubectl"))
//			}
//
//			os.MkdirAll(path.Join(d.Config, "kubectl"), os.ModePerm)
//
//			d.Log("Downloading kubectl binary file")
//
//			if a, err := url.Parse(c.URL); err != nil {
//				return err
//			} else {
//				if a.Scheme == "" {
//					cmdCopy := exec.Command("cp", c.URL, fmt.Sprintf("%s/kubectl/kubectl", d.Config))
//
//					cmdCopy.Stdout, cmdCopy.Stderr = os.Stdout, os.Stderr
//					if err := cmdCopy.Run(); err != nil {
//						return err
//					}
//				} else {
//					cmdDownload := exec.Command("curl", c.URL, "-o", fmt.Sprintf("%s/kubectl/kubectl", d.Config))
//
//					cmdDownload.Stdout, cmdDownload.Stderr = os.Stdout, os.Stderr
//					if err := cmdDownload.Run(); err != nil {
//						return err
//					}
//				}
//			}
//
//			cmdChmod := exec.Command("chmod", "+x", fmt.Sprintf("%s/kubectl/kubectl", d.Config))
//			cmdChmod.Stdout, cmdChmod.Stderr = os.Stdout, os.Stderr
//			if err := cmdChmod.Run(); err != nil {
//				return err
//			}
//
//			d.Log("Generate kubernetes admin ca files")
//			if err := GenerateAdminCAFiles(d.Config); err != nil {
//				return err
//			}
//
//			d.Log("Generate kubectl config file")
//			cmdSetCluster := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
//				fmt.Sprintf("--certificate-authority=%s", path.Join(d.Config, "ssl", "root", "ca.pem")),
//				"--embed-certs=true",
//				fmt.Sprintf("--server=%s", d.Outputs["KUBE_APISERVER"].(string)))
//			cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
//			if err := cmdSetCluster.Run(); err != nil {
//				return err
//			}
//
//			cmdSetCredentials := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-credentials", "admin",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
//				fmt.Sprintf("--client-certificate=%s", path.Join(d.Config, "kubectl", "admin.pem")),
//				"--embed-certs=true",
//				fmt.Sprintf("--client-key=%s", path.Join(d.Config, "kubectl", "admin-key.pem")))
//			cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
//			if err := cmdSetCredentials.Run(); err != nil {
//				return err
//			}
//
//			cmdSetContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-context", "kubernetes",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
//				"--cluster=kubernetes", "--user=admin")
//			cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
//			if err := cmdSetContext.Run(); err != nil {
//				return err
//			}
//
//			cmdUseContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "use-context",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
//				"kubernetes")
//			cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
//			if err := cmdUseContext.Run(); err != nil {
//				return err
//			}
//
//			d.Log("Upload kubectl config file to Kubernetes nodes")
//			if err := UploadKubeConfigFiles(d.Config, d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
//				return err
//			}
//		}
//
//		if c.Binary == "kube-apiserver" {
//			d.Log("Generate Kuberentes Token API file")
//			if err := GenerateTokenFile(d.Config); err != nil {
//				return err
//			}
//
//			d.Log("Generate Kubernetes SSL files and systemd service file")
//			if err := GenerateKuberAPIServerCAFiles(d.Config, masterIp, etcdEndpoints, infra.Version); err != nil {
//				return err
//			}
//
//			d.Log("Upload Kubernetes Token file")
//			if err := UploadTokenFiles(d.Config, d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//
//			d.Log("Upload Kubernetes API Server SSL files and systemd service file")
//			if err := UploadKubeAPIServerCAFiles(d.Config, d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//
//			d.Log("Start Kubernetes API Server")
//			if err := StartKubeAPIServer(d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//
//		}
//
//		if c.Binary == "kube-controller-manager" {
//			d.Log("Generate Kube-controller-manager systemd service file")
//			if err := GenerateKuberControllerManagerFiles(d.Config, masterIp, etcdEndpoints, infra.Version); err != nil {
//				return err
//			}
//
//			d.Log("Upload Kuber-controller-manager systemd service file")
//			if err := UploadKuberControllerFiles(d.Config, d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//
//			d.Log("Start Kube-controller-manager")
//			if err := StartKuberController(d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//		}
//
//		if c.Binary == "kube-scheduler" {
//			d.Log("Generate Kube-scheduler systemd service file")
//			if err := GenerateKuberSchedulerManagerFiles(d.Config, masterIp, etcdEndpoints, infra.Version); err != nil {
//				return err
//			}
//
//			d.Log("Upload Kuber-scheduler systemd service file")
//			if err := UploadKuberSchedulerManagerFiles(d.Config, d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//
//			d.Log("Start Kube-scheduler")
//			if err := StartKuberSchedulerManager(d.Tools.SSH.Private, masterIp); err != nil {
//				return err
//			}
//		}
//
//		if c.Binary == "kubelet" {
//			d.Log("Generate bootstrap.kubeconfig config file")
//			cmdSetCluster := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
//				fmt.Sprintf("--certificate-authority=%s", path.Join(d.Config, "ssl", "root", "ca.pem")),
//				"--embed-certs=true",
//				fmt.Sprintf("--server=%s", d.Outputs["KUBE_APISERVER"].(string)))
//			cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
//			if err := cmdSetCluster.Run(); err != nil {
//				return err
//			}
//
//			cmdSetCredentials := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-credentials", "kubelet-bootstrap",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
//				fmt.Sprintf("--token=%s", template.BooststrapToken))
//			cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
//			if err := cmdSetCredentials.Run(); err != nil {
//				return err
//			}
//
//			cmdSetContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-context", "default",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
//				"--cluster=kubernetes", "--user=kubelet-bootstrap")
//			cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
//			if err := cmdSetContext.Run(); err != nil {
//				return err
//			}
//
//			cmdUseContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "use-context",
//				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
//				"default")
//			cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
//			if err := cmdUseContext.Run(); err != nil {
//				return err
//			}
//
//			d.Log("Upload bootstrap.kubeconfig to all nodes")
//			if err := UploadBootstrapFile(d.Config, d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
//				return err
//			}
//
//			d.Log("Set Kubelet Clusterrolebinding")
//			if err := SetKubeletClusterrolebinding(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)); err != nil {
//				return nil
//			}
//
//			d.Log("Generate Kubelete Systemd template file")
//			if err := GenerateKubeletSystemdFile(d.Config, kubeSlaveNodes, infra.Version); err != nil {
//				return err
//			}
//
//			d.Log("Upload Kubelete Systemd file")
//			if err := UploadKubeletFile(d.Config, d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
//				return err
//			}
//
//			d.Log("Start Kubelete Service")
//			if err := StartKubelet(d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
//				return err
//			}
//
//			time.Sleep(10 * time.Second)
//			d.Log("Time wait 10 seconds for certificate approve")
//			if err := KubeletCertificateApprove(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)); err != nil {
//				return err
//			}
//		}
//
//		if c.Binary == "kube-proxy" {
//			d.Log("Generate Kube Proxy Systemd template file")
//			if err := GenerateKubeProxyFiles(d.Config, kubeSlaveNodes, infra.Version); err != nil {
//				return err
//			}
//
//			for _, ip := range kubeSlaveNodes {
//				d.Log("Generate kube proxy config file")
//				cmdSetCluster := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
//					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
//					fmt.Sprintf("--certificate-authority=%s", path.Join(d.Config, "ssl", "root", "ca.pem")),
//					"--embed-certs=true",
//					fmt.Sprintf("--server=%s", d.Outputs["KUBE_APISERVER"].(string)))
//				cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
//				if err := cmdSetCluster.Run(); err != nil {
//					return err
//				}
//
//				cmdSetCredentials := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-credentials", "kube-proxy",
//					fmt.Sprintf("--client-certificate=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.pem")),
//					fmt.Sprintf("--client-key=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy-key.pem")),
//					"--embed-certs=true",
//					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
//				)
//				cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
//				if err := cmdSetCredentials.Run(); err != nil {
//					return err
//				}
//
//				cmdSetContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-context", "default",
//					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
//					"--cluster=kubernetes", "--user=kube-proxy")
//				cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
//				if err := cmdSetContext.Run(); err != nil {
//					return err
//				}
//
//				cmdUseContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "use-context",
//					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
//					"default")
//				cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
//				if err := cmdUseContext.Run(); err != nil {
//					return err
//				}
//			}
//
//			d.Log("Upload kube-proxy Systemd file")
//			if err := UploadKubeProxyFiles(d.Config, d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
//				return err
//			}
//
//			d.Log("Start kube-proxy Service")
//			if err := StartKubeProxy(d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
