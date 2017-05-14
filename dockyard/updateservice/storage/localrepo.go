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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/containerops/configure"

	"github.com/Huawei/dockyard/updateservice/km"
	"github.com/Huawei/dockyard/utils"
)

var (
	urlRegexp = regexp.MustCompile(`^(.+)/(.+)/(.+)$`)
	// ErrorInvalidLocalLocalRepo occurs when a local url is invalid
	ErrorInvalidLocalLocalRepo = errors.New("Invalid local url")
	// ErrorEmptyLocalRepo occurs when a repo is empty
	ErrorEmptyLocalRepo = errors.New("LocalRepo is empty")
	// ErrorAppNotExist occurs when an app is not exist
	ErrorAppNotExist = errors.New("App is not exist")
)

const (
	defaultMeta      = "meta.json"
	defaultMetaSign  = "meta.sig"
	defaultKeyDir    = "key"
	defaultPubKey    = "pub_key.pem"
	defaultTargetDir = "target"
)

// LocalRepo reprensents a local repository
// if LocalRepo is {
//	Proto: "protoA/VersionB",
//	Path: "/data",
//      Namespace: "containerops",
//      LocalRepository: "official",
//      }
//    add assume there are 'osX/archY/appA' and 'osX/archY/appB' files.
// The local structure will be:
// /data
//   |_ protoA
//       |_ VersionB
//            |_ containerops
//                |_ official
//                |    |_ meta.json
//                |    |_ meta.sig
//                |    |
//                |    |_ target
//                |         |_ hashOfappA
//                |         |_ hashOfappB
//                |
//                |_ pub_key.pem
type LocalRepo struct {
	Proto           string
	Path            string
	Namespace       string
	LocalRepository string

	kmURL string
}

// NewLocalRepo gets repo by a proto, a local storage path and a url
// nr : "namespace/repository" or just 'namespace'
func NewLocalRepo(path string, proto string, nr string) (LocalRepo, error) {
	repo := LocalRepo{Proto: proto, Path: path}
	parts := strings.Split(nr, "/")
	if nr != "" && len(parts) == 1 {
		repo.Namespace = nr
	} else if len(parts) == 2 {
		repo.Namespace = parts[0]
		repo.LocalRepository = parts[1]
	} else {
		return LocalRepo{}, ErrorInvalidLocalLocalRepo
	}

	// create top dir if not exist
	topDir := repo.GetTopDir()
	if !utils.IsDirExist(topDir) {
		if err := os.MkdirAll(topDir, 0777); err != nil {
			return LocalRepo{}, err
		}
	}

	repo.kmURL = ""
	return repo, nil
}

// NewLocalRepoWithKM gets repo by a proto, a local storage path, a url and
//   a keymanager url
// nr : "namespace/repository"
// kmURL: nil means using the km repository defined in configuration
func NewLocalRepoWithKM(path string, proto, nr string, kmURL string) (LocalRepo, error) {
	if kmURL == "" {
		kmURL = configure.GetString("updateserver.keymanager")
	}

	repo, err := NewLocalRepo(path, proto, nr)
	if err != nil {
		return LocalRepo{}, err
	} else if kmURL == "" {
		return repo, err
	}

	err = repo.SetKM(kmURL)
	if err != nil {
		return LocalRepo{}, err
	}

	return repo, nil
}

// SetKM sets the keymanager
func (r *LocalRepo) SetKM(kmURL string) error {
	// pull the public key
	k, err := km.NewKeyManager(kmURL)
	if err != nil {
		return err
	}

	data, err := k.GetPublicKey(r.Proto, r.Namespace)
	if err != nil {
		return err
	}

	keyfile := r.GetPublicKeyFile()
	if !utils.IsDirExist(filepath.Dir(keyfile)) {
		if err := os.MkdirAll(filepath.Dir(keyfile), 0755); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(keyfile, data, 0644); err != nil {
		return err
	}

	r.kmURL = kmURL
	return nil
}

// GetTopDir gets the top directory of a repository
func (r LocalRepo) GetTopDir() string {
	return filepath.Join(r.Path, r.Proto, r.Namespace, r.LocalRepository)
}

// GetMetaFile gets the meta data file url of repository
func (r LocalRepo) GetMetaFile() string {
	return filepath.Join(r.GetTopDir(), defaultMeta)
}

// GetMetaSignFile gets the meta signature file url of repository
func (r LocalRepo) GetMetaSignFile() string {
	return filepath.Join(r.GetTopDir(), defaultMetaSign)
}

// GetPublicKeyFile gets the public key file url of repository
func (r LocalRepo) GetPublicKeyFile() string {
	return filepath.Join(r.Path, r.Proto, r.Namespace, defaultPubKey)
}

// GetMeta gets the meta data of a repository
func (r LocalRepo) GetMeta() ([]byte, error) {
	metaFile := r.GetMetaFile()
	if !utils.IsFileExist(metaFile) {
		return nil, ErrorEmptyLocalRepo
	}

	return ioutil.ReadFile(metaFile)
}

// List lists the applications inside a repository
func (r LocalRepo) List() ([]string, error) {
	data, err := r.GetMeta()
	if err != nil {
		return nil, err
	}

	var meta utils.Meta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		//This may happend in migration, meta struct changes.
		return nil, nil
	}

	var files []string
	for _, m := range meta.Items {
		files = append(files, m.Name)
	}

	return files, nil
}

