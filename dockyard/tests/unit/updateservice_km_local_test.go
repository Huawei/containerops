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

	kml "github.com/Huawei/dockyard/updateservice/km"
)

func loadTestKMLData(t *testing.T) (kml.KeyManager, string) {
	var local kml.KeyManagerLocal
	_, path, _, _ := runtime.Caller(0)
	realPath := filepath.Join(filepath.Dir(path), "testdata")

	l, err := local.New(realPath)
	assert.Nil(t, err, "Fail to setup a local test key manager")

	return l, realPath
}

func TestKMLBasic(t *testing.T) {
	var local kml.KeyManagerLocal

	validURL := "/tmp/containerops_km_cache"
	ok := local.Supported(validURL)
	assert.Equal(t, ok, true, "Fail to get supported status")
	ok = local.Supported("localInvalid://tmp/containerops_km_cache")
	assert.Equal(t, ok, false, "Fail to get supported status")

	_, err := local.New(validURL)
	assert.Nil(t, err, "Fail to setup a local key manager")
}

func TestKMLGetPublicKey(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "us-test-")
	defer os.RemoveAll(tmpPath)
	assert.Nil(t, err, "Fail to create temp dir")

	var local kml.KeyManagerLocal
	l, err := local.New(tmpPath)
	assert.Nil(t, err, "Fail to setup a local test key manager")

	namespace := "containerops"
	_, err = l.GetPublicKey("app/v1", namespace)
	assert.Nil(t, err, "Fail to get public key")
}

func TestKMLSign(t *testing.T) {
	proto := "app/v1"
	namespace := "containerops"
	l, realPath := loadTestKMLData(t)
	testFile := filepath.Join(realPath, "hello.txt")
	testBytes, _ := ioutil.ReadFile(testFile)
	signFile := filepath.Join(realPath, "hello.sig")
	signBytes, _ := ioutil.ReadFile(signFile)

	data, err := l.Sign(proto, namespace, testBytes)
	assert.Nil(t, err, "Fail to sign")
	assert.Equal(t, data, signBytes, "Fail to sign correctly")
}

func TestKMLDecrypt(t *testing.T) {
	proto := "app/v1"
	namespace := "containerops"
	l, realPath := loadTestKMLData(t)
	testFile := filepath.Join(realPath, "hello.txt")
	testBytes, _ := ioutil.ReadFile(testFile)
	testEncryptedFile := filepath.Join(realPath, "hello.encrypt")
	testEncryptedBytes, _ := ioutil.ReadFile(testEncryptedFile)

	testDecryptedBytes, err := l.Decrypt(proto, namespace, testEncryptedBytes)
	assert.Nil(t, err, "Fail to decrypt")
	assert.Equal(t, testDecryptedBytes, testBytes, "Fail to decrypt correctly")
}
