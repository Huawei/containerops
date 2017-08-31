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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"

	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module/objects"
	t "github.com/Huawei/containerops/singular/module/template"
	"github.com/Huawei/containerops/singular/module/tools"
)

const (
	KubeServerConfig = "/etc/kubernetes"
	KubeServerSSL    = "ssl"
)

// DeployKubernetes is function deployment Kubernetes cluster include master and nodes.
// Notes:
//   1. Kubernetes master cluster IP.
//   2. Set kubectl config files.
//   3. Deploy Kubernetes master.
//   4. Deploy Kubernetes nodes.
func DeployKubernetesInCluster(d *objects.Deployment, infra *objects.Infra, stdout io.Writer, timestamp bool) error {
	// TODO Now singular only support one master and multiple nodes architect.
	// TODO So we decide the Kubernetes master IP is NODE_0 .
	masterIP := d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)
	etcdEndpoints := d.Outputs["EtcdEndpoints"].(string)

	// Master IP
	d.Output("MASTER_IP", masterIP)
	d.Output("KUBE_APISERVER", fmt.Sprintf("https://%s:6443", masterIP))

	// Master Nodes
	kubeMasterNodes := []objects.Node{}
	for i := 0; i < infra.Master; i++ {
		kubeMasterNodes = append(kubeMasterNodes, d.Nodes[i])
	}

	// Minion Nodes
	kubeSlaveNodes := []objects.Node{}
	for i := 0; i < infra.Minion; i++ {
		kubeSlaveNodes = append(kubeSlaveNodes, d.Nodes[i])
	}

	// Download binary in master nodes
	for _, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, kubeMasterNodes, stdout, timestamp); err != nil {
			return err
		}
	}

	// Download binary in slave nodes
	for _, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, kubeSlaveNodes, stdout, timestamp); err != nil {
			return err
		}
	}

	for _, c := range infra.Components {
		switch c.Binary {
		case "kubectl":
			//Generate kubectl config file and CA SSL files.
			if err := setKubectlFiles(d.Config, c.URL, masterIP); err != nil {
				return err
			}

			//Upload kubectl config file to Kubernetes nodes
			if err := uploadKubeConfigFiles(d.Config, d.Tools.SSH.Private, kubeSlaveNodes, stdout); err != nil {
				return err
			}
		case "kube-apiserver":
			//Generate Kubernetes Token API file
			if err := generateTokenFile(d.Config); err != nil {
				return err
			}

			//Upload Kubernetes Token file
			if err := uploadTokenFiles(d.Config, d.Tools.SSH.Private, masterIP, stdout); err != nil {
				return err
			}

			//Generate Kubernetes SSL files and systemd service file
			if err := generateKubeAPIServerCAFiles(d.Config, masterIP, etcdEndpoints, infra.Version); err != nil {
				return err
			}

			//Upload Kubernetes API Server SSL files and systemd service file
			if err := uploadKubeAPIServerCAFiles(d.Config, d.Tools.SSH.Private, masterIP, stdout); err != nil {
				return err
			}

			//Start Kubernetes API Server
			if err := startKubeAPIServer(d.Tools.SSH.Private, masterIP); err != nil {
				return err
			}
		case "kube-controller-manager":
			//Generate Kube-controller-manager systemd service file
			if err := generateKubeControllerManagerFiles(d.Config, masterIP, etcdEndpoints, infra.Version); err != nil {
				return err
			}

			//Upload Kuber-controller-manager systemd service file
			if err := uploadKubeControllerFiles(d.Config, d.Tools.SSH.Private, masterIP, stdout); err != nil {
				return err
			}

			//Start Kube-controller-manager
			if err := startKubeController(d.Tools.SSH.Private, masterIP); err != nil {
				return err
			}
		case "kube-scheduler":
			//Generate Kube-scheduler systemd service file
			if err := generateKubeSchedulerFiles(d.Config, masterIP, etcdEndpoints, infra.Version); err != nil {
				return err
			}

			//Upload Kuber-scheduler systemd service file
			if err := uploadKubeSchedulerFiles(d.Config, d.Tools.SSH.Private, masterIP, stdout); err != nil {
				return err
			}

			//Start Kube-scheduler
			if err := startKubeScheduler(d.Tools.SSH.Private, masterIP); err != nil {
				return err
			}
		case "kubelet":
			if err := generateBootstrapFile(d.Config, masterIP); err != nil {
				return err
			}

			if err := uploadBootstrapFile(d.Config, d.Tools.SSH.Private, kubeSlaveNodes, stdout); err != nil {
				return err
			}

			if err := setKubeletClusterrolebinding(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)); err != nil {
				return nil
			}

			if err := generateKubeletSystemdFile(d.Config, kubeSlaveNodes, infra.Version); err != nil {
				return err
			}

			if err := uploadKubeletFile(d.Config, d.Tools.SSH.Private, kubeSlaveNodes, stdout); err != nil {
				return err
			}

			if err := startKubelet(d.Tools.SSH.Private, kubeSlaveNodes, stdout); err != nil {
				return err
			}

			time.Sleep(10 * time.Second)

			if err := kubeletCertificateApprove(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)); err != nil {
				return err
			}
		case "kube-proxy":
			//Generate Kube Proxy Systemd template file
			if err := generateKubeProxyFiles(d.Config, kubeSlaveNodes, infra.Version); err != nil {
				return err
			}

			for _, node := range kubeSlaveNodes {
				if err := generateKubeProxyConfigFile(d.Config, node.IP, masterIP); err != nil {
					return err
				}
			}

			//Upload kube-proxy Systemd file
			if err := uploadKubeProxyFiles(d.Config, d.Tools.SSH.Private, kubeSlaveNodes, stdout); err != nil {
				return err
			}

			//Start kube-proxy Service
			if err := startKubeProxy(d.Tools.SSH.Private, kubeSlaveNodes); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupport kubenetes component: %s", c.Binary)
		}
	}

	return nil
}

