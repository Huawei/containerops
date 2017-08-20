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

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"path"
//	"strings"
//	"text/template"
//
//	"github.com/cloudflare/cfssl/cli"
//	"github.com/cloudflare/cfssl/cli/genkey"
//	"github.com/cloudflare/cfssl/cli/sign"
//	"github.com/cloudflare/cfssl/csr"
//	"github.com/cloudflare/cfssl/signer"
//
//	"github.com/Huawei/containerops/common/utils"
//	t "github.com/Huawei/containerops/singular/module/template"
//	"github.com/Huawei/containerops/singular/module/tools"
//)
//
//func GenerateAdminCAFiles(src string) error {
//	base := path.Join(src, "kubectl")
//
//	caFile := path.Join(src, "ssl", "root", "ca.pem")
//	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
//	configFile := path.Join(src, "ssl", "root", "ca-config.json")
//
//	var tpl bytes.Buffer
//	var err error
//
//	sslTp := template.New("admin-csr")
//	sslTp, _ = sslTp.Parse(t.AdminCATemplate)
//	sslTp.Execute(&tpl, nil)
//	csrFileBytes := tpl.Bytes()
//
//	req := csr.CertificateRequest{
//		KeyRequest: csr.NewBasicKeyRequest(),
//	}
//
//	err = json.Unmarshal(csrFileBytes, &req)
//	if err != nil {
//		return err
//	}
//
//	var key, csrBytes []byte
//	g := &csr.Generator{Validator: genkey.Validator}
//	csrBytes, key, err = g.ProcessRequest(&req)
//	if err != nil {
//		return err
//	}
//
//	c := cli.Config{
//		CAFile:     caFile,
//		CAKeyFile:  caKeyFile,
//		ConfigFile: configFile,
//		Profile:    "kubernetes",
//		Hostname:   "",
//	}
//
//	s, err := sign.SignerFromConfig(c)
//	if err != nil {
//		return err
//	}
//
//	var cert []byte
//	signReq := signer.SignRequest{
//		Request: string(csrBytes),
//		Hosts:   signer.SplitHosts(c.Hostname),
//		Profile: c.Profile,
//	}
//
//	cert, err = s.Sign(signReq)
//	if err != nil {
//		return err
//	}
//
//	err = ioutil.WriteFile(path.Join(base, "admin-csr.json"), csrFileBytes, 0600)
//	err = ioutil.WriteFile(path.Join(base, "admin-key.pem"), key, 0600)
//	err = ioutil.WriteFile(path.Join(base, "admin.csr"), csrBytes, 0600)
//	err = ioutil.WriteFile(path.Join(base, "admin.pem"), cert, 0600)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//


