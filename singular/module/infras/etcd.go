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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

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
	// EtcdMinimalNodes is minimal etcd nodes number.
	EtcdMinimalNodes = 2
	// EtcdServerConfig is etcd config location in the node.
	EtcdServerConfig = "/etc/etcd"
	// EtcdServerSSL is the etcd ssl files folder name in the node.
	// Full path is /etc/etcd/ssl
	EtcdServerSSL = "ssl"
)

// EtcdEndpoint is the etcd node struct.
type EtcdEndpoint struct {
	IP    string
	Name  string
	Nodes string
}

// DeployEtcdCluster deploy etcd cluster.
// Notes:
//   1. Only count master nodes in etcd deploy process.
//   2.
func DeployEtcdCluster(d *objects.Deployment, infra *objects.Infra) error {
	infra.Log("Deploying etcd clusters.")

	// Check master node number.
	if infra.Master > len(d.Nodes) {
		return fmt.Errorf("deploy %s nodes more than %d", infra.Name, len(d.Nodes))
	}

	if infra.Master < EtcdMinimalNodes {
		return fmt.Errorf("etcd node no less than %d nodes", EtcdMinimalNodes)
	}

	// Init nodes, endpoints and adminEndpoints parameters.
	etcdNodes := map[string]string{}
	etcdEndpoints, etcdPeerEndpoints := []string{}, []string{}

	// Get nodes from outputs of Deployment to determine etcd cluster nodes.
	// TODO Now just choose the first nodes of list. Should have a algorithm and filers determined.
	for i := 0; i < infra.Master; i++ {
		// Etcd Notes
		etcdNodes[fmt.Sprintf("etcd-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)

		// Etcd endpoints for client
		etcdEndpoints = append(etcdEndpoints,
			fmt.Sprintf("https://%s:2379", d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)))

		// Etcd admin endpoints for peer
		etcdPeerEndpoints = append(etcdPeerEndpoints,
			fmt.Sprintf("%s=https://%s:2380", fmt.Sprintf("etcd-node-%d", i),
				d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)))
	}

	// Deployment output
	d.Output("EtcdEndpoints", strings.Join(etcdEndpoints, ","))
	d.Output("EtcdPeerEndpoints", strings.Join(etcdPeerEndpoints, ","))

	// Infra output
	infra.Output("EtcdEndpoints", strings.Join(etcdEndpoints, ","))
	infra.Output("EtcdPeerEndpoints", strings.Join(etcdPeerEndpoints, ","))

	// d.Log and infra.log
	infra.Log(fmt.Sprintf("Generating Etcd endpoints environment variable [%s], value is\n [%s]", "EtcdEndpoints", strings.Join(etcdEndpoints, ",")))
	infra.Log(fmt.Sprintf("Generating SSL files and systemd service file for Etcd cluster."))

	// Generate Etcd CA Files
	if err := generateEtcdFiles(d.Config, etcdNodes, strings.Join(etcdPeerEndpoints, ","), infra.Version); err != nil {
		return err
	}

	d.Log(fmt.Sprintf("Uploading SSL files to nodes of Etcd Cluster."))
	if err := uploadEtcdFiles(d.Config, d.Tools.SSH.Private, etcdNodes, tools.DefaultSSHUser); err != nil {
		return err
	}

	d.Log(fmt.Sprintf("Downloading Etcd binary files to nodes of Etcd Cluster."))
	for _, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, etcdNodes); err != nil {
			return err
		}
	}

	d.Log(fmt.Sprintf("Staring Etcd Cluster."))
	if err := startEtcdCluster(d.Tools.SSH.Private, etcdNodes); err != nil {
		return err
	}

	return nil
}

