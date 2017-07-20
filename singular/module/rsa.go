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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/ssh"

	"github.com/Huawei/containerops/common/utils"
)

//
func CreateSSHKeyFiles(config string) (string, string, string, error) {

	sshPath := path.Join(config, "ssh")
	publicFile := path.Join(sshPath, "id_rsa.pub")
	privateFile := path.Join(sshPath, "id_rsa")

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", "", err
	}

	// generate private key
	var private bytes.Buffer
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return "", "", "", err
	}

	// generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	public := ssh.MarshalAuthorizedKey(pub)
	fingerprint := ssh.FingerprintLegacyMD5(pub)

	// Remove exist public and private file
	if utils.IsFileExist(privateFile) == true {
		if err := os.Remove(privateFile); err != nil {
			return "", "", "", err
		}
	}

	if utils.IsFileExist(publicFile) == true {
		if err := os.Remove(publicFile); err != nil {
			return "", "", "", err
		}
	}

	// Save public key and private key file.
	if utils.IsDirExist(sshPath) == false {
		if err := os.MkdirAll(sshPath, os.ModePerm); err != nil {
			return "", "", "", err
		}
	}

	if err := ioutil.WriteFile(privateFile, private.Bytes(), 0400); err != nil {
		return "", "", "", err
	}
	if err := ioutil.WriteFile(publicFile, public, 0600); err != nil {
		return "", "", "", err
	}

	return publicFile, privateFile, fingerprint, nil
}
