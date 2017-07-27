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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"

	"github.com/Huawei/containerops/common/utils"
	t "github.com/Huawei/containerops/singular/module/template"
)

func GenerateCARootFiles(src string) (map[string]string, error) {
	var caConfigTpl, caCsrTpl bytes.Buffer
	var err error

	base := path.Join(src, "ssl", "root")
	if utils.IsDirExist(base) == false {
		os.MkdirAll(base, os.ModePerm)
	}

	files := map[string]string{
		"CaConfigFile":    path.Join(base, "ca-config.json"),
		"CaCsrConfigFile": path.Join(base, "ca-csr.json"),
		"CaPemFile":       path.Join(base, "ca.pem"),
		"CaCsrFile":       path.Join(base, "ca.csr"),
		"CaKeyFile":       path.Join(base, "ca-key.pem"),
	}

	for _, value := range files {
		if utils.IsDirExist(value) == true {
			err = os.Remove(value)
			if err != nil {
				return files, err
			}
		}
	}

	caConfig := template.New("config")
	caConfig, err = caConfig.Parse(t.Root["ca-config"])
	caConfig.Execute(&caConfigTpl, nil)
	err = ioutil.WriteFile(files["CaConfigFile"], caConfigTpl.Bytes(), 0600)
	if err != nil {
		return files, err
	}

	caCsr := template.New("csr")
	caCsr, err = caCsr.Parse(t.Root["ca-csr"])
	caCsr.Execute(&caCsrTpl, nil)
	err = ioutil.WriteFile(files["CaCsrConfigFile"], caCsrTpl.Bytes(), 0600)
	if err != nil {
		return files, err
	}

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}
	err = json.Unmarshal(caCsrTpl.Bytes(), &req)
	if err != nil {
		return files, err
	}

	var key, csrPEM, cert []byte
	cert, csrPEM, key, err = initca.New(&req)
	err = ioutil.WriteFile(files["CaPemFile"], cert, 0600)
	err = ioutil.WriteFile(files["CaCsrFile"], csrPEM, 0600)
	err = ioutil.WriteFile(files["CaKeyFile"], key, 0600)

	return files, err
}

type EtcdEndpoint struct {
	IP    string
	Name  string
	Nodes string
}

func GenerateEtcdFiles(src string, nodes map[string]string, etcdEndpoints string, version string) error {
	base := path.Join(src, "ssl", "etcd")
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, os.ModePerm)

	caFile := path.Join(src, "ssl", "root", "ca.pem")
	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
	configFile := path.Join(src, "ssl", "root", "ca-config.json")

	for name, ip := range nodes {
		if utils.IsDirExist(path.Join(base, ip)) == false {
			os.MkdirAll(path.Join(base, ip), os.ModePerm)
		}

		var tpl bytes.Buffer
		var err error

		node := EtcdEndpoint{
			IP:    ip,
			Name:  name,
			Nodes: etcdEndpoints,
		}

		sslTp := template.New("etcd-csr")
		sslTp, _ = sslTp.Parse(t.EtcdCATemplate[version])
		sslTp.Execute(&tpl, node)
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

		serviceTp := template.New("etcd-systemd")
		serviceTp, _ = serviceTp.Parse(t.EtcdSystemdTemplate[version])
		serviceTp.Execute(&serviceTpl, node)
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(path.Join(base, ip, "etcd-csr.json"), csrFileBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "etcd-key.pem"), key, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "etcd.csr"), csrBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "etcd.pem"), cert, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "etcd.service"), serviceTpFileBytes, 0700)

		if err != nil {
			return err
		}

	}

	return nil
}

func GenerateAdminCAFiles(src string) error {
	base := path.Join(src, "kubectl")

	caFile := path.Join(src, "ssl", "root", "ca.pem")
	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
	configFile := path.Join(src, "ssl", "root", "ca-config.json")

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

	err = ioutil.WriteFile(path.Join(base, "admin-csr.json"), csrFileBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, "admin-key.pem"), key, 0600)
	err = ioutil.WriteFile(path.Join(base, "admin.csr"), csrBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, "admin.pem"), cert, 0600)

	if err != nil {
		return err
	}

	return nil
}