// setKubectlConfig download kubectl and set config file.
func setKubectlFiles(src, link, master string) error {
	//Make kubectl folder
	if utils.IsDirExist(path.Join(src, tools.KubectlFileFolder)) == true {
		os.RemoveAll(path.Join(src, tools.KubectlFileFolder))
	}

	os.MkdirAll(path.Join(src, tools.KubectlFileFolder), os.ModePerm)

	//Download or copy kubectl file
	if a, err := url.Parse(link); err != nil {
		return err
	} else {
		if a.Scheme == "" {
			cmdCopy := exec.Command("cp", link, path.Join(src, tools.KubectlFileFolder, tools.KubectlFile))

			cmdCopy.Stdout, cmdCopy.Stderr = os.Stdout, os.Stderr
			if err := cmdCopy.Run(); err != nil {
				return err
			}
		} else {
			cmdDownload := exec.Command("curl", link, "-o", path.Join(src, tools.KubectlFileFolder, tools.KubectlFile), src)

			cmdDownload.Stdout, cmdDownload.Stderr = os.Stdout, os.Stderr
			if err := cmdDownload.Run(); err != nil {
				return err
			}
		}
	}

	//Change mode +x for kubectl
	cmdChmod := exec.Command("chmod", "+x", path.Join(src, tools.KubectlFileFolder, tools.KubectlFile))
	cmdChmod.Stdout, cmdChmod.Stderr = os.Stdout, os.Stderr
	if err := cmdChmod.Run(); err != nil {
		return err
	}

	//Generate Kubernetes admin CA files
	if err := generateKubeAdminCAFiles(src); err != nil {
		return err
	}

	//Generate kubectl config file
	if err := setKubeConfigFile(src, master); err != nil {
		return err
	}

	return nil
}

// Generate Kubernetes Admin CA files
func generateKubeAdminCAFiles(src string) error {
	base := path.Join(src, tools.KubectlFileFolder)

	caFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	var tpl bytes.Buffer
	var err error

	sslTp := template.New("admin-csr")
	sslTp, _ = sslTp.Parse(t.AdminCATemplate)
	sslTp.Execute(&tpl, nil)
	csrFileBytes := tpl.Bytes()

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}

	err = json.Unmarshal(csrFileBytes, &req)
	if err != nil {
		return err
	}

	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		return err
	}

	c := cli.Config{
		CAFile:     caFile,
		CAKeyFile:  caKeyFile,
		ConfigFile: configFile,
		Profile:    "kubernetes",
		Hostname:   "",
	}

	s, err := sign.SignerFromConfig(c)
	if err != nil {
		return err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(c.Hostname),
		Profile: c.Profile,
	}

	cert, err = s.Sign(signReq)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAdminCSRConfigFile), csrFileBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAdminKeyPemFile), key, 0600)
	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAdminCSRFile), csrBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAdminPemFile), cert, 0600)

	if err != nil {
		return err
	}

	return nil
}

