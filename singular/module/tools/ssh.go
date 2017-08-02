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

const (
	KeyPath         = "ssh"
	PublicFileName  = "id_rsa.pub"
	PrivateFileName = "id_rsa"
	KeyBits4096     = 4096
)

// GenerateSSHKeyFiles create SSH key in @src , default file name is id_rsa and id_rsa.pub.
// Will remove the SSH key files if there are in @src.
// And in GenerateSSHKeyFiles will return fingerprint.
func GenerateSSHKeyFiles(basePath string) (string, string, string, error) {
	sshPath := path.Join(basePath, KeyPath)
	publicFile := path.Join(sshPath, PublicFileName)
	privateFile := path.Join(sshPath, PrivateFileName)

	privateKey, err := rsa.GenerateKey(rand.Reader, KeyBits4096)
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

// OpenSSHKeyFiles open ssh public and private key files to valid.
// Then return the fingerprint.
func OpenSSHKeyFiles(publicFile, privateFile string) (string, string, string, error) {
	if utils.IsFileExist(privateFile) == false {
		return "", "", "", fmt.Errorf("Private key file not exist")
	}

	var privateByte, publicData []byte
	var err error
	var publicKey ssh.PublicKey
	var fingerprint string

	// Read private key file
	privateByte, err = ioutil.ReadFile(privateFile)
	if err != nil {
		return "", "", "", fmt.Errorf("Read privateFile key file error, %s", err.Error())
	}
	// Decode file byte[]
	blockPrivate, _ := pem.Decode(privateByte)
	rsaPrivate, err := x509.ParsePKCS1PrivateKey(blockPrivate.Bytes)

	if rsaPrivate.Validate() != nil {
		return "", "", "", fmt.Errorf("Valid privateFile key error")
	}

	// Get fingerprint from private key
	pub, err := ssh.NewPublicKey(&rsaPrivate.PublicKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	f := ssh.FingerprintLegacyMD5(pub)

	if publicFile != "" {
		// Read public key file, and return fingerprint.
		publicData, err = ioutil.ReadFile(publicFile)
		if err != nil {
			return "", "", "", fmt.Errorf("Read public file key file error, %s", err.Error())
		}

		publicKey, _, _, _, err = ssh.ParseAuthorizedKey(publicData)
		if err != nil {
			return "", "", "", fmt.Errorf("Parse ssh public file key error")
		}

		fingerprint = ssh.FingerprintLegacyMD5(publicKey)
	} else {
		// No public key file provide, will create one form private key.
		publicSSHKey := ssh.MarshalAuthorizedKey(pub)
		publicFile := path.Join(path.Dir(privateFile), PublicFileName)

		if err := ioutil.WriteFile(publicFile, publicSSHKey, 0600); err != nil {
			return "", "", "", err
		}

		fingerprint = f
	}

	// If different fingerprint form public key and private key, will return error.
	if f != fingerprint {
		return "", "", "", fmt.Errorf("Public and Private key is not pair.")
	}

	return publicFile, privateFile, fingerprint, nil
}