func GenerateFlanneldFiles(src string, nodes map[string]string, etcdEndpoints string, version string) error {
	base := path.Join(src, "ssl", "flanneld")
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, os.ModePerm)

	caFile := path.Join(src, "ssl", "root", "ca.pem")
	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
	configFile := path.Join(src, "ssl", "root", "ca-config.json")

	for name, ip := range nodes {
		if utils.IsDirExist(path.Join(base, ip)) == false {
			os.MkdirAll(path.Join(base, ip), os.ModePerm)
		}

		node := EtcdEndpoint{
			IP:    ip,
			Name:  name,
			Nodes: etcdEndpoints,
		}

		var tpl bytes.Buffer
		var err error

		sslTp := template.New("flanneld-csr")
		sslTp, _ = sslTp.Parse(t.FlanneldCATemplate[version])
		sslTp.Execute(&tpl, node)
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

		serviceTp := template.New("flanneld-systemd")
		serviceTp, _ = serviceTp.Parse(t.FlanneldSystemdTemplate[version])
		serviceTp.Execute(&serviceTpl, node)
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(path.Join(base, ip, "flanneld-csr.json"), csrFileBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "flanneld-key.pem"), key, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "flanneld.csr"), csrBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "flanneld.pem"), cert, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "flanneld.service"), serviceTpFileBytes, 0700)

		if err != nil {
			return err
		}

	}

	return nil
}

func UploadCARootFiles(src string, files map[string]string, ip string) error {
	key := path.Join(src, "ssh", "id_rsa")

	initCmd := []string{
		"mkdir -p /etc/kubernetes/ssl",
		"mkdir -p /etc/etcd/ssl",
		"mkdir -p /var/lib/etcd",
		"systemctl stop ufw",
		"systemctl disable ufw",
		"apt-get install -y htop denyhosts python-pip",
		"pip install --upgrade pip",
		"pip install glances",
	}

	if err := utils.SSHCommand("root", key, ip, 22, strings.Join(initCmd, " && "), os.Stdout, os.Stderr); err != nil {
		return err
	}

	for _, f := range files {
		if err := utils.SSHScp("root", key, ip, 22, f, path.Join("/etc/kubernetes/ssl", path.Base(f)), os.Stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

func UploadEtcdCAFiles(src string, nodes map[string]string) error {
	base := path.Join(src, "ssl", "etcd")
	key := path.Join(src, "ssh", "id_rsa")

	if utils.IsDirExist(base) == false {
		return fmt.Errorf("Locate etcd folders %s error.", base)
	}

	for _, ip := range nodes {

		var err error

		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "etcd-csr.json"), "/etc/etcd/ssl/etcd-csr.json", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "etcd-key.pem"), "/etc/etcd/ssl/etcd-key.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "etcd.csr"), "/etc/etcd/ssl/etcd.csr", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "etcd.pem"), "/etc/etcd/ssl/etcd.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "etcd.service"), "/etc/systemd/system/etcd.service", os.Stdout, os.Stderr)

		if err != nil {
			return err
		}
	}

	return nil
}

func StartEtcdCluster(key string, nodes map[string]string) error {
	cmd := "systemctl daemon-reload && systemctl enable etcd && systemctl start --no-block etcd"

	for _, ip := range nodes {
		utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr)
	}

	return nil
}

func UploadFlanneldCAFiles(src string, nodes map[string]string) error {
	base := path.Join(src, "ssl", "flanneld")
	key := path.Join(src, "ssh", "id_rsa")

	if utils.IsDirExist(base) == false {
		return fmt.Errorf("Locate flanneld folders %s error.", base)
	}

	for _, ip := range nodes {

		var err error

		initCmd := []string{
			"mkdir -p /etc/flanneld/ssl",
		}

		err = utils.SSHCommand("root", key, ip, 22, initCmd[0], os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "flanneld-csr.json"), "/etc/flanneld/ssl/flanneld-csr.json", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "flanneld-key.pem"), "/etc/flanneld/ssl/flanneld-key.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "flanneld.csr"), "/etc/flanneld/ssl/flanneld.csr", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "flanneld.pem"), "/etc/flanneld/ssl/flanneld.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "flanneld.service"), "/etc/systemd/system/flanneld.service", os.Stdout, os.Stderr)

		if err != nil {
			return err
		}
	}

	return nil
}

