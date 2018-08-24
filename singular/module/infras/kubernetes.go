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
	"strings"
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

//DeployKubernetes is function deployment Kubernetes cluster include master and nodes.
//Notes:
//  1. Kubernetes master cluster IP.
//  2. Set kubectl config files.
//  3. Deploy Kubernetes master.
//  4. Deploy Kubernetes nodes.
func DeployKubernetesInCluster(d *objects.Deployment, infra *objects.Infra, stdout io.Writer, timestamp bool) error {
	//TODO Now singular only support one master and multiple nodes architect.
	//TODO So we decide the Kubernetes master IP is NODE_0 .
	masterIP := d.Outputs[fmt.Sprintf("NODE_%d", 0)].(string)
	apiServer := fmt.Sprintf("https://%s:6443", masterIP)
	etcdEndpoints := d.Outputs["EtcdEndpoints"].(string)

	//Master IP
	d.Output("MASTER_IP", masterIP)
	d.Output("KUBE_APISERVER", fmt.Sprintf("https://%s:6443", masterIP))

	//Master Nodes
	kubeMasterNodes := []*objects.Node{}
	for i := 0; i < infra.Master; i++ {
		kubeMasterNodes = append(kubeMasterNodes, d.Nodes[i])
	}

	//Minion Nodes
	kubeSlaveNodes := []*objects.Node{}
	for i := 0; i < infra.Minion; i++ {
		kubeSlaveNodes = append(kubeSlaveNodes, d.Nodes[i])
	}

	//Download binary in master nodes
	for _, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, kubeMasterNodes, stdout, timestamp); err != nil {
			return err
		}
	}

	//Download binary in slave nodes
	for _, c := range infra.Components {
		objects.WriteLog(fmt.Sprintf("Download %s in deploy notes", c.Binary), stdout, timestamp, d, infra, c)
		if err := d.DownloadBinaryFile(c.Binary, c.URL, kubeSlaveNodes, stdout, timestamp); err != nil {
			return err
		}
	}

	for _, c := range infra.Components {
		switch c.Binary {
		case "kubectl":
			//Generate kubectl config file and CA SSL files.
			if err := setKubectlFiles(d, c.URL, apiServer, stdout, timestamp); err != nil {
				return err
			}

			//Upload kubectl config file to Kubernetes nodes
			if err := uploadKubeConfigFiles(d, d.Tools.SSH.Private, kubeSlaveNodes, stdout); err != nil {
				return err
			}
		case "kube-apiserver":
			//Generate Kubernetes Token API file
			if files, err := generateTokenFile(d, stdout, timestamp); err != nil {
				return err
			} else {
				//Upload Kubernetes Token file
				if err := uploadTokenFiles(kubeMasterNodes[0].User, files, d.Tools.SSH.Private, masterIP, stdout); err != nil {
					return err
				}
			}

			//Generate Kubernetes SSL files and systemd service file
			if files, err := generateKubeAPIServerFiles(d, masterIP, etcdEndpoints, infra.Version, stdout, timestamp); err != nil {
				return err
			} else {
				//Upload Kubernetes API Server SSL files and systemd service file
				if err := uploadKubeAPIServerCAFiles(files, d, kubeMasterNodes, stdout); err != nil {
					return err
				}

				//Start Kubernetes API Server
				if err := startKubeAPIServer(d, kubeMasterNodes, stdout); err != nil {
					return err
				}
			}
		case "kube-controller-manager":
			//Generate Kube-controller-manager systemd service file
			if files, err := generateKubeControllerManagerFiles(d, masterIP, etcdEndpoints, infra.Version); err != nil {
				return err
			} else {
				//Upload Kube-controller-manager systemd service file
				if err := uploadKubeControllerFiles(files, d, kubeMasterNodes, stdout); err != nil {
					return err
				}

				//Start Kube-controller-manager
				if err := startKubeController(d, kubeMasterNodes, stdout); err != nil {
					return err
				}
			}
		case "kube-scheduler":
			//Generate Kube-scheduler systemd service file
			if files, err := generateKubeSchedulerFiles(d, masterIP, etcdEndpoints, infra.Version); err != nil {
				return err
			} else {
				//Upload Kube-scheduler systemd service file
				if err := uploadKubeSchedulerFiles(files, d, kubeMasterNodes, stdout); err != nil {
					return err
				}

				//Start Kube-scheduler
				if err := startKubeScheduler(d, kubeMasterNodes, stdout); err != nil {
					return err
				}
			}
		case "kubelet":
			//generate bootstrap.kubeconfig config file
			if config, err := generateBootstrapFile(d, apiServer, stdout); err != nil {
				return err
			} else {
				if err := uploadBootstrapFile(config, d, kubeSlaveNodes, stdout); err != nil {
					return err
				}
			}

			//exec kubectl create clusterrolebinding command
			if err := setKubeletClusterrolebinding(d, kubeSlaveNodes, stdout); err != nil {
				return nil
			}

			//generate kubelet systemd service file, then upload to the nodes and start service.
			if files, err := generateKubeletSystemdFile(d, kubeSlaveNodes, infra.Version); err != nil {
				return err
			} else {
				if err := uploadKubeletFile(files, d, kubeSlaveNodes, stdout); err != nil {
					return err
				}

				if err := startKubelet(d, kubeSlaveNodes, stdout); err != nil {
					return err
				}
			}

			//Sleep 60 seconds waiting kubelet service start.
			objects.WriteLog("Time sleep 60 seconds for awaiting server", stdout, timestamp, d)
			time.Sleep(60 * time.Second)

			//Approve kubelet certificate key
			if err := kubeletCertificateApprove(d, kubeSlaveNodes, stdout); err != nil {
				return err
			}
		case "kube-proxy":
			//Generate Kube Proxy Systemd template file
			if files, err := generateKubeProxyFiles(d, kubeSlaveNodes, infra.Version); err != nil {
				return err
			} else {
				if err := generateKubeProxyConfigFile(&files, d, kubeSlaveNodes, apiServer); err != nil {
					return err
				}

				//Upload kube-proxy Systemd file
				if err := uploadKubeProxyFiles(files, d, kubeSlaveNodes, stdout); err != nil {
					return err
				}

				//Start kube-proxy Service
				if err := startKubeProxy(d, kubeSlaveNodes, stdout); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupport kubenetes component: %s", c.Binary)
		}
	}

	return nil
}

