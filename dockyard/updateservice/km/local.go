/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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

package km

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Huawei/dockyard/utils"
)

const (
	localPrefix       = "local"
	defaultPublicKey  = "pub_key.pem"
	defaultPrivateKey = "priv_key.pem"
	defaultBitsSize   = 2048
)

// KeyManagerLocal is the local implementation of a key manager

type KeyManagerLocal struct {
	Path string
}

func init() {
	RegisterKeyManager(localPrefix, &KeyManagerLocal{})
}

// Supported checks if a uri is local path
func (kml *KeyManagerLocal) Supported(uri string) bool {
	if uri == "" {
		return false
	}

	if u, err := url.Parse(uri); err != nil {
		return false
	} else if u.Scheme == "" {
		return true
	}

	return false
}

// New returns a keymanager by a uri
func (kml *KeyManagerLocal) New(uri string) (KeyManager, error) {
	if !kml.Supported(uri) {
		return nil, errors.New("Invalid key manager url, should be local dir")
	}

	kml.Path = uri
	return kml, nil
}

// getKeyDir returns key dir
func (kml *KeyManagerLocal) getKeyDir(proto, namespace string) (string, error) {
	keyDir := filepath.Join(kml.Path, proto, namespace)
	if !isKeyExist(keyDir) {
		err := generateKey(keyDir)
		if err != nil {
			return "", err
		}
	}

	return keyDir, nil
}

// GetPublicKey gets the public key data of a namespace
func (kml *KeyManagerLocal) GetPublicKey(proto string, namespace string) ([]byte, error) {
	keyDir, err := kml.getKeyDir(proto, namespace)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filepath.Join(keyDir, defaultPublicKey))
}

// Sign signs a data of a namespace
func (kml *KeyManagerLocal) Decrypt(proto string, namespace string, data []byte) ([]byte, error) {
	keyDir, err := kml.getKeyDir(proto, namespace)
	if err != nil {
		return nil, err
	}

	privBytes, _ := ioutil.ReadFile(filepath.Join(keyDir, defaultPrivateKey))
	return utils.RSADecrypt(privBytes, data)
}

// Sign signs a data of a namespace
func (kml *KeyManagerLocal) Sign(proto string, namespace string, data []byte) ([]byte, error) {
	keyDir, err := kml.getKeyDir(proto, namespace)
	if err != nil {
		return nil, err
	}

	privBytes, _ := ioutil.ReadFile(filepath.Join(keyDir, defaultPrivateKey))
	return utils.SHA256Sign(privBytes, data)
}

func isKeyExist(keyDir string) bool {
	if !utils.IsFileExist(filepath.Join(keyDir, defaultPrivateKey)) {
		return false
	}

	if !utils.IsFileExist(filepath.Join(keyDir, defaultPublicKey)) {
		return false
	}

	return true
}

func generateKey(keyDir string) error {
	privBytes, pubBytes, err := utils.GenerateRSAKeyPair(defaultBitsSize)
	if err != nil {
		return err
	}

	if !utils.IsDirExist(keyDir) {
		err := os.MkdirAll(keyDir, 0755)
		if err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(filepath.Join(keyDir, defaultPrivateKey), privBytes, 0600); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(keyDir, defaultPublicKey), pubBytes, 0644); err != nil {
		return err
	}

	return nil
}
