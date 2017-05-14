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

package storage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/Huawei/dockyard/utils"
)

const (
	localPrefix = "local"
)

// UpdateServiceStorageLocal is the local file implementation of storage service
type UpdateServiceStorageLocal struct {
	Path string

	kmURL string
}

func init() {
	RegisterStorage(localPrefix, &UpdateServiceStorageLocal{})
}

// Supported checks if a uri is a local path
func (ussl *UpdateServiceStorageLocal) Supported(uri string) bool {
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

// New creates an UpdateServceStorage interface with a local implmentation
func (ussl *UpdateServiceStorageLocal) New(uri string, km string) (UpdateServiceStorage, error) {
	if !ussl.Supported(uri) {
		return nil, fmt.Errorf("invalid url set in StorageLocal.New: %s", uri)
	}

	ussl.Path = uri
	ussl.kmURL = km

	return ussl, nil
}

// String returns 'Path'
func (ussl *UpdateServiceStorageLocal) String() string {
	return ussl.Path
}

// Get the data of an input key. Key is "namespace/repository/os/arch/appname"
func (ussl *UpdateServiceStorageLocal) Get(proto string, key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 5 {
		return nil, fmt.Errorf("invalid key detected in StorageLocal.Get: %s", key)
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, strings.Join(s[:2], "/"), ussl.kmURL)
	if err != nil {
		return nil, err
	}

	return r.Get(strings.Join(s[2:], "/"))
}

// GetMeta gets the metadata of an input key. Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) GetMeta(proto string, key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid key detected in StorageLocal.GetMeta: %s", key)
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	return r.GetMeta()
}

// GetMetaSign gets the meta signature data. Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) GetMetaSign(proto string, key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid key detected in StorageLocal.GetMetaSign")
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	file := r.GetMetaSignFile()
	return ioutil.ReadFile(file)
}

// GetPublicKey gets the public key data. Key is "namespace"
func (ussl *UpdateServiceStorageLocal) GetPublicKey(proto string, key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("invalid key detected in StorageLocal.GetPublicKey")
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	file := r.GetPublicKeyFile()
	return ioutil.ReadFile(file)
}

// Put adds a file with a key. Key is "namespace/repository/os/arch/appname"
func (ussl *UpdateServiceStorageLocal) Put(proto string, key string, content []byte, method utils.EncryptMethod) (string, error) {
	s := strings.Split(key, "/")
	if len(s) != 5 {
		return "", errors.New("invalid key detected in StorageLocal.Put")
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, strings.Join(s[:2], "/"), ussl.kmURL)
	if err != nil {
		return "", err
	}

	return r.Put(strings.Join(s[2:], "/"), content, method)
}

// Delete removes a file by a key. Key is "namespace/repositoryi/os/arch/appname"
func (ussl *UpdateServiceStorageLocal) Delete(proto string, key string) error {
	s := strings.Split(key, "/")
	if len(s) != 5 {
		return errors.New("invalid key detected in StorageLocal.Delete")
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, strings.Join(s[:2], "/"), ussl.kmURL)
	if err != nil {
		return err
	}

	return r.Delete(strings.Join(s[2:], "/"))
}

// List lists the content of a key. Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) List(proto string, key string) ([]string, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid key deteced in StorageLocal.List")
	}

	r, err := NewLocalRepoWithKM(ussl.Path, proto, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	return r.List()
}
