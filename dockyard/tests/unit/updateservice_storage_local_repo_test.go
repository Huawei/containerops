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
package unittest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/updateservice/km"
	"github.com/Huawei/dockyard/updateservice/storage"
	"github.com/Huawei/dockyard/utils"
)

// TestSLRepoBasic
func TestSLRepoBasic(t *testing.T) {
	topDir, err := ioutil.TempDir("", "dus-repo-test-")
	defer os.RemoveAll(topDir)
	assert.Nil(t, err, "Fail to create temp dir")

	proto := "app/v1"
	invalidNRs := []string{"", "a/b/c"}
	for _, invalidNR := range invalidNRs {
		_, err := storage.NewLocalRepo(topDir, proto, invalidNR)
		assert.NotNil(t, err, "Fail to return error while setup an invalid url")
	}

	// new
	validNR := "containerops/official"
	validNamespace := "containerops"
	r, err := storage.NewLocalRepo(topDir, proto, validNR)
	assert.Nil(t, err, "Fail to setup a valid url")
	assert.Equal(t, r.GetTopDir(), filepath.Join(topDir, proto, validNR), "Fail to get the correct top dir")
	assert.Equal(t, r.GetMetaFile(), filepath.Join(topDir, proto, validNR, "meta.json"), "Fail to get the default meta file")

	err = r.SetKM(topDir)
	assert.Nil(t, err, "Fail to set key manager")

	// add
	testData := map[string]string{
		"appA": "This is the content of appA.",
		"appB": "This is the content of appB.",
	}

	for name, value := range testData {
		_, err := r.Put(name, []byte(value), utils.EncryptNone)
		assert.Nil(t, err, "Fail to add a file")
	}

	// list
	names, err := r.List()
	assert.Nil(t, err, "Fail to list repo files")
	assert.Equal(t, len(names), len(testData), "Fail to add or list the same number")
	for _, name := range names {
		_, ok := testData[name]
		assert.Equal(t, ok, true, "Fail to list the correct data")
	}

	// remove
	removeFile := "appA"
	err = r.Delete(removeFile)
	assert.Nil(t, err, "Fail to remove valid file")
	err = r.Delete(removeFile)
	assert.NotNil(t, err, "Fail to remove invalid file")

	// get removed file
	_, err = r.Get(removeFile)
	assert.NotNil(t, err, "Fail to get removed file")

	// update (add with a exist name)
	updateFile := "appB"
	updateContent := "This is the content of updated appB."
	_, err = r.Put(updateFile, []byte(updateContent), utils.EncryptNone)
	assert.Nil(t, err, "Fail to add an exist file")

	res, err := r.Get(updateFile)
	assert.Nil(t, err, "Fail to get file")
	assert.Equal(t, string(res), updateContent)

	// update with encrypt gpg method
	_, err = r.Put(updateFile, []byte(updateContent), utils.EncryptRSA)
	assert.Nil(t, err, "Fail to add an exist file")

	encryptdRes, err := r.Get(updateFile)
	assert.Nil(t, err, "Fail to get file")
	assert.NotEqual(t, string(encryptdRes), updateContent)

	// decrypt by keymanager
	k, err := km.NewKeyManager(topDir)
	assert.Nil(t, err, "Fail to get key manager")
	decryptdRes, err := k.Decrypt(proto, validNamespace, encryptdRes)
	assert.Nil(t, err, "Fail to decrypt")
	assert.Equal(t, string(decryptdRes), updateContent)

	// get meta
	_, err = r.GetMeta()
	assert.Nil(t, err, "Fail to get meta data")

	// get metasign
	metaSignFile := r.GetMetaSignFile()
	ok := utils.IsFileExist(metaSignFile)
	assert.Equal(t, ok, true, "Fail to generate meta sign file")

	// get public key
	pubKeyFile := r.GetPublicKeyFile()
	ok = utils.IsFileExist(pubKeyFile)
	assert.Equal(t, ok, true, "Fail to generate public key file")
}