//setKubectlConfig download kubectl and set config file.
func setKubectlFiles(d *objects.Deployment, link, apiServer string, stdout io.Writer, timestamp bool) error {
	//Make kubectl folder
	if utils.IsDirExist(path.Join(d.Config, tools.KubectlFileFolder)) == true {
		os.RemoveAll(path.Join(d.Config, tools.KubectlFileFolder))
	}
	os.MkdirAll(path.Join(d.Config, tools.KubectlFileFolder), os.ModePerm)

	//Download or copy kubectl file
	if a, err := url.Parse(link); err != nil {
		return err
	} else {
		if a.Scheme == "" {
			cmdCopy := exec.Command("cp", link, path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile))
			cmdCopyStr := fmt.Sprintf("cp %s %s", link, path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile))
			objects.WriteLog(fmt.Sprintf("%s", cmdCopyStr), stdout, timestamp, d)
			cmdCopy.Stdout, cmdCopy.Stderr = stdout, os.Stderr
			if err := cmdCopy.Run(); err != nil {
				return err
			}
		} else {
			cmdDownload := exec.Command("curl", link, "-o", path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile))
			objects.WriteLog(fmt.Sprintf("curl %s -o %s", link, path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile)), stdout, timestamp, d)
			cmdDownload.Stdout, cmdDownload.Stderr = stdout, os.Stderr
			if err := cmdDownload.Run(); err != nil {
				return err
			}
		}
	}

	//Change mode +x for kubectl
	cmdChmod := exec.Command("chmod", "+x", path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile))
	objects.WriteLog(fmt.Sprintf("chmod +x %s", path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile)), stdout, timestamp, d)
	cmdChmod.Stdout, cmdChmod.Stderr = stdout, os.Stderr
	if err := cmdChmod.Run(); err != nil {
		return err
	}

	//Generate Kubernetes admin CA files
	if files, err := generateKubeAdminCAFiles(d.Config); err != nil {
		return err
	} else {
		objects.WriteLog(fmt.Sprintf("Kubernetes CA Admin files [%v]", files), stdout, timestamp, d)
	}

	//Generate kubectl config file
	if file, err := setKubeConfigFile(d, apiServer, stdout); err != nil {
		return err
	} else {
		objects.WriteLog(fmt.Sprintf("kubectl config files [%s]", file), stdout, timestamp, d)
	}

	return nil
}

//generateKubeAdminCAFiles generate Kubernetes Admin CA files
func generateKubeAdminCAFiles(src string) (map[string]string, error) {
	base := path.Join(src, tools.KubectlFileFolder)

	files := map[string]string{
		tools.CAKubeAdminCSRConfigFile: path.Join(base, tools.CAKubeAdminCSRConfigFile),
		tools.CAKubeAdminKeyPemFile:    path.Join(base, tools.CAKubeAdminKeyPemFile),
		tools.CAKubeAdminCSRFile:       path.Join(base, tools.CAKubeAdminCSRFile),
		tools.CAKubeAdminPemFile:       path.Join(base, tools.CAKubeAdminPemFile),
	}

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
		return files, err
	}

	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		return files, err
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
		return files, err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(c.Hostname),
		Profile: c.Profile,
	}

	cert, err = s.Sign(signReq)
	if err != nil {
		return files, err
	}

	err = ioutil.WriteFile(files[tools.CAKubeAdminCSRConfigFile], csrFileBytes, 0600)
	err = ioutil.WriteFile(files[tools.CAKubeAdminKeyPemFile], key, 0600)
	err = ioutil.WriteFile(files[tools.CAKubeAdminCSRFile], csrBytes, 0600)
	err = ioutil.WriteFile(files[tools.CAKubeAdminPemFile], cert, 0600)

	if err != nil {
		return files, err
	}

	return files, nil
}

