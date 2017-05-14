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

package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fernet/fernet-go"
)

type EncryptMethod string

const (
	EncryptRSA          = "rsa"
	EncryptNone         = "none"
	EncryptNotSupported = "not-supported"
)

func NewEncryptMethod(method string) EncryptMethod {
	switch method {
	case string(EncryptRSA):
		return EncryptRSA
	case "":
		return EncryptNone
	case string(EncryptNone):
		return EncryptNone
	}

	return EncryptNotSupported
}

// IsDirExist checks if a path is an existed dir
func IsDirExist(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}

	return fi.IsDir()
}

// IsFileExist checks if a file url is an exist file
func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// Containe checks if a target item exists in the list/map
func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)

	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

// ValidatePassword verifies if a password is valid
func ValidatePassword(password string) error {
	if valida, _ := regexp.MatchString("[:alpha:]", password); valida != true {
		return fmt.Errorf("No alpha character in the password.")
	}

	if valida, _ := regexp.MatchString("[:digit:]", password); valida != true {
		return fmt.Errorf("No digital character in the password.")
	}

	if len(password) < 5 || len(password) > 30 {
		return fmt.Errorf("Password characters length should be between 5 - 30.")
	}

	return nil
}

// EncodeBasicAuth encode username and password into a base64 string
func EncodeBasicAuth(username string, password string) string {
	auth := username + ":" + password
	msg := []byte(auth)
	authorization := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(authorization, msg)
	return string(authorization)
}

// DecodeBasicAuth decode a base64 string into a username and a password
func DecodeBasicAuth(authorization string) (username string, password string, err error) {
	basic := strings.Split(strings.TrimSpace(authorization), " ")
	if len(basic) <= 1 {
		return "", "", err
	}

	decLen := base64.StdEncoding.DecodedLen(len(basic[1]))
	decoded := make([]byte, decLen)
	authByte := []byte(basic[1])
	n, err := base64.StdEncoding.Decode(decoded, authByte)

	if err != nil {
		return "", "", err
	}
	if n > decLen {
		return "", "", fmt.Errorf("Something went wrong decoding auth config")
	}

	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", "", fmt.Errorf("Invalid auth configuration file")
	}

	username = arr[0]
	password = strings.Trim(arr[1], "\x00")

	return username, password, nil
}

// MD5 generates a md value of a key automaticly
func MD5(key string) string {
	md5String := fmt.Sprintf("dockyard %s is a container %d hub", key, time.Now().Unix())
	h := md5.New()
	h.Write([]byte(md5String))

	return hex.EncodeToString(h.Sum(nil))
}

// GenerateRSAKeyPair generate a private key and a public key
func GenerateRSAKeyPair(bits int) ([]byte, []byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)

	if err != nil {
		return nil, nil, err
	}

	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}

	return pem.EncodeToMemory(privBlock), pem.EncodeToMemory(pubBlock), nil
}

// RSAEncrypt encrypts a content by a public key
func RSAEncrypt(keyBytes []byte, contentBytes []byte) ([]byte, error) {
	pubKey, err := getPubKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, contentBytes)
}

// RSADecrypt decrypts content by a private key
func RSADecrypt(keyBytes []byte, contentBytes []byte) ([]byte, error) {
	privKey, err := getPrivKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, privKey, contentBytes)
}

// SHA256Sign signs a content by a private key
func SHA256Sign(keyBytes []byte, contentBytes []byte) ([]byte, error) {
	privKey, err := getPrivKey(keyBytes)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(contentBytes)
	return rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
}

// SHA256Verify verifies if a content is valid by a signed data an a public key
func SHA256Verify(keyBytes []byte, contentBytes []byte, signBytes []byte) error {
	pubKey, err := getPubKey(keyBytes)
	if err != nil {
		return err
	}

	signStr := hex.EncodeToString(signBytes)
	newSignBytes, _ := hex.DecodeString(signStr)
	hashed := sha256.Sum256(contentBytes)
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], newSignBytes)
}

func getPrivKey(privBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privBytes)
	if block == nil {
		return nil, errors.New("Fail to decode private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func getPubKey(pubBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubBytes)
	if block == nil {
		return nil, errors.New("Fail to decode public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Fail get public key from public interface")
	}

	return pubKey, nil
}

// TokenMarshal encrypts data in `v`
func TokenMarshal(v interface{}, key string) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, err
	}

	k, err := fernet.DecodeKey(key)
	if err != nil {
		return nil, err
	}
	return fernet.EncryptAndSign(buf.Bytes(), k)
}

// TokenUnmarshal decryptes a token and save the original data to `v`.
func TokenUnmarshal(token string, key string, v interface{}) error {
	k, err := fernet.DecodeKey(key)
	if err != nil {
		return err
	}

	msg := fernet.VerifyAndDecrypt([]byte(token), time.Hour, []*fernet.Key{k})
	if msg == nil {
		return errors.New("invalid or expired token")
	}

	return json.NewDecoder(bytes.NewBuffer(msg)).Decode(&v)
}
