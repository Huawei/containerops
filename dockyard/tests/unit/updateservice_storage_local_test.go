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
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	sl "github.com/Huawei/dockyard/updateservice/storage"
	"github.com/Huawei/dockyard/utils"
)

func loadSLTestData(t *testing.T) sl.UpdateServiceStorage {
	var local sl.UpdateServiceStorageLocal

	_, path, _, _ := runtime.Caller(0)
	topPath := filepath.Join(filepath.Dir(path), "testdata")
	// In this test, storage dir and key manager dir is the same
	l, err := local.New(topPath, topPath)
	assert.Nil(t, err, "Fail to setup a local test storage")

	return l
}

// TestBasic
func TestSLBasic(t *testing.T) {
	var local sl.UpdateServiceStorageLocal

	validURL := "/tmp/containerops_storage_cache"
	ok := local.Supported(validURL)
	assert.Equal(t, ok, true, "Fail to get supported status")
	ok = local.Supported("localInvalid://tmp/containerops_storage_cache")
	assert.Equal(t, ok, false, "Fail to get supported status")

	l, err := local.New(validURL, "")
	assert.Nil(t, err, "Fail to setup a local storage")
	assert.Equal(t, l.String(), validURL)
}

func TestSLList(t *testing.T) {
	l := loadSLTestData(t)
	key := "containerops/official"
	validCount := 0

	apps, _ := l.List("app/v1", key)
	for _, app := range apps {
		if app == "os/arch/appA" || app == "os/arch/appB" {
			validCount++
		}
	}

	assert.Equal(t, validCount, 2, "Fail to get right apps")
}

func TestSLPutDelete(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "us-test-")
	defer os.RemoveAll(tmpPath)
	assert.Nil(t, err, "Fail to create temp dir")

	proto := "app/v1"
	testData := "this is test DATA, you can put in anything here"

	var local sl.UpdateServiceStorageLocal
	l, err := local.New(tmpPath, tmpPath)
	assert.Nil(t, err, "Fail to setup local repo")

	invalidKey := "containerops/official"
	_, err = l.Put(proto, invalidKey, []byte(testData), utils.EncryptNone)
	assert.NotNil(t, err, "Fail to put with invalid key")

	validKey := "containerops/official/os/arch/appA"
	_, err = l.Put(proto, validKey, []byte(testData), utils.EncryptNone)
	assert.Nil(t, err, "Fail to put key")

	_, err = l.GetMeta(proto, "containerops/official")
	assert.Nil(t, err, "Fail to get meta data")

	getData, err := l.Get(proto, validKey)
	assert.Nil(t, err, "Fail to load file")
	assert.Equal(t, string(getData), testData, "Fail to get correct file")

	err = l.Delete(proto, validKey)
	assert.Nil(t, err, "Fail to remove a file")
	err = l.Delete(proto, validKey)
	assert.NotNil(t, err, "Should return error in removing a non exist file")
}

func TestSLGet(t *testing.T) {
	l := loadSLTestData(t)

	proto := "app/v1"
	namespace := "containerops"
	key := "containerops/official"
	invalidKey := "containerops/official/invalid"

	_, err := l.GetPublicKey(proto, namespace)
	assert.Nil(t, err, "Fail to load public key")
	_, err = l.GetMetaSign(proto, key)
	assert.Nil(t, err, "Fail to load  sign file")

	_, err = l.GetMeta(proto, invalidKey)
	assert.NotNil(t, err, "Fail to get meta from invalid key")
	_, err = l.GetMeta(proto, key)
	assert.Nil(t, err, "Fail to load meta data")

	_, err = l.Get(proto, "invalidinput")
	assert.NotNil(t, err, "Fail to get by invalid key")

	data, err := l.Get(proto, key+"/os/arch/appA")
	expectedData := "This is the content of appA."
	assert.Nil(t, err, "Fail to load file")
	assert.Equal(t, string(data), expectedData, "Fail to get correct file")
}
