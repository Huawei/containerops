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
	"os"
	"path"
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
func DeployFlannelInCluster(d *objects.Deployment, infra *objects.Infra, stdout io.Writer, timestamp bool) error {
	//Get nodes of flanneld
	nodes := []objects.Node{}
	for i := 0; i < infra.Master; i++ {
		nodes = append(nodes, d.Nodes[i])
	}

	//Generate flannel systemd and CA ssl files.
	if files, err := generateFlanneldFiles(d.Config, nodes, d.Outputs["EtcdEndpoints"].(string), infra.Version); err != nil {
		return err
	} else {
		objects.WriteLog(fmt.Sprintf("Flannel CA/Systemd/config files: [%v]", files), stdout, timestamp, d, infra)
		objects.WriteLog(fmt.Sprintf("Upload Flannel CA/Systemd/config files: [%v]", files), stdout, timestamp, d, infra)

		//Upload Flanneld files
		if err := uploadFlanneldFiles(files, d.Tools.SSH.Private, nodes, stdout, timestamp); err != nil {
			return err
		}
	}

	for i, c := range infra.Components {
		if err := d.DownloadBinaryFile(c.Binary, c.URL, nodes, stdout, timestamp); err != nil {
			return err
		}

		if c.Before != "" && i == 0 {
			if err := beforeFlanneldExecute(d.Tools.SSH.Private, d.Outputs[fmt.Sprintf("NODE_%d", i)].(string), c.Before, d.Outputs["EtcdEndpoints"].(string)); err != nil {
				return err
			}
		}
	}

	if err := startFlanneldInCluster(d.Tools.SSH.Private, nodes, stdout, timestamp); err != nil {
		return err
	}

	return nil
}

// Generate Flanneld systemd and CA ssl files.
func generateFlanneldFiles(src string, nodes []objects.Node, etcdEndpoints string, version string) (map[string]map[string]string, error) {
	result := map[string]map[string]string{}

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

	for i, node := range nodes {
		// Mkdir with node ip.
		if utils.IsDirExist(path.Join(base, node.IP)) == false {
			os.MkdirAll(path.Join(base, node.IP), os.ModePerm)
		}

		n := FlanneldEndpoint{
			IP:   node.IP,
			Name: fmt.Sprintf("flanneld-node-%d", i),
			Etcd: etcdEndpoints,
		}

		// generate Flanneld SSL files
		if files, err := generateFlanneldSSLFiles(caFile, caKeyFile, configFile, n, version, base, node.IP); err != nil {
			return result, err
		} else {
			for k, v := range files {
				result[node.IP][k] = v
			}
		}

		// generate Flanneld systemd file
		if files, err := generateFlanneldSystemdFile(n, version, base, node.IP); err != nil {
			return result, err
		} else {
			for k, v := range files {
				result[node.IP][k] = v
			}
		}

	}

	return result, nil
}

// Generate Flanneld SSL files
func generateFlanneldSSLFiles(caFile, caKeyFile, configFile string, node FlanneldEndpoint, version, base, ip string) (map[string]string, error) {
	var tpl bytes.Buffer
	var err error

	files := map[string]string{
		tools.CAFlanneldCSRConfigFile: path.Join(base, ip, tools.CAFlanneldCSRConfigFile),
		tools.CAFlanneldKeyPemFile:    path.Join(base, ip, tools.CAFlanneldKeyPemFile),
		tools.CAFlanneldCSRFile:       path.Join(base, ip, tools.CAFlanneldCSRFile),
		tools.CAFlanneldPemFile:       path.Join(base, ip, tools.CAFlanneldPemFile),
	}

	sslTp := template.New("flanneld-csr")
	sslTp, _ = sslTp.Parse(t.FlanneldCATemplate[version])
	sslTp.Execute(&tpl, node)
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

	err = ioutil.WriteFile(files[tools.CAFlanneldCSRConfigFile], csrFileBytes, 0600)
	err = ioutil.WriteFile(files[tools.CAFlanneldKeyPemFile], key, 0600)
	err = ioutil.WriteFile(files[tools.CAFlanneldCSRFile], csrBytes, 0600)
	err = ioutil.WriteFile(files[tools.CAFlanneldPemFile], cert, 0600)

	if err != nil {
		return files, err
	}

	return files, nil
}

