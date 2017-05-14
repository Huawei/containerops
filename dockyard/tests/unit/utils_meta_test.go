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
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/utils"
)

// TestMetaEncrypt
func TestMetaEncrypt(t *testing.T) {
	var item utils.MetaItem

	item.SetEncryption(utils.EncryptRSA)
	assert.Equal(t, true, item.GetEncryption() == utils.EncryptRSA, "Fail to set/get entrypt method")
}

// TestMetaItemGenerate
func TestMetaItemGenerate(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(path), "testdata")

	testContentFile := filepath.Join(dir, "hello.txt")
	testHashFile := filepath.Join(dir, "hello.hash")
	contentByte, _ := ioutil.ReadFile(testContentFile)
	hashByte, _ := ioutil.ReadFile(testHashFile)
	metaItem := utils.GenerateMetaItem("hello.txt", contentByte)
	assert.Equal(t, metaItem.GetHash(), strings.TrimSpace(string(hashByte)), "Fail to get correct hash value")
}

// TestMetaTime
func TestMetaTime(t *testing.T) {
	test1 := "test1"
	test1Byte := []byte("test1 byte")
	metaItem1 := utils.GenerateMetaItem(test1, test1Byte)
	metaItem2 := metaItem1
	assert.Equal(t, metaItem1, metaItem2, "Fail to compare metaItem, should be the same")

	metaItem2.SetCreated(metaItem2.GetCreated().Add(time.Hour * 1))
	cmp := metaItem1.Compare(metaItem2)
	assert.Equal(t, cmp < 0, true, "Fail to compare metaItem, should be smaller")

	assert.Equal(t, metaItem2.IsExpired(), false, "Fail to get expired information")
	metaItem2.SetExpired(time.Now().Add(time.Hour * (-1)))
	assert.Equal(t, metaItem2.IsExpired(), true, "Fail to get expired information")
}