// Generate kubectl config file
func setKubeConfigFile(src, master string) error {
	base := path.Join(src, tools.KubectlFileFolder)
	adminPem := path.Join(base, tools.CAKubeAdminPemFile)
	adminKey := path.Join(base, tools.CAKubeAdminKeyPemFile)

	kubectl := path.Join(src, tools.KubectlFileFolder, tools.KubectlFile)
	config := path.Join(src, tools.KubectlFileFolder, tools.KubectlConfigFile)
	caFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)

	//Generate kubectl config file
	cmdSetCluster := exec.Command(kubectl, "config", "set-cluster", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", config),
		fmt.Sprintf("--certificate-authority=%s", caFile),
		"--embed-certs=true",
		fmt.Sprintf("--server=%s", master))
	cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
	if err := cmdSetCluster.Run(); err != nil {
		return err
	}

	cmdSetCredentials := exec.Command(kubectl, "config", "set-credentials", "admin",
		fmt.Sprintf("--kubeconfig=%s", config),
		fmt.Sprintf("--client-certificate=%s", adminPem),
		"--embed-certs=true",
		fmt.Sprintf("--client-key=%s", adminKey))
	cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
	if err := cmdSetCredentials.Run(); err != nil {
		return err
	}

	cmdSetContext := exec.Command(kubectl, "config", "set-context", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", config),
		"--cluster=kubernetes", "--user=admin")
	cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
	if err := cmdSetContext.Run(); err != nil {
		return err
	}

	cmdUseContext := exec.Command(kubectl, "config", "use-context",
		fmt.Sprintf("--kubeconfig=%s", config),
		"kubernetes")
	cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
	if err := cmdUseContext.Run(); err != nil {
		return err
	}

	return nil
}

// Upload Kubectl config file and ca ssl files.
func uploadKubeConfigFiles(src, key string, nodes []objects.Node, stdout io.Writer) error {
	base := path.Join(src, tools.KubectlFileFolder)
	config := path.Join(src, tools.KubectlFileFolder, tools.KubectlConfigFile)

	for _, node := range nodes {
		var err error

		err = utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, "mkdir -p /root/.kube", stdout, os.Stderr)

		_, err = tools.DownloadComponent(config, "/root/.kube/config", node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAdminCSRConfigFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAdminCSRConfigFile), node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAdminKeyPemFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAdminKeyPemFile), node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAdminCSRFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAdminCSRFile), node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAdminPemFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAdminPemFile), node.IP, key, node.User, stdout)

		if err != nil {
			return err
		}
	}

	return nil
}

// Generate Kubectl Token file
func generateTokenFile(src string) error {
	if utils.IsDirExist(path.Join(src, "kubectl")) == false {
		os.MkdirAll(path.Join(src, "kubectl"), os.ModePerm)
	}

	var tokenTpl bytes.Buffer

	tokenTp := template.New("token")
	tokenTp, _ = tokenTp.Parse("{{.Token}},kubelet-bootstrap,10001,\"system:kubelet-bootstrap\"")
	tokenTp.Execute(&tokenTpl, map[string]string{"Token": t.BooststrapToken})
	tokenTpFileBytes := tokenTpl.Bytes()

	if err := ioutil.WriteFile(path.Join(src, tools.KubectlFileFolder, tools.KubeTokenCSVFile), tokenTpFileBytes, 0700); err != nil {
		return err
	}

	return nil
}

// Upload Token CSV file
func uploadTokenFiles(src, key, ip string, stdout io.Writer) error {
	file := path.Join(src, tools.KubectlFileFolder, tools.KubeTokenCSVFile)

	if _, err := tools.DownloadComponent(file, path.Join(KubeServerConfig, tools.KubeTokenCSVFile), ip, key, tools.DefaultSSHUser, stdout); err != nil {
		return err
	}

	return nil
}

// Kubernetes master
type KubeMaster struct {
	MasterIP string
	Nodes    string
}