func BeforeFlanneldExecute(key, ip, tplString, etcdEndpoints string) error {
	node := EtcdEndpoint{
		Nodes: etcdEndpoints,
	}

	var tpl bytes.Buffer

	sslTp := template.New("before")
	sslTp, _ = sslTp.Parse(tplString)
	sslTp.Execute(&tpl, node)
	cmd := string(tpl.Bytes()[:])

	utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr)

	return nil
}

func StartFlanneldCluster(key string, nodes map[string]string) error {
	cmd := "systemctl daemon-reload && systemctl enable flanneld && systemctl start --no-block flanneld"

	for _, ip := range nodes {
		utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr)
	}

	return nil
}

func GenerateDockerFiles(src string, nodes map[string]string, version string) error {
	base := path.Join(src, "ssl", "docker")
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, os.ModePerm)

	for _, ip := range nodes {
		if utils.IsDirExist(path.Join(base, ip)) == false {
			os.MkdirAll(path.Join(base, ip), os.ModePerm)
		}

		var serviceTpl bytes.Buffer
		var err error

		serviceTp := template.New("docker-systemd")
		serviceTp, err = serviceTp.Parse(t.DockerSystemdTemplate[version])
		serviceTp.Execute(&serviceTpl, nil)
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(path.Join(base, ip, "docker.service"), serviceTpFileBytes, 0700)

		if err != nil {
			return err
		}

	}

	return nil
}

func UploadDockerCAFiles(src string, nodes map[string]string) error {
	base := path.Join(src, "ssl", "docker")
	key := path.Join(src, "ssh", "id_rsa")

	if utils.IsDirExist(base) == false {
		return fmt.Errorf("Locate docker folders %s error.", base)
	}

	for _, ip := range nodes {

		var err error

		initCmd := []string{
			"apt-get update",
			"apt-get dist-upgrade",
			"apt-get -y autoremove",
			"apt-get install -y bridge-utils aufs-tools cgroupfs-mount libltdl7",
		}

		err = utils.SSHCommand("root", key, ip, 22, strings.Join(initCmd, " && "), os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "docker.service"), "/etc/systemd/system/docker.service", os.Stdout, os.Stderr)

		if err != nil {
			return err
		}
	}

	return nil
}

