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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	homeDir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module/service"
	t "github.com/Huawei/containerops/singular/module/template"
)

// JSON export deployment data
func (d *Deployment) JSON() ([]byte, error) {
	return json.Marshal(&d)
}

//
func (d *Deployment) YAML() ([]byte, error) {
	return yaml.Marshal(&d)
}

//
func (d *Deployment) URIs() (namespace, repository, name string, err error) {
	array := strings.Split(d.URI, "/")

	if len(array) != 3 {
		return "", "", "", fmt.Errorf("Invalid deployment URI: %s", d.URI)
	}

	namespace, repository, name = array[0], array[1], array[2]

	return namespace, repository, name, nil
}

// TODO filter the log print with different color.
func (d *Deployment) Log(log string) {
	d.Logs = append(d.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if d.Verbose == true {
		if d.Timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (d *Deployment) Output(key, value string) {
	if d.Outputs == nil {
		d.Outputs = map[string]interface{}{}
	}

	d.Outputs[key] = value
}

// ParseFromFile
func (d *Deployment) ParseFromFile(t string, verbose, timestamp bool) error {
	if data, err := ioutil.ReadFile(t); err != nil {
		d.Log(fmt.Sprintf("Read deployment template file %s error: %s", t, err.Error()))
		return err
	} else {
		if err := yaml.Unmarshal(data, &d); err != nil {
			d.Log(fmt.Sprintf("Unmarshal the template file error: %s", err.Error()))
			return err
		}

		d.Verbose, d.Timestamp = verbose, timestamp

		if err := d.InitConfigPath(""); err != nil {
			return err
		}
	}

	return nil
}

func (d *Deployment) InitConfigPath(path string) error {
	if path == "" {
		home, _ := homeDir.Dir()
		d.Config = fmt.Sprintf("%s/.containerops/singular", home)
	}

	if utils.IsDirExist(d.Config) == false {
		if err := os.MkdirAll(d.Config, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// Check Sequence: CheckServiceAuth -> TODO Check Other?
func (d *Deployment) Check() error {
	if err := d.CheckServiceAuth(); err != nil {
		return fmt.Errorf("check template or configuration error: %s ", err.Error())
	}

	return nil
}

// CheckServiceAuth
func (d *Deployment) CheckServiceAuth() error {
	if d.Service.Provider == "" || d.Service.Token == "" {
		if common.Singular.Provider == "" || common.Singular.Token == "" {
			return fmt.Errorf("Should provide infra service and auth token in %s", "deploy template, or configuration file")
		} else {
			d.Service.Provider, d.Service.Token = common.Singular.Provider, common.Singular.Token
		}
	}

	return nil
}

// Check SSH private and public key files
func (d *Deployment) CheckSSHKey() error {
	if utils.IsFileExist(d.Tools.SSH.Public) == false || utils.IsFileExist(d.Tools.SSH.Private) {
		return fmt.Errorf("Should provide SSH public and private key files in deploy process")
	}

	return nil
}

// Deploy Sequence: Preparing SSH Key files -> Preparing VM -> Preparing SSL root Key files -> Deploy Etcd
//   -> Deploy flannel -> Deploy k8s Master -> Deploy k8s node -> TODO Deploy other...
func (d *Deployment) Deploy() error {
	// Preparing SSH Keys
	if d.Tools.SSH.Public == "" || d.Tools.SSH.Private == "" {
		if public, private, fingerprint, err := CreateSSHKeyFiles(d.Config); err != nil {
			return err
		} else {
			d.Log(fmt.Sprintf(
				"Generate SSH key files successfully, fingerprint is %s", fingerprint))
			d.Tools.SSH.Public, d.Tools.SSH.Private, d.Tools.SSH.Fingerprint = public, private, fingerprint
		}
	}

	switch d.Service.Provider {
	case "digitalocean":
		do := new(service.DigitalOcean)
		do.Token = d.Service.Token
		do.Region, do.Size, do.Image = d.Service.Region, d.Service.Size, d.Service.Image

		// Init DigitalOcean API client.
		do.InitClient()
		d.Log("Init DigitalOcean API client.")

		// Upload ssh public key
		if err := do.UpdateSSHKey(d.Tools.SSH.Public); err != nil {
			return err
		}
		d.Log("Uploaded SSH public key to the DigitalOcean.")

		d.Log("Creating Droplet in DigitalOcean.")
		// Create DigitalOcean Droplets
		if err := do.CreateDroplet(d.Nodes, d.Tools.SSH.Fingerprint); err != nil {
			return err
		}

		d.Log("Created Droplets successfully: ")
		i := 0
		for key, _ := range do.Droplets {
			d.Log(fmt.Sprintf("Node %d ip: %s", i, key))
			d.Output(fmt.Sprintf("NODE_%d", i), key)
			i += 1
		}

		fmt.Println(Red("Waiting 60 seconds for preparing droplets..."))
		time.Sleep(60 * time.Second)

		// Generate CA Root files
		if roots, err := GenerateCARootFiles(d.Config); err != nil {
			return err
		} else {
			d.Log("CA Root files generated successfully.")

			for key, value := range roots {
				d.Output(key, value)
			}

			for ip, _ := range do.Droplets {
				d.Log(fmt.Sprintf("Upload SSL Root files to Droplet[%s] and init environments.", ip))
				if err := UploadCARootFiles(d.Config, roots, ip); err != nil {
					return err
				}
			}

		}

		for _, infra := range d.Infras {
			switch infra.Name {
			case "etcd":
				if err := d.DeployEtcd(infra); err != nil {
					return err
				}
			case "flannel":
				if err := d.DeployFlannel(infra); err != nil {
					return err
				}
			case "docker":
				if err := d.DeployDocker(infra); err != nil {
					return err
				}
			case "kubernetes":
				if err := d.DeployKubernetes(infra); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unsupport infrastruction software: %s", infra)
			}

		}

	default:
		return fmt.Errorf("Unsupport service provide: %s", d.Service.Provider)

	}

	return nil
}

// DeployEtcd is function deployment etcd cluster.
// Notes:
//   1. Only count master nodes in etcd deploy process.
//   2.
func (d *Deployment) DeployEtcd(infra Infra) error {
	d.Log("Deploying etcd clusters.")

	if infra.Nodes.Master > d.Nodes {
		return fmt.Errorf("Deploy %s nodes more than %d", infra.Name, d.Nodes)
	}

	etcdNodes := map[string]string{}
	etcdEndpoints, etcdAdminEndpoints := []string{}, []string{}

	for i := 0; i < infra.Nodes.Master; i++ {
		etcdNodes[fmt.Sprintf("etcd-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
		etcdEndpoints = append(etcdEndpoints,
			fmt.Sprintf("https://%s:2379", d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)))
		etcdAdminEndpoints = append(etcdAdminEndpoints,
			fmt.Sprintf("%s=https://%s:2380", fmt.Sprintf("etcd-node-%d", i),
				d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)))
	}

	d.Output("EtcdEndpoints", strings.Join(etcdEndpoints, ","))
	d.Log(fmt.Sprintf("Generating Etcd endpoints environment variable [%s], value is\n [%s]", "EtcdEndpoints", strings.Join(etcdEndpoints, ",")))

	d.Log(fmt.Sprintf("Generating SSL files and systemd service file for Etcd cluster."))
	if err := GenerateEtcdFiles(d.Config, etcdNodes, strings.Join(etcdAdminEndpoints, ","), infra.Version); err != nil {
		return err
	} else {
		d.Log(fmt.Sprintf("Uploading SSL files to nodes of Etcd Cluster."))
		if err := UploadEtcdCAFiles(d.Config, etcdNodes); err != nil {
			return err
		}

		d.Log(fmt.Sprintf("Downloading Etcd binary files to nodes of Etcd Cluster."))
		for _, c := range infra.Components {
			if err := d.DownloadBinaryFile(c.Binary, c.URL, etcdNodes); err != nil {
				return err
			}
		}

		d.Log(fmt.Sprintf("Staring Etcd Cluster."))
		if err := StartEtcdCluster(d.Tools.SSH.Private, etcdNodes); err != nil {
			return err
		}

	}

	return nil
}

func (d *Deployment) DownloadBinaryFile(file, url string, nodes map[string]string) error {
	for _, ip := range nodes {
		downloadCmd := fmt.Sprintf("curl %s -o /usr/local/bin/%s", url, file)
		chmodCmd := fmt.Sprintf("chmod +x /usr/local/bin/%s", file)

		d.Log(fmt.Sprintf("Downloading %s to Node[%s]", file, ip))
		if err := utils.SSHCommand("root", d.Tools.SSH.Private, ip, 22, downloadCmd, os.Stdout, os.Stderr); err != nil {
			return err
		}

		d.Log(fmt.Sprintf("Change %s mode in Node[%s]", file, ip))
		if err := utils.SSHCommand("root", d.Tools.SSH.Private, ip, 22, chmodCmd, os.Stdout, os.Stderr); err != nil {
			return err
		}

	}

	return nil
}

func (d *Deployment) DeployFlannel(infra Infra) error {
	flanneldNodes := map[string]string{}
	for i := 0; i < infra.Nodes.Master; i++ {
		flanneldNodes[fmt.Sprintf("flanneld-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
	}

	d.Log(fmt.Sprintf("Generating SSL files and systemd service file for Flanneld."))
	if err := GenerateFlanneldFiles(d.Config, flanneldNodes, d.Outputs["EtcdEndpoints"].(string), infra.Version); err != nil {
		return err
	} else {
		d.Log(fmt.Sprintf("Uploading SSL files and systemd service to nodes of Flanneld."))
		if err := UploadFlanneldCAFiles(d.Config, flanneldNodes); err != nil {
			return err
		}

		for i, c := range infra.Components {
			d.Log(fmt.Sprintf("Downloading Flanneld binary files to Nodes."))
			if err := d.DownloadBinaryFile(c.Binary, c.URL, flanneldNodes); err != nil {
				return err
			}

			if c.Before != "" && i == 0 {
				d.Log(fmt.Sprintf("Execute Flanneld before scripts: %s", c.Before))
				if err := BeforeFlanneldExecute(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", i)].(string), c.Before, d.Outputs["EtcdEndpoints"].(string)); err != nil {
					return err
				}
			}
		}

		d.Log(fmt.Sprintf("Staring Flanneld Service."))
		if err := StartFlanneldCluster(d.Tools.SSH.Private, flanneldNodes); err != nil {
			return err
		}
	}

	return nil
}

func (d *Deployment) DeployDocker(infra Infra) error {
	dockerNodes := map[string]string{}
	for i := 0; i < infra.Nodes.Master; i++ {
		dockerNodes[fmt.Sprintf("docker-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
	}

	d.Log(fmt.Sprintf("Generating SSL files and systemd service file for Docker."))
	if err := GenerateDockerFiles(d.Config, dockerNodes, infra.Version); err != nil {
		return err
	} else {
		d.Log(fmt.Sprintf("Uploading SSL files and systemd service to nodes of Docker."))
		if err := UploadDockerCAFiles(d.Config, dockerNodes); err != nil {
			return err
		}

		for _, c := range infra.Components {
			d.Log(fmt.Sprintf("Downloading Docker binary files to Nodes."))
			if err := d.DownloadBinaryFile(c.Binary, c.URL, dockerNodes); err != nil {
				return err
			}

			if c.Before != "" {
				for _, ip := range dockerNodes {
					d.Log(fmt.Sprintf("Execute Docker before scripts: %s in %s", ip, c.Before))
					if err := BeforeDockerExecute(d.Tools.SSH.Private, ip, c.Before); err != nil {
						return err
					}
				}
			}
		}

		for _, ip := range dockerNodes {
			d.Log(fmt.Sprintf("Start Docker in %s", ip))
			if err := StartDockerDaemon(d.Tools.SSH.Private, ip); err != nil {
				return err
			}
		}

		for _, c := range infra.Components {
			if c.After != "" {
				for _, ip := range dockerNodes {
					d.Log(fmt.Sprintf("Execute Docker After scripts: %s in %s", c.After, ip))
					if err := AfterDockerExecute(d.Tools.SSH.Private, ip, c.After); err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}

// DeployKubernetes is function deployment Kubernetes cluster include master and nodes.
// Notes:
//   1. Kubernetes master cluster IP.
//   2. Set kubectl config files.
//   3. Deploy Kubernetes master.
//   4. Deploy Kubernetes nodes.
func (d *Deployment) DeployKubernetes(infra Infra) error {
	// TODO Now singular only support one master and multiple nodes architect.
	// TODO So we decide the Kubernetes master IP is NODE_0 .
	masterIp := d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)
	etcdEndpoints := d.Outputs["EtcdEndpoints"].(string)

	d.Output("MASTER_IP", masterIp)
	d.Output("KUBE_APISERVER", fmt.Sprintf("https://%s:6443", masterIp))

	kubeMasterNodes := map[string]string{}
	for i := 0; i < infra.Nodes.Master; i++ {
		kubeMasterNodes[fmt.Sprintf("kube-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
	}

	kubeSlaveNodes := map[string]string{}
	for i := 0; i < infra.Nodes.Node; i++ {
		kubeSlaveNodes[fmt.Sprintf("kube-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
	}

	for _, c := range infra.Components {
		d.Log(fmt.Sprintf("Download %s binary files", c.Binary))
		if err := d.DownloadBinaryFile(c.Binary, c.URL, kubeSlaveNodes); err != nil {
			return err
		}
	}

	for _, c := range infra.Components {
		if c.Binary == "kubectl" {
			if utils.IsDirExist(path.Join(d.Config, "kubectl")) == true {
				os.RemoveAll(path.Join(d.Config, "kubectl"))
			}

			os.MkdirAll(path.Join(d.Config, "kubectl"), os.ModePerm)

			d.Log("Downloading kubectl binary file")
			cmdDownload := exec.Command("curl", c.URL, "-o", fmt.Sprintf("%s/kubectl/kubectl", d.Config))
			cmdDownload.Stdout, cmdDownload.Stderr = os.Stdout, os.Stderr
			if err := cmdDownload.Run(); err != nil {
				return err
			}

			cmdChmod := exec.Command("chmod", "+x", fmt.Sprintf("%s/kubectl/kubectl", d.Config))
			cmdChmod.Stdout, cmdChmod.Stderr = os.Stdout, os.Stderr
			if err := cmdChmod.Run(); err != nil {
				return err
			}

			d.Log("Genearate kubernetes admin ca files")
			if err := GenerateAdminCAFiles(d.Config); err != nil {
				return err
			}

			d.Log("Generate kubectl config file")
			cmdSetCluster := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
				fmt.Sprintf("--certificate-authority=%s", path.Join(d.Config, "ssl", "root", "ca.pem")),
				"--embed-certs=true",
				fmt.Sprintf("--server=%s", d.Outputs["KUBE_APISERVER"].(string)))
			cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
			if err := cmdSetCluster.Run(); err != nil {
				return err
			}

			cmdSetCredentials := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-credentials", "admin",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
				fmt.Sprintf("--client-certificate=%s", path.Join(d.Config, "kubectl", "admin.pem")),
				"--embed-certs=true",
				fmt.Sprintf("--client-key=%s", path.Join(d.Config, "kubectl", "admin-key.pem")))
			cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
			if err := cmdSetCredentials.Run(); err != nil {
				return err
			}

			cmdSetContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-context", "kubernetes",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
				"--cluster=kubernetes", "--user=admin")
			cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
			if err := cmdSetContext.Run(); err != nil {
				return err
			}

			cmdUseContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "use-context",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "config")),
				"kubernetes")
			cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
			if err := cmdUseContext.Run(); err != nil {
				return err
			}

			d.Log("Upload kubectl config file to Kubernetes nodes")
			if err := UploadKubeConfigFiles(d.Config, kubeSlaveNodes); err != nil {
				return err
			}
		}

		if c.Binary == "kube-apiserver" {
			d.Log("Generate Kuberentes Token API file")
			if err := GenerateTokenFile(d.Config); err != nil {
				return err
			}

			d.Log("Generate Kubernetes SSL files and systemd service file")
			if err := GenerateKuberAPIServerCAFiles(d.Config, masterIp, etcdEndpoints, infra.Version); err != nil {
				return err
			}

			d.Log("Upload Kubernetes Token file")
			if err := UploadTokenFiles(d.Config, masterIp); err != nil {
				return err
			}

			d.Log("Upload Kubernetes API Server SSL files and systemd service file")
			if err := UploadKubeAPIServerCAFiles(d.Config, masterIp); err != nil {
				return err
			}

			d.Log("Start Kubernetes API Server")
			if err := StartKubeAPIServer(d.Config, masterIp); err != nil {
				return err
			}

		}

		if c.Binary == "kube-controller-manager" {
			d.Log("Generate Kube-controller-manager systemd service file")
			if err := GenerateKuberControllerManagerFiles(d.Config, masterIp, etcdEndpoints, infra.Version); err != nil {
				return err
			}

			d.Log("Upload Kuber-controller-manager systemd service file")
			if err := UploadKuberControllerFiles(d.Config, masterIp); err != nil {
				return err
			}

			d.Log("Start Kube-controller-manager")
			if err := StartKuberController(d.Config, masterIp); err != nil {
				return err
			}
		}

		if c.Binary == "kube-scheduler" {
			d.Log("Generate Kube-scheduler systemd service file")
			if err := GenerateKuberSchedulerManagerFiles(d.Config, masterIp, etcdEndpoints, infra.Version); err != nil {
				return err
			}

			d.Log("Upload Kuber-scheduler systemd service file")
			if err := UploadKuberSchedulerManagerFiles(d.Config, masterIp); err != nil {
				return err
			}

			d.Log("Start Kube-scheduler")
			if err := StartKuberSchedulerManager(d.Config, masterIp); err != nil {
				return err
			}
		}

		if c.Binary == "kubelet" {
			d.Log("Generate bootstrap.kubeconfig config file")
			cmdSetCluster := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
				fmt.Sprintf("--certificate-authority=%s", path.Join(d.Config, "ssl", "root", "ca.pem")),
				"--embed-certs=true",
				fmt.Sprintf("--server=%s", d.Outputs["KUBE_APISERVER"].(string)))
			cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
			if err := cmdSetCluster.Run(); err != nil {
				return err
			}

			cmdSetCredentials := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-credentials", "kubelet-bootstrap",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
				fmt.Sprintf("--token=%s", t.BooststrapToken))
			cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
			if err := cmdSetCredentials.Run(); err != nil {
				return err
			}

			cmdSetContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-context", "default",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
				"--cluster=kubernetes", "--user=kubelet-bootstrap")
			cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
			if err := cmdSetContext.Run(); err != nil {
				return err
			}

			cmdUseContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "use-context",
				fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "kubectl", "bootstrap.kubeconfig")),
				"default")
			cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
			if err := cmdUseContext.Run(); err != nil {
				return err
			}

			d.Log("Upload bootstrap.kubeconfig to all nodes")
			if err := UploadBootstrapFile(d.Config, kubeSlaveNodes); err != nil {
				return err
			}

			d.Log("Set Kubelet Clusterrolebinding")
			if err := SetKubeletClusterrolebinding(d.Config, d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)); err != nil {
				return nil
			}

			d.Log("Generate Kubelete Systemd template file")
			if err := GenerateKubeletSystemdFile(d.Config, kubeSlaveNodes, infra.Version); err != nil {
				return err
			}

			d.Log("Upload Kubelete Systemd file")
			if err := UploadKubeletFile(d.Config, kubeSlaveNodes); err != nil {
				return err
			}

			d.Log("Start Kubelete Service")
			if err := StartKubelet(d.Config, kubeSlaveNodes); err != nil {
				return err
			}

			time.Sleep(10 * time.Second)
			d.Log("Time wait 10 seconds for certificate approve")
			if err := KubeletCertificateApprove(d.Config, d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)); err != nil {
				return err
			}
		}

		if c.Binary == "kube-proxy" {
			d.Log("Generate Kube Proxy Systemd template file")
			if err := GenerateKubeProxyFiles(d.Config, kubeSlaveNodes, infra.Version); err != nil {
				return err
			}

			for _, ip := range kubeSlaveNodes {
				d.Log("Generate kube proxy config file")
				cmdSetCluster := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
					fmt.Sprintf("--certificate-authority=%s", path.Join(d.Config, "ssl", "root", "ca.pem")),
					"--embed-certs=true",
					fmt.Sprintf("--server=%s", d.Outputs["KUBE_APISERVER"].(string)))
				cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
				if err := cmdSetCluster.Run(); err != nil {
					return err
				}

				cmdSetCredentials := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-credentials", "kube-proxy",
					fmt.Sprintf("--client-certificate=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.pem")),
					fmt.Sprintf("--client-key=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy-key.pem")),
					"--embed-certs=true",
					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
				)
				cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
				if err := cmdSetCredentials.Run(); err != nil {
					return err
				}

				cmdSetContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "set-context", "default",
					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
					"--cluster=kubernetes", "--user=kube-proxy")
				cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
				if err := cmdSetContext.Run(); err != nil {
					return err
				}

				cmdUseContext := exec.Command(path.Join(d.Config, "kubectl", "kubectl"), "config", "use-context",
					fmt.Sprintf("--kubeconfig=%s", path.Join(d.Config, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
					"default")
				cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
				if err := cmdUseContext.Run(); err != nil {
					return err
				}
			}

			d.Log("Upload kube-proxy Systemd file")
			if err := UploadKubeProxyFiles(d.Config, kubeSlaveNodes); err != nil {
				return err
			}

			d.Log("Start kube-proxy Service")
			if err := StartKubeProxy(d.Config, kubeSlaveNodes); err != nil {
				return err
			}
		}
	}

	return nil
}