//
//func UploadKubeConfigFiles(src, key string, nodes map[string]string) error {
//	config := path.Join(src, "kubectl", "config")
//
//	for _, ip := range nodes {
//		var err error
//
//		err = utils.SSHCommand("root", key, ip, 22, "mkdir -p /root/.kube", os.Stdout, os.Stderr)
//		err = tools.DownloadComponent(config, "/root/.kube/config", ip, key)
//		err = tools.DownloadComponent(path.Join(src, "kubectl", "admin.csr"), "/etc/kubernetes/ssl/admin.csr", ip, key)
//		err = tools.DownloadComponent(path.Join(src, "kubectl", "admin-csr.json"), "/etc/kubernetes/ssl/admin-csr.json", ip, key)
//		err = tools.DownloadComponent(path.Join(src, "kubectl", "admin-key.pem"), "/etc/kubernetes/ssl/admin-key.pem", ip, key)
//		err = tools.DownloadComponent(path.Join(src, "kubectl", "admin.pem"), "/etc/kubernetes/ssl/admin.pem", ip, key)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//type KubeMaster struct {
//	MasterIP string
//	Nodes    string
//}
//
//func GenerateKuberAPIServerCAFiles(src string, masterIP, etcdEndpoints string, version string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//	if utils.IsDirExist(base) == true {
//		os.RemoveAll(base)
//	}
//
//	os.MkdirAll(base, os.ModePerm)
//
//	caFile := path.Join(src, "ssl", "root", "ca.pem")
//	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
//	configFile := path.Join(src, "ssl", "root", "ca-config.json")
//
//	master := KubeMaster{
//		MasterIP: masterIP,
//		Nodes:    etcdEndpoints,
//	}
//
//	var tpl bytes.Buffer
//	var err error
//
//	sslTp := template.New("kube-csr")
//	sslTp, _ = sslTp.Parse(t.KubernetesCATemplate[version])
//	sslTp.Execute(&tpl, master)
//	csrFileBytes := tpl.Bytes()
//
//	req := csr.CertificateRequest{
//		KeyRequest: csr.NewBasicKeyRequest(),
//	}
//
//	err = json.Unmarshal(csrFileBytes, &req)
//	if err != nil {
//		return err
//	}
//
//	var key, csrBytes []byte
//	g := &csr.Generator{Validator: genkey.Validator}
//	csrBytes, key, err = g.ProcessRequest(&req)
//	if err != nil {
//		return err
//	}
//
//	c := cli.Config{
//		CAFile:     caFile,
//		CAKeyFile:  caKeyFile,
//		ConfigFile: configFile,
//		Profile:    "kubernetes",
//		Hostname:   "",
//	}
//
//	s, err := sign.SignerFromConfig(c)
//	if err != nil {
//		return err
//	}
//
//	var cert []byte
//	signReq := signer.SignRequest{
//		Request: string(csrBytes),
//		Hosts:   signer.SplitHosts(c.Hostname),
//		Profile: c.Profile,
//	}
//
//	cert, err = s.Sign(signReq)
//	if err != nil {
//		return err
//	}
//
//	var serviceTpl bytes.Buffer
//
//	serviceTp := template.New("kube-systemd")
//	serviceTp, _ = serviceTp.Parse(t.KubernetesAPIServerSystemdTemplate[version])
//	serviceTp.Execute(&serviceTpl, master)
//	serviceTpFileBytes := serviceTpl.Bytes()
//
//	err = ioutil.WriteFile(path.Join(base, "kubernetes-csr.json"), csrFileBytes, 0600)
//	err = ioutil.WriteFile(path.Join(base, "kubernetes-key.pem"), key, 0600)
//	err = ioutil.WriteFile(path.Join(base, "kubernetes.csr"), csrBytes, 0600)
//	err = ioutil.WriteFile(path.Join(base, "kubernetes.pem"), cert, 0600)
//	err = ioutil.WriteFile(path.Join(base, "kube-apiserver.service"), serviceTpFileBytes, 0700)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func UploadKubeAPIServerCAFiles(src, key, ip string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	var err error
//
//	err = tools.DownloadComponent(path.Join(base, "kubernetes-csr.json"), "/etc/kubernetes/ssl/kubernetes-csr.json", ip, key)
//	err = tools.DownloadComponent(path.Join(base, "kubernetes-key.pem"), "/etc/kubernetes/ssl/kubernetes-key.pem", ip, key)
//	err = tools.DownloadComponent(path.Join(base, "kubernetes.csr"), "/etc/kubernetes/ssl/kubernetes.csr", ip, key)
//	err = tools.DownloadComponent(path.Join(base, "kubernetes.pem"), "/etc/kubernetes/ssl/kubernetes.pem", ip, key)
//	err = tools.DownloadComponent(path.Join(base, "kube-apiserver.service"), "/etc/systemd/system/kube-apiserver.service", ip, key)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func StartKubeAPIServer(key, ip string) error {
//	cmd := "systemctl daemon-reload && systemctl enable kube-apiserver && systemctl start --no-block kube-apiserver"
//
//	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GenerateKuberControllerManagerFiles(src string, masterIP, etcdEndpoints string, version string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	master := KubeMaster{
//		MasterIP: masterIP,
//		Nodes:    etcdEndpoints,
//	}
//
//	var serviceTpl bytes.Buffer
//
//	serviceTp := template.New("kube-control-systemd")
//	serviceTp, _ = serviceTp.Parse(t.KubernetesControllerManagerSystemdTemplate[version])
//	serviceTp.Execute(&serviceTpl, master)
//	serviceTpFileBytes := serviceTpl.Bytes()
//
//	if err := ioutil.WriteFile(path.Join(base, "kube-controller-manager.service"), serviceTpFileBytes, 0700); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func UploadKuberControllerFiles(src, key, ip string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	if err := tools.DownloadComponent(path.Join(base, "kube-controller-manager.service"), "/etc/systemd/system/kube-controller-manager.service", ip, key); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func StartKuberController(key, ip string) error {
//	cmd := "systemctl daemon-reload && systemctl enable kube-controller-manager && systemctl start --no-block kube-controller-manager"
//
//	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GenerateKuberSchedulerManagerFiles(src string, masterIP, etcdEndpoints string, version string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	master := KubeMaster{
//		MasterIP: masterIP,
//		Nodes:    etcdEndpoints,
//	}
//
//	var serviceTpl bytes.Buffer
//
//	serviceTp := template.New("kube-scheduler-systemd")
//	serviceTp, _ = serviceTp.Parse(t.KubernetesSchedulerSystemdTemplate[version])
//	serviceTp.Execute(&serviceTpl, master)
//	serviceTpFileBytes := serviceTpl.Bytes()
//
//	if err := ioutil.WriteFile(path.Join(base, "kube-scheduler.service"), serviceTpFileBytes, 0700); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func UploadKuberSchedulerManagerFiles(src, key, ip string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	if err := tools.DownloadComponent(path.Join(base, "kube-scheduler.service"), "/etc/systemd/system/kube-scheduler.service", ip, key); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func StartKuberSchedulerManager(key, ip string) error {
//	cmd := "systemctl daemon-reload && systemctl enable kube-scheduler && systemctl start --no-block kube-scheduler"
//
//	if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func UploadBootstrapFile(src, key string, nodes map[string]string) error {
//	config := path.Join(src, "kubectl", "bootstrap.kubeconfig")
//
//	for _, ip := range nodes {
//		if err := tools.DownloadComponent(config, "/etc/kubernetes/bootstrap.kubeconfig", ip, key); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func GenerateKubeletSystemdFile(src string, nodes map[string]string, version string) error {
//	for _, ip := range nodes {
//		kubeNode := map[string]string{
//			"IP": ip,
//		}
//
//		base := path.Join(src, "ssl", "kubernetes", ip)
//		if utils.IsDirExist(base) == true {
//			os.RemoveAll(base)
//		}
//
//		os.MkdirAll(base, os.ModePerm)
//
//		var serviceTpl bytes.Buffer
//
//		serviceTp := template.New("kubelete-systemd")
//		serviceTp, _ = serviceTp.Parse(t.KubeletSystemdTemplate[version])
//		serviceTp.Execute(&serviceTpl, kubeNode)
//		serviceTpFileBytes := serviceTpl.Bytes()
//
//		if err := ioutil.WriteFile(path.Join(base, "kubelet.service"), serviceTpFileBytes, 0700); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func UploadKubeletFile(src, key string, nodes map[string]string) error {
//	for _, ip := range nodes {
//		file := path.Join(src, "ssl", "kubernetes", ip, "kubelet.service")
//
//		if err := utils.SSHCommand("root", key, ip, 22, "mkdir -p /var/lib/kubelet", os.Stdout, os.Stderr); err != nil {
//			return err
//		}
//
//		if err := tools.DownloadComponent(file, "/etc/systemd/system/kubelet.service", ip, key); err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//}
//
//func SetKubeletClusterrolebinding(key, ip string) error {
//	if err := utils.SSHCommand("root", key, ip, 22, "kubectl create clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap", os.Stdout, os.Stderr); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func KubeletCertificateApprove(key, ip string) error {
//	if err := utils.SSHCommand("root", key, ip, 22, "kubectl certificate approve `kubectl get csr -o name`", os.Stdout, os.Stderr); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func StartKubelet(key string, nodes map[string]string) error {
//	for _, ip := range nodes {
//
//		cmd := "systemctl daemon-reload && systemctl enable kubelet && systemctl start --no-block kubelet"
//
//		if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func GenerateTokenFile(src string) error {
//	if utils.IsDirExist(path.Join(src, "kubectl")) == false {
//		os.MkdirAll(path.Join(src, "kubectl"), os.ModePerm)
//	}
//
//	var tokenTpl bytes.Buffer
//
//	tokenTp := template.New("token")
//	tokenTp, _ = tokenTp.Parse("{{.Token}},kubelet-bootstrap,10001,\"system:kubelet-bootstrap\"")
//	tokenTp.Execute(&tokenTpl, map[string]string{"Token": t.BooststrapToken})
//	tokenTpFileBytes := tokenTpl.Bytes()
//
//	if err := ioutil.WriteFile(path.Join(src, "kubectl", "token.csv"), tokenTpFileBytes, 0700); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func UploadTokenFiles(src, key, ip string) error {
//	file := path.Join(src, "kubectl", "token.csv")
//
//	if err := tools.DownloadComponent(file, "/etc/kubernetes/token.csv", ip, key); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GenerateKubeProxyFiles(src string, nodes map[string]string, version string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	caFile := path.Join(src, "ssl", "root", "ca.pem")
//	caKeyFile := path.Join(src, "ssl", "root", "ca-key.pem")
//	configFile := path.Join(src, "ssl", "root", "ca-config.json")
//
//	for _, ip := range nodes {
//		if utils.IsDirExist(path.Join(base, ip)) == false {
//			os.MkdirAll(path.Join(base, ip), os.ModePerm)
//		}
//
//		var tpl bytes.Buffer
//		var err error
//
//		sslTp := template.New("proxy-csr")
//		sslTp, _ = sslTp.Parse(t.KubeProxyCATemplate[version])
//		sslTp.Execute(&tpl, nil)
//		csrFileBytes := tpl.Bytes()
//
//		req := csr.CertificateRequest{
//			KeyRequest: csr.NewBasicKeyRequest(),
//		}
//
//		err = json.Unmarshal(csrFileBytes, &req)
//		if err != nil {
//			return err
//		}
//
//		var key, csrBytes []byte
//		g := &csr.Generator{Validator: genkey.Validator}
//		csrBytes, key, err = g.ProcessRequest(&req)
//		if err != nil {
//			return err
//		}
//
//		c := cli.Config{
//			CAFile:     caFile,
//			CAKeyFile:  caKeyFile,
//			ConfigFile: configFile,
//			Profile:    "kubernetes",
//			Hostname:   "",
//		}
//
//		s, err := sign.SignerFromConfig(c)
//		if err != nil {
//			return err
//		}
//
//		var cert []byte
//		signReq := signer.SignRequest{
//			Request: string(csrBytes),
//			Hosts:   signer.SplitHosts(c.Hostname),
//			Profile: c.Profile,
//		}
//
//		cert, err = s.Sign(signReq)
//		if err != nil {
//			return err
//		}
//
//		var serviceTpl bytes.Buffer
//
//		serviceTp := template.New("proxy-systemd")
//		serviceTp, _ = serviceTp.Parse(t.KubeProxySystemdTemplate[version])
//		serviceTp.Execute(&serviceTpl, map[string]string{"IP": ip})
//		serviceTpFileBytes := serviceTpl.Bytes()
//
//		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy-csr.json"), csrFileBytes, 0600)
//		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy-key.pem"), key, 0600)
//		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy.csr"), csrBytes, 0600)
//		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy.pem"), cert, 0600)
//		err = ioutil.WriteFile(path.Join(base, ip, "kube-proxy.service"), serviceTpFileBytes, 0700)
//
//		if err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//}
//
//func UploadKubeProxyFiles(src, key string, nodes map[string]string) error {
//	base := path.Join(src, "ssl", "kubernetes")
//
//	for _, ip := range nodes {
//		var err error
//
//		err = utils.SSHCommand("root", key, ip, 22, "mkdir -p /var/lib/kube-proxy", os.Stdout, os.Stderr)
//		err = tools.DownloadComponent(path.Join(base, ip, "kube-proxy-csr.json"), "/etc/kubernetes/ssl/kube-proxy-csr.json", ip, key)
//		err = tools.DownloadComponent(path.Join(base, ip, "kube-proxy-key.pem"), "/etc/kubernetes/ssl/kube-proxy-key.pem", ip, key)
//		err = tools.DownloadComponent(path.Join(base, ip, "kube-proxy.csr"), "/etc/kubernetes/ssl/kube-proxy.csr", ip, key)
//		err = tools.DownloadComponent(path.Join(base, ip, "kube-proxy.pem"), "/etc/kubernetes/ssl/kube-proxy.pem", ip, key)
//		err = tools.DownloadComponent(path.Join(base, ip, "kube-proxy.service"), "/etc/systemd/system/kube-proxy.service", ip, key)
//		err = tools.DownloadComponent(path.Join(base, ip, "kube-proxy.kubeconfig"), "/etc/kubernetes/kube-proxy.kubeconfig", ip, key)
//
//		if err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//}
//
//func StartKubeProxy(key string, nodes map[string]string) error {
//	for _, ip := range nodes {
//		cmd := "systemctl daemon-reload && systemctl enable kube-proxy && systemctl start --no-block kube-proxy"
//
//		if err := utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr); err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//}
