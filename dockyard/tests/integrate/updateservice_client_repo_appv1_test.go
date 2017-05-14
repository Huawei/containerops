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
package integratetest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/updateservice/client"
	"github.com/Huawei/dockyard/utils"
)

func getTestURL() string {
	// Start a dockyard web server and set the enviornment like this:
	//     $ export US_TEST_SERVER=http://localhost:1234
	server := os.Getenv("US_TEST_SERVER")
	if server == "" {
		return ""
	}

	//TODO: need to clean the repo in the server after finish the testing
	namespace := "namespace-" + uuid.NewV4().String()
	repository := "repository-" + uuid.NewV4().String()
	return fmt.Sprintf("%s/%s/%s", server, namespace, repository)
}

// TestOper tests add/get/getmeta/getmetasign/list
func TestOper(t *testing.T) {
	var appV1Repo client.UpdateClientAppV1Repo

	validURL := getTestURL()

	// Skip the test if the testing enviornment is not ready
	if validURL == "" {
		fmt.Printf("Skip the '%s' test since the testing enviornment is not ready.\n", "List")
		return
	}

	testFiles := []string{"osA/archA/appA", "osB/archB/appB"}

	f, _ := appV1Repo.New(validURL)
	defer func() {
		for _, tf := range testFiles {
			err := f.Delete(tf)
			assert.Nil(t, err, "Fail to delete file")
		}
	}()

	// Init the data and also test the put function
	_, path, _, _ := runtime.Caller(0)
	for _, tf := range testFiles {
		file := filepath.Join(filepath.Dir(path), "testdata", tf)
		content, _ := ioutil.ReadFile(file)
		err := f.Put(tf, content, utils.EncryptNone)
		assert.Nil(t, err, "Fail to put file")
	}

	// Test list
	l, err := f.List()
	assert.Nil(t, err, "Fail to list")
	assert.Equal(t, len(l), 2, "Fail to list or something wrong in put")
	ok := (l[0] == testFiles[0] && l[1] == testFiles[1]) || (l[0] == testFiles[1] && l[1] == testFiles[0])
	assert.Equal(t, true, ok, "Fail to list the correct data")

	// Test get file
	fileBytes, err := f.GetFile(testFiles[0])
	assert.Nil(t, err, "Fail to get file")
	expectedBytes, _ := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "testdata", testFiles[0]))
	assert.Equal(t, fileBytes, expectedBytes, "Fail to get the correct data")

	// Test get meta
	metaBytes, err := f.GetMeta()
	assert.Nil(t, err, "Fail to get meta file")

	// Test get metasign
	signBytes, err := f.GetMetaSign()
	assert.Nil(t, err, "Fail to get meta signature file")

	// Test get public key
	pubkeyBytes, err := f.GetPublicKey()
	assert.Nil(t, err, "Fail to get public key file")

	// VIP: Verify meta/sign with public to make real sure that everything works perfect
	err = utils.SHA256Verify(pubkeyBytes, metaBytes, signBytes)
	assert.Nil(t, err, "Fail to verify the meta data")
}