// Get gets the data of an application
func (r LocalRepo) Get(name string) ([]byte, error) {
	var meta utils.Meta
	metaFile := r.GetMetaFile()
	if !utils.IsFileExist(metaFile) {
		return nil, ErrorEmptyLocalRepo
	}
	data, err := ioutil.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &meta)
	if err != nil {
		//This may happend in migration, meta struct changes.
		os.Remove(metaFile)
		return nil, nil
	}

	for i := range meta.Items {
		if meta.Items[i].Name != name {
			continue
		}

		dataFile := filepath.Join(r.GetTopDir(), defaultTargetDir, meta.Items[i].Hash)
		return ioutil.ReadFile(dataFile)
	}

	return nil, ErrorAppNotExist
}

// Put adds an application to a repository
func (r LocalRepo) Put(name string, content []byte, method utils.EncryptMethod) (string, error) {
	topDir := r.GetTopDir()
	if !utils.IsDirExist(topDir) {
		if err := os.MkdirAll(topDir, 0777); err != nil {
			return "", err
		}
	}

	var meta utils.Meta
	metaFile := r.GetMetaFile()
	if utils.IsFileExist(metaFile) {
		data, err := ioutil.ReadFile(metaFile)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(data, &meta)
		if err != nil {
			//This may happend in migration, meta struct changes.
			os.Remove(metaFile)
		}

	}

	encryptContent, err := r.encrypt(method, content)
	if err != nil {
		return "", err
	}

	item := utils.GenerateMetaItem(name, encryptContent)
	item.SetEncryption(method)

	// Using the 'hash' value to rename the original file
	dataFileName := item.GetHash()
	dataFile := filepath.Join(topDir, defaultTargetDir, dataFileName)
	if !utils.IsDirExist(filepath.Dir(dataFile)) {
		if err := os.MkdirAll(filepath.Dir(dataFile), 0777); err != nil {
			return "", err
		}
	}

	// write data
	if err := ioutil.WriteFile(dataFile, encryptContent, 0644); err != nil {
		return "", err
	}

	// get meta content
	exist := false
	for i := range meta.Items {
		if meta.Items[i].Name == name {
			meta.Items[i] = item
			exist = true
		}
	}
	if !exist {
		meta.Items = append(meta.Items, item)
	}

	// write meta data
	err = r.saveMeta(meta)
	if err != nil {
		os.Remove(dataFile)
		return "", err
	}

	return dataFile, nil
}

func (r LocalRepo) encrypt(method utils.EncryptMethod, content []byte) ([]byte, error) {
	switch method {
	case utils.EncryptRSA:
		pubBytes, err := ioutil.ReadFile(r.GetPublicKeyFile())
		if err != nil {
			return nil, err
		}
		return utils.RSAEncrypt(pubBytes, content)
	default:
		return content, nil
	}
}

func (r LocalRepo) saveMeta(meta utils.Meta) error {
	meta.Updated = time.Now()
	metaContent, _ := json.Marshal(meta)
	metaFile := r.GetMetaFile()
	err := ioutil.WriteFile(metaFile, metaContent, 0644)
	if err != nil {
		return err
	}

	// write sign file
	err = r.saveSign(metaContent)
	if err != nil {
		os.Remove(metaFile)
		return err
	}

	return nil
}

func (r LocalRepo) saveSign(metaContent []byte) error {
	if r.kmURL == "" {
		return nil
	}

	k, _ := km.NewKeyManager(r.kmURL)
	signContent, _ := k.Sign(r.Proto, r.Namespace, metaContent)
	signFile := r.GetMetaSignFile()
	if err := ioutil.WriteFile(signFile, signContent, 0644); err != nil {
		return err
	}

	return nil
}

// Delete removes an application from a repository
func (r LocalRepo) Delete(name string) error {
	var meta utils.Meta
	metaFile := r.GetMetaFile()
	if !utils.IsFileExist(metaFile) {
		return ErrorEmptyLocalRepo
	}
	data, err := ioutil.ReadFile(metaFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &meta)
	if err != nil {
		//This may happend in migration, meta struct changes.
		os.Remove(metaFile)
		return nil
	}

	for i := range meta.Items {
		if meta.Items[i].Name != name {
			continue
		}

		dataFile := filepath.Join(r.GetTopDir(), defaultTargetDir, meta.Items[i].Hash)
		if err := os.Remove(dataFile); err != nil {
			return err
		}

		meta.Items = append(meta.Items[:i], meta.Items[i+1:]...)
		if err := r.saveMeta(meta); err != nil {
			return err
		}
		return nil
	}

	return ErrorAppNotExist
}
