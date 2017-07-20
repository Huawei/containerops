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
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	homeDir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module/service"
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

// Deploy Sequence: Preparing SSH Key files -> Preparing VM -> Preparing SSL Key files -> Deploy Etcd
//   -> Deploy flannel -> Deploy k8s Master -> Deploy k8s node -> TODO Deploy other...
func (d *Deployment) Deploy() error {
	// Preparing SSH Keys
	if d.Tools.SSH.Public == "" || d.Tools.SSH.Private == "" {
		if public, private, fingerprint, err := CreateSSHKeyFiles(d.Config); err != nil {
			return err
		} else {
			d.Tools.SSH.Public, d.Tools.SSH.Private, d.Tools.SSH.Fingerprint = public, private, fingerprint
		}
	}

	switch d.Service.Provider {
	case "digitalocean":
		do := new(service.DigitalOcean)
		do.Token = d.Service.Token
		do.Region, do.Size, do.Image = d.Service.Region, d.Service.Size, d.Service.Image

		do.InitClient()

		if err := do.UpdateSSHKey(d.Tools.SSH.Public); err != nil {
			return err
		}

		if err := do.CreateDroplet(d.Nodes, d.Tools.SSH.Fingerprint); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupport service provide: %s", d.Service.Provider)

	}

	return nil
}