//setKubeConfigFile generate kubectl config file
func setKubeConfigFile(d *objects.Deployment, apiServer string, stdout io.Writer) (string, error) {
	base := path.Join(d.Config, tools.KubectlFileFolder)
	adminPem := path.Join(base, tools.CAKubeAdminPemFile)
	adminKey := path.Join(base, tools.CAKubeAdminKeyPemFile)

	kubectl := path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile)
	config := path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlConfigFile)
	caFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)

	cmdSetCluster := exec.Command(kubectl, "config", "set-cluster", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", config),
		fmt.Sprintf("--certificate-authority=%s", caFile),
		"--embed-certs=true",
		fmt.Sprintf("--server=%s", apiServer))
	cmdSetCluster.Stdout, cmdSetCluster.Stderr = stdout, os.Stderr
	if err := cmdSetCluster.Run(); err != nil {
		return caFile, err
	}

	cmdSetCredentials := exec.Command(kubectl, "config", "set-credentials", "admin",
		fmt.Sprintf("--kubeconfig=%s", config),
		fmt.Sprintf("--client-certificate=%s", adminPem),
		"--embed-certs=true",
		fmt.Sprintf("--client-key=%s", adminKey))
	cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = stdout, os.Stderr
	if err := cmdSetCredentials.Run(); err != nil {
		return caFile, err
	}

	cmdSetContext := exec.Command(kubectl, "config", "set-context", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", config),
		"--cluster=kubernetes", "--user=admin")
	cmdSetContext.Stdout, cmdSetContext.Stderr = stdout, os.Stderr
	if err := cmdSetContext.Run(); err != nil {
		return caFile, err
	}

	cmdUseContext := exec.Command(kubectl, "config", "use-context",
		fmt.Sprintf("--kubeconfig=%s", config),
		"kubernetes")
	cmdUseContext.Stdout, cmdUseContext.Stderr = stdout, os.Stderr
	if err := cmdUseContext.Run(); err != nil {
		return caFile, err
	}

	return caFile, nil
}