// Generate Kube API Server CA SSL files.
func generateKubeAPIServerCAFiles(src string, masterIP, etcdEndpoints string, version string) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, os.ModePerm)

	caFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	master := KubeMaster{
		MasterIP: masterIP,
		Nodes:    etcdEndpoints,
	}

	var tpl bytes.Buffer
	var err error

	sslTp := template.New("kube-csr")
	sslTp, _ = sslTp.Parse(t.KubernetesCATemplate[version])
	sslTp.Execute(&tpl, master)
	csrFileBytes := tpl.Bytes()

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}

	err = json.Unmarshal(csrFileBytes, &req)
	if err != nil {
		return err
	}

	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		return err
	}

	c := cli.Config{
		CAFile:     caFile,
		CAKeyFile:  caKeyFile,
		ConfigFile: configFile,
		Profile:    "kubernetes",
		Hostname:   "",
	}

	s, err := sign.SignerFromConfig(c)
	if err != nil {
		return err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(c.Hostname),
		Profile: c.Profile,
	}

	cert, err = s.Sign(signReq)
	if err != nil {
		return err
	}

	var serviceTpl bytes.Buffer

	serviceTp := template.New("kube-systemd")
	serviceTp, _ = serviceTp.Parse(t.KubernetesAPIServerSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, master)
	serviceTpFileBytes := serviceTpl.Bytes()

	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAPIServerCSRConfigFile), csrFileBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAPIServerKeyPemFile), key, 0600)
	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAPIServerCSRFile), csrBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, tools.CAKubeAPIServerPemFile), cert, 0600)

	err = ioutil.WriteFile(path.Join(base, tools.KubeAPIServerSystemdFile), serviceTpFileBytes, 0700)

	if err != nil {
		return err
	}

	return nil
}

// Upload Kube API Server systemd file and CA SSL file.
func uploadKubeAPIServerCAFiles(src, key, ip string, stdout io.Writer) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	var err error

	_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAPIServerCSRConfigFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAPIServerCSRConfigFile), ip, key, tools.DefaultSSHUser, stdout)
	_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAPIServerKeyPemFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAPIServerKeyPemFile), ip, key, tools.DefaultSSHUser, stdout)
	_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAPIServerCSRFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAPIServerCSRFile), ip, key, tools.DefaultSSHUser, stdout)
	_, err = tools.DownloadComponent(path.Join(base, tools.CAKubeAPIServerPemFile), path.Join(KubeServerConfig, KubeServerSSL, tools.CAKubeAPIServerPemFile), ip, key, tools.DefaultSSHUser, stdout)
	_, err = tools.DownloadComponent(path.Join(base, tools.KubeAPIServerSystemdFile), path.Join(tools.SytemdServerPath, tools.KubeAPIServerSystemdFile), ip, key, tools.DefaultSSHUser, stdout)

	if err != nil {
		return err
	}

	return nil
}

func startKubeAPIServer(key, ip string) error {
	cmd := "systemctl daemon-reload && systemctl enable kube-apiserver && systemctl start --no-block kube-apiserver"

	if err := utils.SSHCommand(tools.DefaultSSHUser, key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func generateKubeControllerManagerFiles(src string, masterIP, etcdEndpoints string, version string) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	master := KubeMaster{
		MasterIP: masterIP,
		Nodes:    etcdEndpoints,
	}

	var serviceTpl bytes.Buffer

	serviceTp := template.New("kube-control-systemd")
	serviceTp, _ = serviceTp.Parse(t.KubernetesControllerManagerSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, master)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(path.Join(base, "kube-controller-manager.service"), serviceTpFileBytes, 0700); err != nil {
		return err
	}

	return nil
}

func uploadKubeControllerFiles(src, key, ip string, stdout io.Writer) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	if _, err := tools.DownloadComponent(path.Join(base, "kube-controller-manager.service"), "/etc/systemd/system/kube-controller-manager.service", ip, key, tools.DefaultSSHUser, stdout); err != nil {
		return err
	}

	return nil
}

func startKubeController(key, ip string) error {
	cmd := "systemctl daemon-reload && systemctl enable kube-controller-manager && systemctl start --no-block kube-controller-manager"

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func generateKubeSchedulerFiles(src string, masterIP, etcdEndpoints string, version string) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	master := KubeMaster{
		MasterIP: masterIP,
		Nodes:    etcdEndpoints,
	}

	var serviceTpl bytes.Buffer

	serviceTp := template.New("kube-scheduler-systemd")
	serviceTp, _ = serviceTp.Parse(t.KubernetesSchedulerSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, master)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(path.Join(base, "kube-scheduler.service"), serviceTpFileBytes, 0700); err != nil {
		return err
	}

	return nil
}

func uploadKubeSchedulerFiles(src, key, ip string, stdout io.Writer) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	if _, err := tools.DownloadComponent(path.Join(base, "kube-scheduler.service"), "/etc/systemd/system/kube-scheduler.service", ip, key, tools.DefaultSSHUser, stdout); err != nil {
		return err
	}

	return nil
}

