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
	"testing"

	"github.com/fernet/fernet-go"
	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/utils"
)

// TestEncryptMethod
func TestEncryptMethod(t *testing.T) {
	cases := []struct {
		data     string
		expected utils.EncryptMethod
	}{
		{"rsa", utils.EncryptRSA},
		{"", utils.EncryptNone},
		{"anyother", utils.EncryptNotSupported},
	}

	for _, c := range cases {
		assert.Equal(t, utils.NewEncryptMethod(c.data), c.expected, "Fail to get encrypt method")
	}
}

// TestRSAGenerateEnDe
func TestRSAGenerateEnDe(t *testing.T) {
	privBytes, pubBytes, err := utils.GenerateRSAKeyPair(1024)
	assert.Nil(t, err, "Fail to genereate RSA Key Pair")

	testData := []byte("This is the testdata for encrypt and decryp")
	encrypted, err := utils.RSAEncrypt(pubBytes, testData)
	assert.Nil(t, err, "Fail to encrypt data")
	decrypted, err := utils.RSADecrypt(privBytes, encrypted)
	assert.Nil(t, err, "Fail to decrypt data")
	assert.Equal(t, testData, decrypted, "Fail to get correct data after en/de")
}

// TestSHA256Sign
func TestSHA256Sign(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(path), "testdata")

	testPrivFile := filepath.Join(dir, "rsa_private_key.pem")
	testContentFile := filepath.Join(dir, "hello.txt")
	testSignFile := filepath.Join(dir, "hello.sig")

	privBytes, _ := ioutil.ReadFile(testPrivFile)
	signBytes, _ := ioutil.ReadFile(testSignFile)
	contentBytes, _ := ioutil.ReadFile(testContentFile)
	testBytes, err := utils.SHA256Sign(privBytes, contentBytes)
	assert.Nil(t, err, "Fail to sign")
	assert.Equal(t, testBytes, signBytes, "Fail to get valid sign data ")
}

// TestSHA256Verify
func TestSHA256Verify(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(path), "testdata")

	testPubFile := filepath.Join(dir, "rsa_public_key.pem")
	testContentFile := filepath.Join(dir, "hello.txt")
	testSignFile := filepath.Join(dir, "hello.sig")

	pubBytes, _ := ioutil.ReadFile(testPubFile)
	signBytes, _ := ioutil.ReadFile(testSignFile)
	contentBytes, _ := ioutil.ReadFile(testContentFile)
	err := utils.SHA256Verify(pubBytes, contentBytes, signBytes)
	assert.Nil(t, err, "Fail to verify valid signed data")
	err = utils.SHA256Verify(pubBytes, []byte("Invalid content data"), signBytes)
	assert.NotNil(t, err, "Fail to verify invalid signed data")
}

func TestTokenMarshalUnmarshal(t *testing.T) {
	var fkey fernet.Key
	fkey.Generate()
	key := string(fkey.Encode())
	invalidKey := "invalidKey"

	var retInt int
	var testInt int
	testInt = 1024

	intResult, err := utils.TokenMarshal(testInt, invalidKey)
	assert.NotNil(t, err, "Fail to marshal int with invalid key")
	err = utils.TokenUnmarshal(string(intResult), key, &retInt)
	assert.NotNil(t, err, "Fail to unmarshal int with invalid key")

	intResult, err = utils.TokenMarshal(testInt, key)
	assert.Nil(t, err, "Fail to marshal int")
	err = utils.TokenUnmarshal(string(intResult), key, &retInt)
	assert.Nil(t, err, "Fail to unmarshal int")
	assert.Equal(t, testInt, retInt, "Fail to get the original int data")

	var retStr string
	var testStr string
	testStr = "hello, world"
	strResult, err := utils.TokenMarshal(testStr, key)
	assert.Nil(t, err, "Fail to marshal string")
	err = utils.TokenUnmarshal(string(strResult), key, &retStr)
	assert.Nil(t, err, "Fail to unmarshal string")
	assert.Equal(t, testStr, retStr, "Fail to get the original string data")
}