//uploadKubeConfigFiles upload Kubectl config file and ca ssl files.
func uploadKubeConfigFiles(d *objects.Deployment, key string, nodes []*objects.Node, stdout io.Writer) error {
	base := path.Join(d.Config, tools.KubectlFileFolder)

	for _, node := range nodes {
		var cmd, dest string

		if node.User == tools.DefaultSSHUser {
			cmd = fmt.Sprintf("mkdir -p /%s/.kube", tools.DefaultSSHUser)
			dest = fmt.Sprintf("/%s/.kube/config", tools.DefaultSSHUser)
		} else {
			cmd = fmt.Sprintf("mkdir -p /home/%s/.kube", node.User)
			dest = fmt.Sprintf("/home/%s/.kube", node.User)
		}

		files := []map[string]string{
			{
				"src":  path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlConfigFile),
				"dest": dest,
			},
			{
				"src":  path.Join(base, tools.CAKubeAdminCSRConfigFile),
				"dest": path.Join(tools.KubeServerConfig, tools.KubeServerSSL, tools.CAKubeAdminCSRConfigFile),
			},
			{
				"src":  path.Join(base, tools.CAKubeAdminKeyPemFile),
				"dest": path.Join(tools.KubeServerConfig, tools.KubeServerSSL, tools.CAKubeAdminKeyPemFile),
			},
			{
				"src":  path.Join(base, tools.CAKubeAdminCSRFile),
				"dest": path.Join(tools.KubeServerConfig, tools.KubeServerSSL, tools.CAKubeAdminCSRFile),
			},
			{
				"src":  path.Join(base, tools.CAKubeAdminPemFile),
				"dest": path.Join(tools.KubeServerConfig, tools.KubeServerSSL, tools.CAKubeAdminPemFile),
			},
		}

		if err := utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, []string{cmd}, stdout, os.Stderr); err != nil {
			return err
		}

		if err := tools.DownloadComponent(files, node.IP, key, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//generateTokenFile generate Kubectl Token file
func generateTokenFile(d *objects.Deployment, stdout io.Writer, timestamp bool) (map[string]string, error) {
	if utils.IsDirExist(path.Join(d.Config, "kubectl")) == false {
		os.MkdirAll(path.Join(d.Config, "kubectl"), os.ModePerm)
	}

	var tokenTpl bytes.Buffer

	files := map[string]string{
		tools.KubeTokenCSVFile: path.Join(d.Config, tools.KubectlFileFolder, tools.KubeTokenCSVFile),
	}

	//TODO generate token string
	tokenTp := template.New("token")
	tokenTp, _ = tokenTp.Parse("{{.Token}},kubelet-bootstrap,10001,\"system:kubelet-bootstrap\"")
	tokenTp.Execute(&tokenTpl, map[string]string{"Token": t.BooststrapToken})
	tokenTpFileBytes := tokenTpl.Bytes()

	objects.WriteLog(fmt.Sprintf("Write kubenetes token csv file %s", files[tools.KubeTokenCSVFile]), stdout, timestamp, d)
	if err := ioutil.WriteFile(files[tools.KubeTokenCSVFile], tokenTpFileBytes, 0700); err != nil {
		return files, err
	}

	return files, nil
}

//uploadTokenFiles upload Token CSV file
func uploadTokenFiles(user string, f map[string]string, key, ip string, stdout io.Writer) error {
	files := []map[string]string{{
		"src":  f[tools.KubeTokenCSVFile],
		"dest": path.Join(tools.KubeServerConfig, tools.KubeTokenCSVFile),
	}}

	if err := tools.DownloadComponent(files, ip, key, user, stdout); err != nil {
		return err
	}

	return nil
}

//KubeMaster is kubernetes master template struct
type KubeMaster struct {
	MasterIP string
	Nodes    string
}

//generateKubeAPIServerFiles generate Kube API Server CA SSL files.
func generateKubeAPIServerFiles(d *objects.Deployment, masterIP, etcdEndpoints string, version string, stdout io.Writer, timestamp bool) (map[string]string, error) {
	base := path.Join(d.Config, tools.CAFilesFolder, tools.CAKubernetesFolder)
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, os.ModePerm)

	caFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	master := KubeMaster{
		MasterIP: masterIP,
		Nodes:    etcdEndpoints,
	}

	files := map[string]string{
		tools.CAKubeAPIServerCSRConfigFile: path.Join(base, tools.CAKubeAPIServerCSRConfigFile),
		tools.CAKubeAPIServerKeyPemFile:    path.Join(base, tools.CAKubeAPIServerKeyPemFile),
		tools.CAKubeAPIServerCSRFile:       path.Join(base, tools.CAKubeAPIServerCSRFile),
		tools.CAKubeAPIServerPemFile:       path.Join(base, tools.CAKubeAPIServerPemFile),
		tools.KubeAPIServerSystemdFile:     path.Join(base, tools.KubeAPIServerSystemdFile),
	}

	var tpl bytes.Buffer
	var err error

	// Handle the version
	tplContent := t.GetK8sCATemplate(version)
	if tplContent == "" {
		return files, fmt.Errorf("No kubernetes ca template for version %s or its parent version", version)
	}

	sslTp := template.New("kube-csr")
	sslTp, _ = sslTp.Parse(tplContent)
	sslTp.Execute(&tpl, master)
	csrFileBytes := tpl.Bytes()

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}

	err = json.Unmarshal(csrFileBytes, &req)
	if err != nil {
		return files, err
	}

	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		return files, err
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
		return files, err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(c.Hostname),
		Profile: c.Profile,
	}

	cert, err = s.Sign(signReq)
	if err != nil {
		return files, err
	}

	var serviceTpl bytes.Buffer

	tplContent = t.GetK8sAPIServerSystemdTemplate(version)
	if tplContent == "" {
		return files, fmt.Errorf("No kube-apiserver systemd template for version %s or its parent version", version)
	}
	serviceTp := template.New("kube-systemd")
	serviceTp, _ = serviceTp.Parse(tplContent)
	serviceTp.Execute(&serviceTpl, master)
	serviceTpFileBytes := serviceTpl.Bytes()

	err = ioutil.WriteFile(files[tools.CAKubeAPIServerCSRConfigFile], csrFileBytes, 0600)
	err = ioutil.WriteFile(files[tools.CAKubeAPIServerKeyPemFile], key, 0600)
	err = ioutil.WriteFile(files[tools.CAKubeAPIServerCSRFile], csrBytes, 0600)
	err = ioutil.WriteFile(files[tools.CAKubeAPIServerPemFile], cert, 0600)
	err = ioutil.WriteFile(files[tools.KubeAPIServerSystemdFile], serviceTpFileBytes, 0700)

	if err != nil {
		return files, err
	}

	return files, nil
}