// Generate Etcd CA SSL and Systemd service Files
func generateEtcdFiles(src string, nodes map[string]string, etcdEndpoints string, version string) error {
	// If ca file exist, remove it.
	sslBase := path.Join(src, tools.CAFilesFolder, tools.CAEtcdFolder)
	if utils.IsDirExist(sslBase) == true {
		os.RemoveAll(sslBase)
	}

	// Mkdir ssl folder
	os.MkdirAll(sslBase, os.ModePerm)

	// If service folder, remove it.
	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceEtcdFolder)
	if utils.IsDirExist(serviceBase) == true {
		os.RemoveAll(serviceBase)
	}

	// Mkdir ssl folder
	os.MkdirAll(serviceBase, os.ModePerm)

	// CA root files
	caFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	// Loop etcd nodes and generate CA files.
	for name, ip := range nodes {
		// Mkdir with node ip.
		if utils.IsDirExist(path.Join(sslBase, ip)) == false {
			os.MkdirAll(path.Join(sslBase, ip), os.ModePerm)
		}

		if utils.IsDirExist(path.Join(serviceBase, ip)) == false {
			os.MkdirAll(path.Join(serviceBase, ip), os.ModePerm)
		}

		node := EtcdEndpoint{
			IP:    ip,
			Name:  name,
			Nodes: etcdEndpoints,
		}

		// Generate Etcd SSL files
		if err := generateEtcdSSLFiles(caFile, caKeyFile, configFile, node, version, sslBase, ip); err != nil {
			return err
		}

		// Generate Etcd Systemd File
		if err := generateEtcdServiceFile(node, version, serviceBase, ip); err != nil {
			return err
		}
	}

	return nil
}

// generateEtcdSSLFiles generate ssl file with node information.
func generateEtcdSSLFiles(caFile, caKeyFile, configFile string, node EtcdEndpoint, version, base, ip string) error {
	var tpl bytes.Buffer
	var err error

	// Generate csr file
	sslTp := template.New("etcd-csr")
	sslTp, _ = sslTp.Parse(t.EtcdCATemplate[version])
	sslTp.Execute(&tpl, node)
	csrFileBytes := tpl.Bytes()

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}

	// Unmarshal csr to certificate request
	err = json.Unmarshal(csrFileBytes, &req)
	if err != nil {
		return err
	}

	// Generate key file and others.
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

	err = ioutil.WriteFile(path.Join(base, ip, tools.CAEtcdCSRConfigFile), csrFileBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, ip, tools.CAEtcdKeyPemFile), key, 0600)
	err = ioutil.WriteFile(path.Join(base, ip, tools.CAEtcdCSRFile), csrBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, ip, tools.CAEtcdPemFile), cert, 0600)

	if err != nil {
		return err
	}

	return nil
}

// generateEtcdServiceFile generate systemd service file
func generateEtcdServiceFile(node EtcdEndpoint, version, base, ip string) error {
	var serviceTpl bytes.Buffer

	serviceTp := template.New("etcd-systemd")
	serviceTp, _ = serviceTp.Parse(t.EtcdSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, node)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(path.Join(base, ip, tools.ServiceEtcdFile), serviceTpFileBytes, 0700); err != nil {
		return err
	}

	return nil
}

// upload Etcd SSL files and systemd file to nodes
func uploadEtcdFiles(src, key string, nodes map[string]string, user string) error {
	sslBase := path.Join(src, tools.CAFilesFolder, tools.CAEtcdFolder)
	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceEtcdFolder)

	if utils.IsDirExist(sslBase) == false || utils.IsDirExist(serviceBase) {
		return fmt.Errorf("Locate etcd folders %s or %s error", sslBase, serviceBase)
	}

	for _, ip := range nodes {
		var err error

		// Upload SSL files
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAEtcdCSRConfigFile), path.Join(EtcdServerConfig, EtcdServerSSL, tools.CAEtcdCSRConfigFile), ip, key, user)
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAEtcdKeyPemFile), path.Join(EtcdServerConfig, EtcdServerSSL, tools.CAEtcdKeyPemFile), ip, key, user)
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAEtcdCSRFile), path.Join(EtcdServerConfig, EtcdServerSSL, tools.CAEtcdCSRFile), ip, key, user)
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAEtcdPemFile), path.Join(EtcdServerConfig, EtcdServerSSL, tools.CAEtcdPemFile), ip, key, user)

		// Upload Service file
		err = tools.DownloadComponent(path.Join(serviceBase, ip, tools.ServiceEtcdFile), path.Join(tools.SytemdServerPath, tools.ServiceEtcdFile), ip, key, user)

		if err != nil {
			return err
		}
	}

	return nil
}

func startEtcdCluster(key string, nodes map[string]string) error {
	cmd := "systemctl daemon-reload && systemctl enable etcd && systemctl start --no-block etcd"

	for _, ip := range nodes {
		utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr)
	}

	return nil
}