func startKubeScheduler(key, ip string) error {
	cmd := "systemctl daemon-reload && systemctl enable kube-scheduler && systemctl start --no-block kube-scheduler"

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func generateBootstrapFile(src, master string) error {
	//Generate bootstrap.kubeconfig config file
	cmdSetCluster := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "kubectl", "bootstrap.kubeconfig")),
		fmt.Sprintf("--certificate-authority=%s", path.Join(src, "ssl", "root", "ca.pem")),
		"--embed-certs=true",
		fmt.Sprintf("--server=%s", master))
	cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
	if err := cmdSetCluster.Run(); err != nil {
		return err
	}

	cmdSetCredentials := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "set-credentials", "kubelet-bootstrap",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "kubectl", "bootstrap.kubeconfig")),
		fmt.Sprintf("--token=%s", t.BooststrapToken))
	cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
	if err := cmdSetCredentials.Run(); err != nil {
		return err
	}

	cmdSetContext := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "set-context", "default",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "kubectl", "bootstrap.kubeconfig")),
		"--cluster=kubernetes", "--user=kubelet-bootstrap")
	cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
	if err := cmdSetContext.Run(); err != nil {
		return err
	}

	cmdUseContext := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "use-context",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "kubectl", "bootstrap.kubeconfig")),
		"default")
	cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
	if err := cmdUseContext.Run(); err != nil {
		return err
	}

	return nil
}

func uploadBootstrapFile(src, key string, nodes []objects.Node, stdout io.Writer) error {
	config := path.Join(src, "kubectl", "bootstrap.kubeconfig")

	for _, node := range nodes {
		if _, err := tools.DownloadComponent(config, "/etc/kubernetes/bootstrap.kubeconfig", node.IP, key, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

func generateKubeletSystemdFile(src string, nodes []objects.Node, version string) error {
	for _, node := range nodes {
		kubeNode := map[string]string{
			"IP": node.IP,
		}

		base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder, node.IP)
		if utils.IsDirExist(base) == true {
			os.RemoveAll(base)
		}

		os.MkdirAll(base, os.ModePerm)

		var serviceTpl bytes.Buffer

		serviceTp := template.New("kubelete-systemd")
		serviceTp, _ = serviceTp.Parse(t.KubeletSystemdTemplate[version])
		serviceTp.Execute(&serviceTpl, kubeNode)
		serviceTpFileBytes := serviceTpl.Bytes()

		if err := ioutil.WriteFile(path.Join(base, "kubelet.service"), serviceTpFileBytes, 0700); err != nil {
			return err
		}
	}
	return nil
}

func uploadKubeletFile(src, key string, nodes []objects.Node, stdout io.Writer) error {
	for _, node := range nodes {
		file := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder, node.IP, "kubelet.service")

		if err := utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, "mkdir -p /var/lib/kubelet", stdout, os.Stderr); err != nil {
			return err
		}

		if _, err := tools.DownloadComponent(file, "/etc/systemd/system/kubelet.service", node.IP, key, node.User, stdout); err != nil {
			return err
		}

	}

	return nil
}

