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

package signature

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
)

//VerifyACISignature is
func VerifyACISignature(aci, asc, pubkeys string) error {
	files, err := ioutil.ReadDir(pubkeys)
	if err != nil {
		return fmt.Errorf("Read pubkeys directory failed: %v", err.Error())
	}

	if len(files) <= 0 {
		return fmt.Errorf("No pubkey file found in %v", pubkeys)
	}

	var keyring openpgp.EntityList
	for _, file := range files {
		pubkeyfile, err := os.Open(pubkeys + "/" + file.Name())
		if err != nil {
			return err
		}
		defer pubkeyfile.Close()

		keyList, err := openpgp.ReadArmoredKeyRing(pubkeyfile)
		if err != nil {
			return err
		}

		if len(keyList) < 1 {
			return fmt.Errorf("Missing opengpg entity")
		}

		keyring = append(keyring, keyList[0])
	}

	acifile, err := os.Open(aci)
	if err != nil {
		return fmt.Errorf("Open ACI file failed: %v", err.Error())
	}
	defer acifile.Close()

	ascfile, err := os.Open(asc)
	if err != nil {
		return fmt.Errorf("Open signature file failed: %v", err.Error())
	}
	defer ascfile.Close()

	if _, err := acifile.Seek(0, 0); err != nil {
		return fmt.Errorf("Seek ACI file failed: %v", err)
	}
	if _, err := ascfile.Seek(0, 0); err != nil {
		return fmt.Errorf("Seek signature file: %v", err)
	}

	//Verify detached signature which default is ASCII format
	_, err = openpgp.CheckArmoredDetachedSignature(keyring, acifile, ascfile)

	if err == io.EOF {
		if _, err := acifile.Seek(0, 0); err != nil {
			return fmt.Errorf("Seek ACI file failed: %v", err)
		}
		if _, err := ascfile.Seek(0, 0); err != nil {
			return fmt.Errorf("Seek signature file: %v", err)
		}

		//try to verify detached signature with binary format
		_, err = openpgp.CheckDetachedSignature(keyring, acifile, ascfile)
	}

	if err == io.EOF {
		return fmt.Errorf("Signature format is invalid")
	}

	return err
}