//Generate flanneld systemd file
func generateFlanneldSystemdFile(node FlanneldEndpoint, version, base, ip string) (map[string]string, error) {
	var serviceTpl bytes.Buffer
	files := map[string]string{
		tools.ServiceFlanneldFile: path.Join(base, ip, tools.ServiceFlanneldFile),
	}

	serviceTp := template.New("flanneld-systemd")
	serviceTp, _ = serviceTp.Parse(t.FlanneldSystemdTemplate[version])
	serviceTp.Execute(&serviceTpl, node)
	serviceTpFileBytes := serviceTpl.Bytes()

	if err := ioutil.WriteFile(files[tools.ServiceFlanneldFile], serviceTpFileBytes, 0700); err != nil {
		return files, err
	}

	return files, nil
}

//Upload flanneld SSL files and Systemd file
func uploadFlanneldFiles(files map[string]map[string]string, key string, nodes []objects.Node, stdout io.Writer, timestamp bool) error {
	for _, node := range nodes {
		var err error
		var cmd string

		// Mkdir flanneld ssl folder in server
		initCmd := []string{
			"mkdir -p /etc/flanneld/ssl",
		}

		err = utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, initCmd[0], stdout, os.Stderr)

		// Upload CA SSL files
		cmd, err = tools.DownloadComponent(files[node.IP][tools.CAFlanneldCSRConfigFile], path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldCSRConfigFile), node.IP, key, node.User, stdout)
		objects.WriteLog(
			fmt.Sprintf("upload %s to %s@%s with %s", files[node.IP][tools.CAFlanneldCSRConfigFile], node.IP, path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldCSRConfigFile), cmd),
			stdout, timestamp, &node)

		cmd, err = tools.DownloadComponent(files[node.IP][tools.CAFlanneldKeyPemFile], path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldKeyPemFile), node.IP, key, node.User, stdout)
		objects.WriteLog(
			fmt.Sprintf("upload %s to %s@%s with %s", files[node.IP][tools.CAFlanneldKeyPemFile], node.IP, path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldKeyPemFile), cmd),
			stdout, timestamp, &node)

		cmd, err = tools.DownloadComponent(files[node.IP][tools.CAFlanneldCSRFile], path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldCSRFile), node.IP, key, node.User, stdout)
		objects.WriteLog(
			fmt.Sprintf("upload %s to %s@%s with %s", files[node.IP][tools.CAFlanneldCSRFile], node.IP, path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldCSRFile), cmd),
			stdout, timestamp, &node)

		cmd, err = tools.DownloadComponent(files[node.IP][tools.CAFlanneldPemFile], path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldPemFile), node.IP, key, node.User, stdout)
		objects.WriteLog(
			fmt.Sprintf("upload %s to %s@%s with %s", files[node.IP][tools.CAFlanneldPemFile], node.IP, path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.CAFlanneldPemFile), cmd),
			stdout, timestamp, &node)

		// Upload Systemd file
		cmd, err = tools.DownloadComponent(files[node.IP][tools.ServiceFlanneldFile], path.Join(tools.SystemdServerPath, tools.ServiceFlanneldFile), node.IP, key, node.User, stdout)
		objects.WriteLog(
			fmt.Sprintf("upload %s to %s@%s with %s", files[node.IP][tools.ServiceFlanneldFile], node.IP, path.Join(FlanneldServerConfig, FlanneldServerSSL, tools.ServiceFlanneldFile), cmd),
			stdout, timestamp, &node)

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

func startFlanneldInCluster(key string, nodes []objects.Node, stdout io.Writer, timestamp bool) error {
	cmd := "systemctl daemon-reload && systemctl enable flanneld && systemctl start --no-block flanneld"

	for _, node := range nodes {
		utils.SSHCommand(node.User, key, node.IP, tools.DefaultSSHPort, cmd, stdout, os.Stderr)
		objects.WriteLog(fmt.Sprintf("%s start flanneld in node %s", cmd, node.IP), stdout, timestamp, &node)
	}

	return nil
}
