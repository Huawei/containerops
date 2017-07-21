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
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/Huawei/containerops/common/utils"
	t "github.com/Huawei/containerops/singular/module/template"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
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
	IP string
}

func GenerateEtcdCAFiles(src string, IPs []string) (string, error) {

	//for _, ip := range IPs {
	//
	//}
	return "", nil
}
