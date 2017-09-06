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

package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"

	"github.com/Huawei/containerops/common/utils"
	t "github.com/Huawei/containerops/singular/module/template"
)

//GenerateCARootFiles generate root files from template.
func GenerateCARootFiles(src string) (map[string]string, error) {
	var caConfigTpl, caCsrTpl bytes.Buffer
	var err error

	//mkdir for ca root files
	base := path.Join(src, CAFilesFolder, CARootFilesFolder)
	if utils.IsDirExist(base) == false {
		os.MkdirAll(base, os.ModePerm)
	}

	files := map[string]string{
		CARootConfigFile:    path.Join(base, CARootConfigFile),
		CARootCSRConfigFile: path.Join(base, CARootCSRConfigFile),
		CARootPemFile:       path.Join(base, CARootPemFile),
		CARootCSRFile:       path.Join(base, CARootCSRFile),
		CARootKeyFile:       path.Join(base, CARootKeyFile),
	}

	//Remove exist ca files
	for _, value := range files {
		if utils.IsDirExist(value) == true {
			err = os.Remove(value)
			if err != nil {
				return files, err
			}
		}
	}

	//Generate ca-config.json
	caConfig := template.New("config")
	caConfig, err = caConfig.Parse(t.CARootTemplate[t.CARootConfig])
	caConfig.Execute(&caConfigTpl, nil)
	err = ioutil.WriteFile(files[CARootConfigFile], caConfigTpl.Bytes(), 0600)
	if err != nil {
		return files, err
	}

	//Generate ca-csr.json
	caCsr := template.New("csr")
	caCsr, err = caCsr.Parse(t.CARootTemplate[t.CARootCSR])
	caCsr.Execute(&caCsrTpl, nil)
	err = ioutil.WriteFile(files[CARootCSRConfigFile], caCsrTpl.Bytes(), 0600)
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

	//Generate ca.pem, ca.csr, ca-key.pem files
	var key, csrPEM, cert []byte
	cert, csrPEM, key, err = initca.New(&req)
	err = ioutil.WriteFile(files[CARootPemFile], cert, 0600)
	err = ioutil.WriteFile(files[CARootCSRFile], csrPEM, 0600)
	err = ioutil.WriteFile(files[CARootKeyFile], key, 0600)

	return files, err
}

//UploadCARootFiles upload root ca files to the node.
func UploadCARootFiles(key string, roots map[string]string, ip, user string, stdout io.Writer) error {
	files := []map[string]string{}

	for _, f := range roots {
		file := map[string]string{}
		file["src"] = f
		file["dest"] = path.Join(KubeServerConfig, KubeServerSSL, path.Base(f))
		files = append(files, file)
	}

	if err := DownloadComponent(files, ip, key, user, stdout); err != nil {
		return err
	}

	return nil
}