//uploadKubeAPIServerCAFiles upload Kube API Server systemd file and CA SSL file.
func uploadKubeAPIServerCAFiles(f map[string]string, d *objects.Deployment, masters []*objects.Node, stdout io.Writer) error {
	for _, node := range masters {
		files := []map[string]string{}

		for k, file := range f {
			if k == tools.KubeAPIServerSystemdFile {
				files = append(files, map[string]string{
					"src":  file,
					"dest": path.Join(tools.SystemdServerPath, tools.KubeAPIServerSystemdFile),
				})
			} else {
				files = append(files, map[string]string{
					"src":  file,
					"dest": path.Join(tools.KubeServerConfig, tools.KubeServerSSL, k),
				})
			}
		}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//startKubeAPIServer start Kube-APIServer in the master nodes.
func startKubeAPIServer(d *objects.Deployment, masters []*objects.Node, stdout io.Writer) error {
	commands := []string{
		"systemctl daemon-reload",
		"systemctl enable kube-apiserver",
		"systemctl start --no-block kube-apiserver",
	}

	for _, node := range masters {
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, commands, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

//generateKubeControllerManagerFiles generate kube-controller-manager systemd service file.
func generateKubeControllerManagerFiles(d *objects.Deployment, masterIP, etcdEndpoints string, version string) (map[string]string, error) {
	base := path.Join(d.Config, tools.CAFilesFolder, tools.CAKubernetesFolder)

	files := map[string]string{
		tools.KubeControllerManagerSystemdFile: path.Join(base, tools.KubeControllerManagerSystemdFile),
	}

	master := KubeMaster{
		MasterIP: masterIP,
		Nodes:    etcdEndpoints,
	}

	var serviceTpl bytes.Buffer

	tplContent := t.GetK8sControllerManagerSystemdTemplate(version)
	if tplContent == "" {
		return files, fmt.Errorf("No kube-controller-manager template for version %s or its parent version", version)
	}
	serviceTp := template.New("kube-control-systemd")
	serviceTp, _ = serviceTp.Parse(tplContent)
	serviceTp.Execute(&serviceTpl, master)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(files[tools.KubeControllerManagerSystemdFile], serviceTpFileBytes, 0700); err != nil {
		return files, err
	}

	return files, nil
}

//uploadKubeControllerFiles upload kube-controller-manager
func uploadKubeControllerFiles(f map[string]string, d *objects.Deployment, masters []*objects.Node, stdout io.Writer) error {
	for _, node := range masters {
		files := []map[string]string{{
			"src":  f[tools.KubeControllerManagerSystemdFile],
			"dest": path.Join(tools.SystemdServerPath, tools.KubeControllerManagerSystemdFile),
		}}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//startKubeController start kube-controller-manager in the master nodes.
func startKubeController(d *objects.Deployment, masters []*objects.Node, stdout io.Writer) error {
	commands := []string{
		"systemctl daemon-reload",
		"systemctl enable kube-controller-manager",
		"systemctl start --no-block kube-controller-manager",
	}

	for _, node := range masters {
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, commands, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

//generateKubeSchedulerFiles generate kube-scheduler systemd service file.
func generateKubeSchedulerFiles(d *objects.Deployment, masterIP, etcdEndpoints string, version string) (map[string]string, error) {
	base := path.Join(d.Config, tools.CAFilesFolder, tools.CAKubernetesFolder)

	files := map[string]string{
		tools.KubeSchedulerSystemdFile: path.Join(base, tools.KubeSchedulerSystemdFile),
	}

	master := KubeMaster{
		MasterIP: masterIP,
		Nodes:    etcdEndpoints,
	}

	var serviceTpl bytes.Buffer

	tplContent := t.GetK8sSchedulerSystemdTemplate(version)
	if tplContent == "" {
		return files, fmt.Errorf("No kube-scheduler template for version %s or its parent version", version)
	}
	serviceTp := template.New("kube-scheduler-systemd")
	serviceTp, _ = serviceTp.Parse(tplContent)
	serviceTp.Execute(&serviceTpl, master)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(files[tools.KubeSchedulerSystemdFile], serviceTpFileBytes, 0700); err != nil {
		return files, err
	}

	return files, nil
}

//uploadKubeSchedulerFiles
func uploadKubeSchedulerFiles(files map[string]string, d *objects.Deployment, masters []*objects.Node, stdout io.Writer) error {
	for _, node := range masters {
		files := []map[string]string{{
			"src":  files[tools.KubeSchedulerSystemdFile],
			"dest": path.Join(tools.SystemdServerPath, tools.KubeSchedulerSystemdFile),
		}}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//startKubeScheduler start kube-scheduler in the master nodes.
func startKubeScheduler(d *objects.Deployment, masters []*objects.Node, stdout io.Writer) error {
	commands := []string{
		"systemctl daemon-reload",
		"systemctl enable kube-scheduler",
		"systemctl start --no-block kube-scheduler",
	}

	for _, node := range masters {
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, commands, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

//generateBootstrapFile generate bootstrap.kubeconfig config file
func generateBootstrapFile(d *objects.Deployment, apiServer string, stdout io.Writer) (string, error) {
	kubectl := path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile)
	kubeconfig := path.Join(d.Config, tools.KubectlFileFolder, tools.KubeBootstrapConfig)
	caFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)

	//Generate bootstrap.kubeconfig config file
	cmdSetCluster := exec.Command(kubectl, "config", "set-cluster", "kubernetes",
		fmt.Sprintf("--kubeconfig=%s", kubeconfig),
		fmt.Sprintf("--certificate-authority=%s", caFile),
		"--embed-certs=true",
		fmt.Sprintf("--server=%s", apiServer))
	cmdSetCluster.Stdout, cmdSetCluster.Stderr = stdout, os.Stderr
	if err := cmdSetCluster.Run(); err != nil {
		return kubeconfig, err
	}

	cmdSetCredentials := exec.Command(kubectl, "config", "set-credentials", "kubelet-bootstrap",
		fmt.Sprintf("--kubeconfig=%s", kubeconfig),
		fmt.Sprintf("--token=%s", t.BooststrapToken))
	cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = stdout, os.Stderr
	if err := cmdSetCredentials.Run(); err != nil {
		return kubeconfig, err
	}

	cmdSetContext := exec.Command(kubectl, "config", "set-context", "default",
		fmt.Sprintf("--kubeconfig=%s", kubeconfig),
		"--cluster=kubernetes", "--user=kubelet-bootstrap")
	cmdSetContext.Stdout, cmdSetContext.Stderr = stdout, os.Stderr
	if err := cmdSetContext.Run(); err != nil {
		return kubeconfig, err
	}

	cmdUseContext := exec.Command(kubectl, "config", "use-context",
		fmt.Sprintf("--kubeconfig=%s", kubeconfig),
		"default")
	cmdUseContext.Stdout, cmdUseContext.Stderr = stdout, os.Stderr
	if err := cmdUseContext.Run(); err != nil {
		return kubeconfig, err
	}

	return kubeconfig, nil
}

//uploadBootstrapFile upload bootstrap.kubeconfig to kubernetes slave nodes.
func uploadBootstrapFile(file string, d *objects.Deployment, kubeSlaveNodes []*objects.Node, stdout io.Writer) error {
	for _, node := range kubeSlaveNodes {
		files := []map[string]string{{
			"src":  file,
			"dest": path.Join(tools.KubeServerConfig, tools.KubeBootstrapConfig),
		}}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//setKubeletClusterrolebinding exec kubectl create clusterrolebinding command in the first slave node in Kubernetes clusters.
func setKubeletClusterrolebinding(d *objects.Deployment, nodes []*objects.Node, stdout io.Writer) error {
	// Make sure the kube-apiserver is ready to serve(retry for 60 seconds at most)
	masterIP := nodes[0].IP
	if err := utils.WaitForHostPort(masterIP, 6443, 3, 20); err != nil {
		return err
	}
	time.Sleep(time.Second * 5)
	for i, node := range nodes {
		if i == 0 {
			retries := 0
			for {
				cmd := "/usr/local/bin/kubectl create clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap"
				if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, []string{cmd}, stdout, os.Stderr); err == nil {
					break
				} else if retries >= 10 {
					return err
				}
				retries++
			}
		}
	}

	return nil
}

//generateKubeletSystemdFile generate kubelet systemd service file.
func generateKubeletSystemdFile(d *objects.Deployment, nodes []*objects.Node, version string) (map[string]map[string]string, error) {
	files := map[string]map[string]string{}

	for _, node := range nodes {

		cmd := "which iptables"
		var buf bytes.Buffer
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, []string{cmd}, &buf, ioutil.Discard); err != nil {
			return nil, err
		}
		iptablesAbsolutePath := strings.TrimSpace(buf.String())

		kubeNode := map[string]string{
			"IP":       node.IP,
			"iptables": iptablesAbsolutePath,
		}

		files[node.IP] = map[string]string{}

		base := path.Join(d.Config, tools.CAFilesFolder, tools.CAKubernetesFolder, node.IP)
		if utils.IsDirExist(base) == true {
			os.RemoveAll(base)
		}

		os.MkdirAll(base, os.ModePerm)
		files[node.IP][tools.KubeletSystemdFile] = path.Join(base, tools.KubeletSystemdFile)

		var serviceTpl bytes.Buffer

		tplContent := t.GetKubeletSystemdTemplate(version)
		if tplContent == "" {
			return files, fmt.Errorf("No kubelet systemd template for version %s or its parent version", version)
		}
		serviceTp := template.New("kubelet-systemd")
		serviceTp, _ = serviceTp.Parse(tplContent)
		serviceTp.Execute(&serviceTpl, kubeNode)
		serviceTpFileBytes := serviceTpl.Bytes()

		if err := ioutil.WriteFile(files[node.IP][tools.KubeletSystemdFile], serviceTpFileBytes, 0700); err != nil {
			return files, err
		}
	}

	return files, nil
}

//uploadKubeletFile upload kubelet systemd service file to slave nodes.
func uploadKubeletFile(files map[string]map[string]string, d *objects.Deployment, kubeSlaveNodes []*objects.Node, stdout io.Writer) error {
	for _, node := range kubeSlaveNodes {
		cmd := "mkdir -p /var/lib/kubelet"
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, []string{cmd}, stdout, os.Stderr); err != nil {
			return err
		}

		files := []map[string]string{{
			"src":  files[node.IP][tools.KubeletSystemdFile],
			"dest": path.Join(tools.SystemdServerPath, tools.KubeletSystemdFile),
		},
		}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//startKubelet start kubelet service in the slave nodes.
func startKubelet(d *objects.Deployment, kubeSlaveNodes []*objects.Node, stdout io.Writer) error {
	for _, node := range kubeSlaveNodes {
		commands := []string{
			"systemctl daemon-reload",
			"systemctl enable kubelet",
			"systemctl start --no-block kubelet",
		}

		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, commands, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

//kubeletCertificateApprove approve kubelet certificate
func kubeletCertificateApprove(d *objects.Deployment, nodes []*objects.Node, stdout io.Writer) error {
	for i, node := range nodes {
		if i == 0 {
			cmd := "/usr/local/bin/kubectl certificate approve `/usr/local/bin/kubectl get csr -o name`"
			if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, []string{cmd}, stdout, os.Stderr); err != nil {
				return err
			}
		}
	}

	return nil
}

//generateKubeProxyFiles generate CA files and systemd file.
func generateKubeProxyFiles(d *objects.Deployment, kubeSlaveNodes []*objects.Node, version string) (map[string]map[string]string, error) {
	base := path.Join(d.Config, tools.CAFilesFolder, tools.CAKubernetesFolder)

	caFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	files := map[string]map[string]string{}

	for _, node := range kubeSlaveNodes {
		if utils.IsDirExist(path.Join(base, node.IP)) == false {
			os.MkdirAll(path.Join(base, node.IP), os.ModePerm)
		}

		var tpl bytes.Buffer
		var err error

		files[node.IP] = map[string]string{}
		files[node.IP][tools.CAKubeProxyServerCSRConfigFile] = path.Join(base, node.IP, tools.CAKubeProxyServerCSRConfigFile)
		files[node.IP][tools.CAKubeProxyServerKeyPemFile] = path.Join(base, node.IP, tools.CAKubeProxyServerKeyPemFile)
		files[node.IP][tools.CAKubeProxyServerCSR] = path.Join(base, node.IP, tools.CAKubeProxyServerCSR)
		files[node.IP][tools.CAKubeProxyServerPemFile] = path.Join(base, node.IP, tools.CAKubeProxyServerPemFile)
		files[node.IP][tools.KubeProxySystemdFiles] = path.Join(base, node.IP, tools.KubeProxySystemdFiles)

		tplContent := t.GetKubeProxyCATemplate(version)
		if tplContent == "" {
			return files, fmt.Errorf("No kube-proxy ca template for version %s or its parent version", version)
		}
		sslTp := template.New("proxy-csr")
		sslTp, _ = sslTp.Parse(tplContent)
		sslTp.Execute(&tpl, nil)
		csrFileBytes := tpl.Bytes()

		req := csr.CertificateRequest{
			KeyRequest: csr.NewBasicKeyRequest(),
		}

		err = json.Unmarshal(csrFileBytes, &req)
		if err != nil {
			return files, err
		}

		var key, csrBytes []byte
		g := &csr.Generator{Validator: genkey.Validator}
		csrBytes, key, err = g.ProcessRequest(&req)
		if err != nil {
			return files, err
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
			return files, err
		}

		var cert []byte
		signReq := signer.SignRequest{
			Request: string(csrBytes),
			Hosts:   signer.SplitHosts(c.Hostname),
			Profile: c.Profile,
		}

		cert, err = s.Sign(signReq)
		if err != nil {
			return files, err
		}

		var serviceTpl bytes.Buffer

		tplContent = t.GetKubeProxySystemdTemplate(version)
		if tplContent == "" {
			return files, fmt.Errorf("No kube-proxy systemd template for version %s or its parent version", version)
		}
		serviceTp := template.New("proxy-systemd")
		serviceTp, _ = serviceTp.Parse(tplContent)
		serviceTp.Execute(&serviceTpl, map[string]string{"IP": node.IP})
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(files[node.IP][tools.CAKubeProxyServerCSRConfigFile], csrFileBytes, 0600)
		err = ioutil.WriteFile(files[node.IP][tools.CAKubeProxyServerKeyPemFile], key, 0600)
		err = ioutil.WriteFile(files[node.IP][tools.CAKubeProxyServerCSR], csrBytes, 0600)
		err = ioutil.WriteFile(files[node.IP][tools.CAKubeProxyServerPemFile], cert, 0600)
		err = ioutil.WriteFile(files[node.IP][tools.KubeProxySystemdFiles], serviceTpFileBytes, 0700)

		if err != nil {
			return files, err
		}

	}

	return files, nil
}

//generateKubeProxyConfigFile generate kube-proxy.kubeconfig file.
func generateKubeProxyConfigFile(files *map[string]map[string]string, d *objects.Deployment, kubeSlaveNodes []*objects.Node, apiServer string) error {
	kubectl := path.Join(d.Config, tools.KubectlFileFolder, tools.KubectlFile)
	caFile := path.Join(d.Config, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)

	for _, node := range kubeSlaveNodes {
		config := path.Join(d.Config, tools.CAFilesFolder, tools.CAKubernetesFolder, node.IP, tools.KubeProxyConfigFile)
		(*files)[node.IP][tools.KubeProxyConfigFile] = config

		cmdSetCluster := exec.Command(kubectl, "config", "set-cluster", "kubernetes",
			fmt.Sprintf("--kubeconfig=%s", config),
			fmt.Sprintf("--certificate-authority=%s", caFile),
			"--embed-certs=true",
			fmt.Sprintf("--server=%s", apiServer))
		cmdSetCluster.Stdout, cmdSetCluster.Stderr = os.Stdout, os.Stderr
		if err := cmdSetCluster.Run(); err != nil {
			return err
		}

		cmdSetCredentials := exec.Command(kubectl, "config", "set-credentials", "kube-proxy",
			fmt.Sprintf("--client-certificate=%s", (*files)[node.IP][tools.CAKubeProxyServerPemFile]),
			fmt.Sprintf("--client-key=%s", (*files)[node.IP][tools.CAKubeProxyServerKeyPemFile]),
			"--embed-certs=true",
			fmt.Sprintf("--kubeconfig=%s", config),
		)
		cmdSetCredentials.Stdout, cmdSetCredentials.Stderr = os.Stdout, os.Stderr
		if err := cmdSetCredentials.Run(); err != nil {
			return err
		}

		cmdSetContext := exec.Command(kubectl, "config", "set-context", "default",
			fmt.Sprintf("--kubeconfig=%s", config), "--cluster=kubernetes", "--user=kube-proxy")
		cmdSetContext.Stdout, cmdSetContext.Stderr = os.Stdout, os.Stderr
		if err := cmdSetContext.Run(); err != nil {
			return err
		}

		cmdUseContext := exec.Command(kubectl, "config", "use-context",
			fmt.Sprintf("--kubeconfig=%s", config),
			"default")
		cmdUseContext.Stdout, cmdUseContext.Stderr = os.Stdout, os.Stderr
		if err := cmdUseContext.Run(); err != nil {
			return err
		}

	}

	return nil
}

//uploadKubeProxyFiles upload kube-proxy CA files, systemd service file and kube-proxy.kubeconfig file to the nodes.
func uploadKubeProxyFiles(f map[string]map[string]string, d *objects.Deployment, kubeSlaveNodes []*objects.Node, stdout io.Writer) error {
	for _, node := range kubeSlaveNodes {
		cmd := "mkdir -p /var/lib/kube-proxy"
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, []string{cmd}, stdout, os.Stderr); err != nil {
			return err
		}

		files := []map[string]string{}
		for k, file := range f[node.IP] {
			switch k {
			case tools.KubeProxySystemdFiles:
				files = append(files, map[string]string{
					"src":  file,
					"dest": path.Join(tools.SystemdServerPath, tools.KubeProxySystemdFiles),
				})
			case tools.KubeProxyConfigFile:
				files = append(files, map[string]string{
					"src":  file,
					"dest": path.Join(tools.KubeServerConfig, tools.KubeProxyConfigFile),
				})
			default:
				files = append(files, map[string]string{
					"src":  file,
					"dest": path.Join(EtcdServerConfig, EtcdServerSSL, k),
				})
			}
		}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}
	}

	return nil
}

//startKubeProxy start kube-proxy service in the nodes.
func startKubeProxy(d *objects.Deployment, kubeSlaveNodes []*objects.Node, stdout io.Writer) error {
	commands := []string{
		"systemctl daemon-reload",
		"systemctl enable kube-proxy",
		"systemctl start --no-block kube-proxy",
	}

	for _, node := range kubeSlaveNodes {
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, commands, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}
