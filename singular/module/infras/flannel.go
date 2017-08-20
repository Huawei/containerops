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
	//"strings"
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
	// EtcdServerConfig is flanneld config location in the node.
	FlanneldServerConfig = "/etc/flanneld"
	// EtcdServerSSL is the flanneld ssl files folder name in the node.
	// Full path is /etc/flanneld/ssl
	FlanneldServerSSL = "ssl"
)

// FlanneldEndpoint is the etcd node struct.
type FlanneldEndpoint struct {
	IP   string
	Name string
	Etcd string
}

// DeployFlannelInCluster is deploy flannel in cluster.
func DeployFlannelInCluster(d *objects.Deployment, infra *objects.Infra) error {
	//Get nodes of flanneld
	flanneldNodes := map[string]string{}
	for i := 0; i < infra.Master; i++ {
		flanneldNodes[fmt.Sprintf("flanneld-node-%d", i)] = d.Outputs[fmt.Sprintf("NODE_%d", i)].(string)
	}

	//Generate flannel systemd and CA ssl files.
	d.Log(fmt.Sprintf("Generating SSL files and systemd service file for Flanneld."))
	if err := generateFlanneldFiles(d.Config, flanneldNodes, d.Outputs["EtcdEndpoints"].(string), infra.Version); err != nil {
		return err
	}

	//Upload Flanneld files
	d.Log(fmt.Sprintf("Uploading SSL files and systemd service to nodes of Flanneld."))
	if err := uploadFlanneldFiles(d.Config, d.Tools.SSH.Private, flanneldNodes, tools.DefaultSSHUser); err != nil {
		return err
	}

	for i, c := range infra.Components {
		d.Log(fmt.Sprintf("Downloading Flanneld binary files to Nodes."))
		if err := d.DownloadBinaryFile(c.Binary, c.URL, flanneldNodes); err != nil {
			return err
		}

		if c.Before != "" && i == 0 {
			d.Log(fmt.Sprintf("Execute Flanneld before scripts: %s", c.Before))
			if err := beforeFlanneldExecute(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", i)].(string), c.Before, d.Outputs["EtcdEndpoints"].(string)); err != nil {
				return err
			}
		}
	}

	d.Log(fmt.Sprintf("Staring Flanneld Service."))
	if err := startFlanneldInCluster(d.Tools.SSH.Private, flanneldNodes); err != nil {
		return err
	}

	return nil
}

// Generate Flanneld systemd and CA ssl files.
func generateFlanneldFiles(src string, nodes map[string]string, etcdEndpoints string, version string) error {
	// If ca file exist, remove it.
	base := path.Join(src, tools.CAFilesFolder, tools.CAFlanneldFolder)
	if utils.IsDirExist(base) == true {
		os.RemoveAll(base)
	}

	// Mkdir ssl folder
	os.MkdirAll(base, os.ModePerm)

	// If service folder, remove it.
	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceFlanneldFolder)
	if utils.IsDirExist(serviceBase) == true {
		os.RemoveAll(serviceBase)
	}

	// Mkdir ssl folder
	os.MkdirAll(serviceBase, os.ModePerm)

	// CA root files
	caFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootPemFile)
	caKeyFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootKeyFile)
	configFile := path.Join(src, tools.CAFilesFolder, tools.CARootFilesFolder, tools.CARootConfigFile)

	for name, ip := range nodes {
		// Mkdir with node ip.
		if utils.IsDirExist(path.Join(base, ip)) == false {
			os.MkdirAll(path.Join(base, ip), os.ModePerm)
		}

		node := FlanneldEndpoint{
			IP:   ip,
			Name: name,
			Etcd: etcdEndpoints,
		}

		// generate Flanneld SSL files
		if err := generateFlanneldSSLFiles(caFile, caKeyFile, configFile, node, version, base, ip); err != nil {
			return err
		}

		// generate Flanneld systemd file
		if err := generateFlanneldSystemdFile(node, version, base, ip); err != nil {
			return err
		}

	}

	return nil
}

// Generate Flanneld SSL files
func generateFlanneldSSLFiles(caFile, caKeyFile, configFile string, node FlanneldEndpoint, version, base, ip string) error {
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

	err = ioutil.WriteFile(path.Join(base, ip, tools.CAFlanneldCSRConfigFile), csrFileBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, ip, tools.CAFlanneldKeyPemFile), key, 0600)
	err = ioutil.WriteFile(path.Join(base, ip, tools.CAFlanneldCSRFile), csrBytes, 0600)
	err = ioutil.WriteFile(path.Join(base, ip, tools.CAFlanneldPemFile), cert, 0600)

	if err != nil {
		return err
	}

	return nil
}

// Generate flanneld systemd file
func generateFlanneldSystemdFile(node FlanneldEndpoint, version, base, ip string) error {
	var serviceTpl bytes.Buffer

	serviceTp := template.New("flanneld-systemd")
	serviceTp, _ = serviceTp.Parse(t.FlanneldSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, node)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(path.Join(base, ip, tools.ServiceFlanneldFile), serviceTpFileBytes, 0700); err != nil {
		return err
	}

	return nil
}

// Upload flanneld SSL files and Systemd file
func uploadFlanneldFiles(src, key string, nodes map[string]string, user string) error {
	sslBase := path.Join(src, tools.CAFilesFolder, tools.CAFlanneldFolder)
	serviceBase := path.Join(src, tools.ServiceFilesFolder, tools.ServiceFlanneldFolder)

	if utils.IsDirExist(sslBase) == false || utils.IsDirExist(serviceBase) == false {
		return fmt.Errorf("Locate flanneld folders %s or  error", sslBase, serviceBase)
	}

	for _, ip := range nodes {
		var err error

		// Mkdir flanneld ssl folder in server
		initCmd := []string{
			"mkdir -p /etc/flanneld/ssl",
		}

		err = utils.SSHCommand("root", key, ip, 22, initCmd[0], os.Stdout, os.Stderr)

		// Upload CA SSL files
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAFlanneldCSRConfigFile), path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldCSRConfigFile), ip, key, user)
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAFlanneldKeyPemFile), path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldKeyPemFile), ip, key, user)
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAFlanneldCSRFile), path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldCSRFile), ip, key, user)
		err = tools.DownloadComponent(path.Join(sslBase, ip, tools.CAFlanneldPemFile), path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldPemFile), ip, key, user)

		// Upload Systemd file
		err = tools.DownloadComponent(path.Join(serviceBase, ip, tools.ServiceFlanneldFile), path.Join(tools.SytemdServerPath, tools.ServiceFlanneldFile), ip, key, user)

		if err != nil {
			return err
		}
	}

	return nil
}

func beforeFlanneldExecute(key, ip, tplString, etcdEndpoints string) error {
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

func startFlanneldInCluster(key string, nodes map[string]string) error {
	cmd := "systemctl daemon-reload && systemctl enable flanneld && systemctl start --no-block flanneld"

	for _, ip := range nodes {
		utils.SSHCommand("root", key, ip, 22, cmd, os.Stdout, os.Stderr)
	}

	return nil
}