func setKubeletClusterrolebinding(key, ip string) error {
	if err := utils.SSHCommand("root", key, ip, 22, "kubectl create clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func kubeletCertificateApprove(key, ip string) error {
	if err := utils.SSHCommand("root", key, ip, 22, "kubectl certificate approve `kubectl get csr -o name`", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func startKubelet(key string, nodes []objects.Node, stdout io.Writer) error {
	for _, node := range nodes {
		cmd := "systemctl daemon-reload && systemctl enable kubelet && systemctl start --no-block kubelet"

		if err := utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, cmd, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

func generateKubeProxyFiles(src string, nodes []objects.Node, version string) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	caFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	for _, node := range nodes {
		if utils.IsDirExist(path.Join(base, node.IP)) == false {
			os.MkdirAll(path.Join(base, node.IP), os.ModePerm)
		}

		var tpl bytes.Buffer
		var err error

		sslTp := template.New("proxy-csr")
		sslTp, _ = sslTp.Parse(t.KubeProxyCATemplate[version])
		sslTp.Execute(&tpl, nil)
		csrFileBytes := tpl.Bytes()

		req := csr.CertificateRequest{
			KeyRequest: csr.NewBasicKeyRequest(),
		}

		err = json.Unmarshal(csrFileBytes, &req)
		if err != nil {
			return err
		}

		var key, csrBytes []byte
		g := &csr.Generator{Validator: genkey.Validator}
		csrBytes, key, err = g.ProcessRequest(&req)
		if err != nil {
			return err
		}

		c := cli.Config{
			CAFile:     caFile,
			CAKeyFile:  caKeyFile,
			ConfigFile: configFile,
			Profile:    "kubernetes",
			Hostname:   "",
		}

		s, err := sign.SignerFromConfig(c)
		if err != nil {
			return err
		}

		var cert []byte
		signReq := signer.SignRequest{
			Request: string(csrBytes),
			Hosts:   signer.SplitHosts(c.Hostname),
			Profile: c.Profile,
		}

		cert, err = s.Sign(signReq)
		if err != nil {
			return err
		}

		var serviceTpl bytes.Buffer

		serviceTp := template.New("proxy-systemd")
		serviceTp, _ = serviceTp.Parse(t.KubeProxySystemdTemplate[version])
		serviceTp.Execute(&serviceTpl, map[string]string{"IP": node.IP})
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(path.Join(base, node.IP, "kube-proxy-csr.json"), csrFileBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, node.IP, "kube-proxy-key.pem"), key, 0600)
		err = ioutil.WriteFile(path.Join(base, node.IP, "kube-proxy.csr"), csrBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, node.IP, "kube-proxy.pem"), cert, 0600)
		err = ioutil.WriteFile(path.Join(base, node.IP, "kube-proxy.service"), serviceTpFileBytes, 0700)

		if err != nil {
			return err
		}

	}

	return nil
}

func uploadKubeProxyFiles(src, key string, nodes []objects.Node, stdout io.Writer) error {
	base := path.Join(src, tools.CAFilesFolder, tools.CAKubernetesFolder)

	for _, node := range nodes {
		var err error

		err = utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, "mkdir -p /var/lib/kube-proxy", stdout, os.Stderr)
		_, err = tools.DownloadComponent(path.Join(base, node.IP, "kube-proxy-csr.json"), "/etc/kubernetes/ssl/kube-proxy-csr.json", node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, node.IP, "kube-proxy-key.pem"), "/etc/kubernetes/ssl/kube-proxy-key.pem", node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, node.IP, "kube-proxy.csr"), "/etc/kubernetes/ssl/kube-proxy.csr", node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, node.IP, "kube-proxy.pem"), "/etc/kubernetes/ssl/kube-proxy.pem", node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, node.IP, "kube-proxy.service"), "/etc/systemd/system/kube-proxy.service", node.IP, key, node.User, stdout)
		_, err = tools.DownloadComponent(path.Join(base, node.IP, "kube-proxy.kubeconfig"), "/etc/kubernetes/kube-proxy.kubeconfig", node.IP, key, node.User, stdout)

		if err != nil {
			return err
		}

	}

	return nil
}

func startKubeProxy(key string, nodes []objects.Node) error {
	for _, node := range nodes {
		cmd := "systemctl daemon-reload && systemctl enable kube-proxy && systemctl start --no-block kube-proxy"

		if err := utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, cmd, os.Stdout, os.Stderr); err != nil {
			return err
		}

	}

	return nil
}

func generateKubeProxyConfigFile(src, ip, master string) error {
	cmdSetCluster := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "set-cluster", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
		fmt.Sprintf("--certificate-authority=%s", path.Join(src, "ssl", "root", "ca.pem")),
		"--embed-certs=true",
		fmt.Sprintf("--server=%s", master))
	cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
	if err := cmdSetCluster.Run(); err != nil {
		return err
	}

	cmdSetCredentials := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "set-credentials", "kube-proxy",
		fmt.Sprintf("--client-certificate=%s", path.Join(src, "ssl", "kubernetes", ip, "kube-proxy.pem")),
		fmt.Sprintf("--client-key=%s", path.Join(src, "ssl", "kubernetes", ip, "kube-proxy-key.pem")),
		"--embed-certs=true",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
	)
	cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
	if err := cmdSetCredentials.Run(); err != nil {
		return err
	}

	cmdSetContext := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "set-context", "default",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
		"--cluster=kubernetes", "--user=kube-proxy")
	cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
	if err := cmdSetContext.Run(); err != nil {
		return err
	}

	cmdUseContext := exec.Command(path.Join(src, "kubectl", "kubectl"), "config", "use-context",
		fmt.Sprintf("--kubeconfig=%s", path.Join(src, "ssl", "kubernetes", ip, "kube-proxy.kubeconfig")),
		"default")
	cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
	if err := cmdUseContext.Run(); err != nil {
		return err
	}

	return nil
}