func BeforeDockerExecute(key, ip, cmd string) error {
	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func StartDockerDaemon(key, ip string) error {
	cmd := "systemctl daemon-reload && systemctl enable docker && systemctl start --no-block docker"

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func AfterDockerExecute(key, ip, cmd string) error {

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func UploadKubeConfigFiles(src string, nodes map[string]string) error {
	key := path.Join(src, "ssh", "id_rsa")
	config := path.Join(src, "kubectl", "config")

	for _, ip := range nodes {
		var err error

		err = utils.SSHCommand("root", key, ip, 22, "mkdir -p /root/.kube", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, config, "/root/.kube/config", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(src, "kubectl", "admin.csr"), "/etc/kubernetes/ssl/admin.csr", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(src, "kubectl", "admin-csr.json"), "/etc/kubernetes/ssl/admin-csr.json", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(src, "kubectl", "admin-key.pem"), "/etc/kubernetes/ssl/admin-key.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(src, "kubectl", "admin.pem"), "/etc/kubernetes/ssl/admin.pem", os.Stdout, os.Stderr)

		if err != nil {
			return err
		}
	}

	return nil
}

type KubeMaster struct {
	MasterIP string
	Nodes    string
}

func GenerateKuberAPIServerCAFiles(src string, masterIP, etcdEndpoints string, version string) error {
	base := path.Join(src, "ssl", "kubernetes")
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	os.MkdirAll(base, os.ModePerm)

	caFile := path.Join(src, "ssl", "root", "ca.pem")
	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
	configFile := path.Join(src, "ssl", "root", "ca-config.json")

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

	err = ioutil.WriteFile(path.Join(base, "kubernetes-csr.json"), csrFileBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, "kubernetes-key.pem"), key, 0600)
	err = ioutil.WriteFile(path.Join(base, "kubernetes.csr"), csrBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, "kubernetes.pem"), cert, 0600)
	err = ioutil.WriteFile(path.Join(base, "kube-apiserver.service"), serviceTpFileBytes, 0700)

	if err != nil {
		return err
	}

	return nil
}

func UploadKubeAPIServerCAFiles(src, ip string) error {
	base := path.Join(src, "ssl", "kubernetes")
	key := path.Join(src, "ssh", "id_rsa")

	var err error

	err = utils.SSHScp("root", key, ip, 22, path.Join(base, "kubernetes-csr.json"), "/etc/kubernetes/ssl/kubernetes-csr.json", os.Stdout, os.Stderr)
	err = utils.SSHScp("root", key, ip, 22, path.Join(base, "kubernetes-key.pem"), "/etc/kubernetes/ssl/kubernetes-key.pem", os.Stdout, os.Stderr)
	err = utils.SSHScp("root", key, ip, 22, path.Join(base, "kubernetes.csr"), "/etc/kubernetes/ssl/kubernetes.csr", os.Stdout, os.Stderr)
	err = utils.SSHScp("root", key, ip, 22, path.Join(base, "kubernetes.pem"), "/etc/kubernetes/ssl/kubernetes.pem", os.Stdout, os.Stderr)
	err = utils.SSHScp("root", key, ip, 22, path.Join(base, "kube-apiserver.service"), "/etc/systemd/system/kube-apiserver.service", os.Stdout, os.Stderr)

	if err != nil {
		return err
	}

	return nil
}

func StartKubeAPIServer(src, ip string) error {
	key := path.Join(src, "ssh", "id_rsa")

	cmd := "systemctl daemon-reload && systemctl enable kube-apiserver && systemctl start --no-block kube-apiserver"

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func GenerateKuberControllerManagerFiles(src string, masterIP, etcdEndpoints string, version string) error {
	base := path.Join(src, "ssl", "kubernetes")

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

func UploadKuberControllerFiles(src, ip string) error {
	base := path.Join(src, "ssl", "kubernetes")
	key := path.Join(src, "ssh", "id_rsa")

	if err := utils.SSHScp("root", key, ip, 22, path.Join(base, "kube-controller-manager.service"), "/etc/systemd/system/kube-controller-manager.service", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func StartKuberController(src, ip string) error {
	key := path.Join(src, "ssh", "id_rsa")

	cmd := "systemctl daemon-reload && systemctl enable kube-controller-manager && systemctl start --no-block kube-controller-manager"

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func GenerateKuberSchedulerManagerFiles(src string, masterIP, etcdEndpoints string, version string) error {
	base := path.Join(src, "ssl", "kubernetes")

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

func UploadKuberSchedulerManagerFiles(src, ip string) error {
	base := path.Join(src, "ssl", "kubernetes")
	key := path.Join(src, "ssh", "id_rsa")

	if err := utils.SSHScp("root", key, ip, 22, path.Join(base, "kube-scheduler.service"), "/etc/systemd/system/kube-scheduler.service", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func StartKuberSchedulerManager(src, ip string) error {
	key := path.Join(src, "ssh", "id_rsa")

	cmd := "systemctl daemon-reload && systemctl enable kube-scheduler && systemctl start --no-block kube-scheduler"

	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func UploadBootstrapFile(src string, nodes map[string]string) error {
	config := path.Join(src, "kubectl", "bootstrap.kubeconfig")
	key := path.Join(src, "ssh", "id_rsa")

	for _, ip := range nodes {
		if err := utils.SSHScp("root", key, ip, 22, config, "/etc/kubernetes/bootstrap.kubeconfig", os.Stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

func GenerateKubeletSystemdFile(src string, nodes map[string]string, version string) error {
	for _, ip := range nodes {
		kubeNode := map[string]string{
			"IP": ip,
		}

		base := path.Join(src, "ssl", "kubernetes", ip)
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

func UploadKubeletFile(src string, nodes map[string]string) error {
	key := path.Join(src, "ssh", "id_rsa")

	for _, ip := range nodes {
		file := path.Join(src, "ssl", "kubernetes", ip, "kubelet.service")

		if err := utils.SSHCommand("root", key, ip, 22, "mkdir -p /var/lib/kubelet", os.Stdout, os.Stderr); err != nil {
			return err
		}
		if err := utils.SSHScp("root", key, ip, 22, file, "/etc/systemd/system/kubelet.service", os.Stdout, os.Stderr); err != nil {
			return err
		}

	}

	return nil
}

func SetKubeletClusterrolebinding(src, ip string) error {
	key := path.Join(src, "ssh", "id_rsa")

	if err := utils.SSHCommand("root", key, ip, 22, "kubectl create clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func KubeletCertificateApprove(src, ip string) error {
	key := path.Join(src, "ssh", "id_rsa")

	if err := utils.SSHCommand("root", key, ip, 22, "kubectl certificate approve `kubectl get csr -o name`", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func StartKubelet(src string, nodes map[string]string) error {
	key := path.Join(src, "ssh", "id_rsa")

	for _, ip := range nodes {

		cmd := "systemctl daemon-reload && systemctl enable kubelet && systemctl start --no-block kubelet"

		if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

func GenerateTokenFile(src string) error {
	if utils.IsDirExist(path.Join(src, "kubectl")) == false {
		os.MkdirAll(path.Join(src, "kubectl"), os.ModePerm)
	}

	var tokenTpl bytes.Buffer

	tokenTp := template.New("token")
	tokenTp, _ = tokenTp.Parse("{{.Token}},kubelet-bootstrap,10001,\"system:kubelet-bootstrap\"")
	tokenTp.Execute(&tokenTpl, map[string]string{"Token": t.BooststrapToken})
	tokenTpFileBytes := tokenTpl.Bytes()

	if err := ioutil.WriteFile(path.Join(src, "kubectl", "token.csv"), tokenTpFileBytes, 0700); err != nil {
		return err
	}

	return nil
}

func UploadTokenFiles(src, ip string) error {
	file := path.Join(src, "kubectl", "token.csv")
	key := path.Join(src, "ssh", "id_rsa")

	if err := utils.SSHScp("root", key, ip, 22, file, "/etc/kubernetes/token.csv", os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func GenerateKubeProxyFiles(src string, nodes map[string]string, version string) error {
	base := path.Join(src, "ssl", "kubernetes")

	caFile := path.Join(src, "ssl", "root", "ca.pem")
	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
	configFile := path.Join(src, "ssl", "root", "ca-config.json")

	for _, ip := range nodes {
		if utils.IsDirExist(path.Join(base, ip)) == false {
			os.MkdirAll(path.Join(base, ip), os.ModePerm)
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
		serviceTp.Execute(&serviceTpl, map[string]string{"IP": ip})
		serviceTpFileBytes := serviceTpl.Bytes()

		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy-csr.json"), csrFileBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy-key.pem"), key, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy.csr"), csrBytes, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy.pem"), cert, 0600)
		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy.service"), serviceTpFileBytes, 0700)

		if err != nil {
			return err
		}

	}

	return nil
}

func UploadKubeProxyFiles(src string, nodes map[string]string) error {
	key := path.Join(src, "ssh", "id_rsa")
	base := path.Join(src, "ssl", "kubernetes")

	for _, ip := range nodes {
		var err error

		err = utils.SSHCommand("root", key, ip, 22, "mkdir -p /var/lib/kube-proxy", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "kube-proxy-csr.json"), "/etc/kubernetes/ssl/kube-proxy-csr.json", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "kube-proxy-key.pem"), "/etc/kubernetes/ssl/kube-proxy-key.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "kube-proxy.csr"), "/etc/kubernetes/ssl/kube-proxy.csr", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "kube-proxy.pem"), "/etc/kubernetes/ssl/kube-proxy.pem", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "kube-proxy.service"), "/etc/systemd/system/kube-proxy.service", os.Stdout, os.Stderr)
		err = utils.SSHScp("root", key, ip, 22, path.Join(base, ip, "kube-proxy.kubeconfig"), "/etc/kubernetes/kube-proxy.kubeconfig", os.Stdout, os.Stderr)

		if err != nil {
			return err
		}

	}

	return nil
}

func StartKubeProxy(src string, nodes map[string]string) error {
	key := path.Join(src, "ssh", "id_rsa")

	for _, ip := range nodes {
		cmd := "systemctl daemon-reload && systemctl enable kube-proxy && systemctl start --no-block kube-proxy"

		if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
			return err
		}

	}

	return nil
}
