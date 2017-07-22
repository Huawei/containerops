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
	"log"
	"net"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"golang.org/x/crypto/ssh"

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
		err = ioutil.WriteFile(path.Join(base, ip, "etcd.service"), serviceTpFileBytes, 0600)

		if err != nil {
			return err
		}

	}

	return nil
}

func UploadEtcdCAFiles(src string, nodes map[string]string) error {
	base := path.Join(src, "ssl", "etcd")
	if utils.IsDirExist(base) == false {
		return fmt.Errorf("Locate etcd folders %s error.", base)
	}

	for _, ip := range nodes {

		time.Sleep(10 * time.Second)

		initCmd := []string{
			"mkdir -p /etc/kubernetes/ssl",
			"mkdir -p /etc/etcd/ssl",
			"mkdir -p /var/lib/etcd",
			"systemctl stop ufw",
			"systemctl disable ufw",
			"apt-get update",
			"apt-get dist-upgrade",
			"apt-get install -y bridge-utils htop denyhosts python-pip aufs-tools cgroupfs-mount libltdl7",
			"pip install --upgrade pip",
			"pip install glances",
		}

		session, err := connect("root", "", ip, 22)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Run(strings.Join(initCmd, " && "))

	}

	return nil
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func connect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/meaglith/.containerops/singular/ssh/id_rsa"),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 0,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}
