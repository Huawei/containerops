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

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/Huawei/dockyard/utils"
)

const (
	appV1Protocal = "appv1"
	appV1Restful  = "app/v1"
)

var (
	appv1Regexp = regexp.MustCompile(`^(.+)://(.+)/(.+)/(.+)$`)
)

// UpdateClientAppV1Repo represents the 'appV1' repo
type UpdateClientAppV1Repo struct {
	Runmode   string
	Site      string
	Namespace string
	Repo      string
}

func init() {
	RegisterRepo(appV1Protocal, &UpdateClientAppV1Repo{})
}

// Supported checks if a protocal is "appv1"
func (ap *UpdateClientAppV1Repo) Supported(protocal string) bool {
	return protocal == appV1Protocal
}

// New parses 'http://containerops/dockyard.me/containerops/dockyard/offical' and get
//	Site:       "containerops/dockyard.me"
//      Namespace:  "containerops/dockyard"
//      Repo:       "offical"
func (ap *UpdateClientAppV1Repo) New(url string) (UpdateClientRepo, error) {
	parts := appv1Regexp.FindStringSubmatch(url)
	if len(parts) != 5 {
		return nil, ErrorsUCRepoInvalid
	}
	ap.Runmode = parts[1]
	ap.Site = parts[2]
	ap.Namespace = parts[3]
	ap.Repo = parts[4]

	return ap, nil
}

// NRString returns 'namespace/repo'
func (ap UpdateClientAppV1Repo) NRString() string {
	return fmt.Sprintf("%s/%s", ap.Namespace, ap.Repo)
}

// String returns the full appV1 url
func (ap UpdateClientAppV1Repo) String() string {
	return fmt.Sprintf("%s://%s/%s/%s", ap.Runmode, ap.Site, ap.Namespace, ap.Repo)
}

func (ap UpdateClientAppV1Repo) generateURL() string {
	return fmt.Sprintf("%s://%s/%s/%s/%s", ap.Runmode, ap.Site, appV1Restful, ap.Namespace, ap.Repo)
}

// List lists the applications of a remote repository
func (ap UpdateClientAppV1Repo) List() ([]string, error) {
	url := fmt.Sprintf("%s/list", ap.generateURL())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type httpRet struct {
		Message string
		Content []string
	}

	var ret httpRet
	err = json.Unmarshal(respBody, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Content, nil
}

// GetFile gets the application data by its name
func (ap UpdateClientAppV1Repo) GetFile(fullname string) ([]byte, error) {
	// fullname:  os/arch/app
	if len(strings.Split(fullname, "/")) != 3 {
		return nil, errors.New("Invalid fullname in appV1 GetFile")
	}

	url := fmt.Sprintf("%s/%s", ap.generateURL(), fullname)
	return ap.getFromURL(url)
}

// GetMetaSign gets the meta signature data of a repository
func (ap UpdateClientAppV1Repo) GetMetaSign() ([]byte, error) {
	url := fmt.Sprintf("%s/metasign", ap.generateURL())
	return ap.getFromURL(url)
}

// GetMeta gets the meta data of a repository
func (ap UpdateClientAppV1Repo) GetMeta() ([]byte, error) {
	url := fmt.Sprintf("%s/meta", ap.generateURL())
	return ap.getFromURL(url)
}

// GetPublicKey gets the public key data of a repository
func (ap UpdateClientAppV1Repo) GetPublicKey() ([]byte, error) {
	url := fmt.Sprintf("%s://%s/%s/%s/pubkey", ap.Runmode, ap.Site, appV1Restful, ap.Namespace)
	return ap.getFromURL(url)
}

func (ap UpdateClientAppV1Repo) getFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return respBody, nil
}

// Put adds an application with a name to a repository
func (ap UpdateClientAppV1Repo) Put(name string, content []byte, method utils.EncryptMethod) error {
	url := fmt.Sprintf("%s/%s", ap.generateURL(), name)
	r := bytes.NewReader(content)
	req, err := http.NewRequest("PUT", url, r)
	req.Header.Set("Dockyard-Encrypt-Method", string(method))

	if err != nil {
		return err
	}
	//TODO: set head
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

// Delete removes an application with a name from a repository
func (ap UpdateClientAppV1Repo) Delete(name string) error {
	url := fmt.Sprintf("%s/%s", ap.generateURL(), name)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
