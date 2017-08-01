/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

package tools

import (
	"os"
	"testing"

	"github.com/Huawei/containerops/common/utils"
)

func TestSSHKeys(t *testing.T) {
	// Clean test data
	if utils.IsFileExist("/tmp/ssh/id_rsa") == true {
		os.Remove("/tmp/ssh/id_rsa")
	}
	if utils.IsFileExist("/tmp/ssh/id_rsa.pub") == true {
		os.Remove("/tmp/ssh/id_rsa.pub")
	}

	//
	if publicKeyFile, privateKeyFile, fingerprint, err := GenerateSSHKeyFiles("/tmp"); err != nil {
		t.Errorf("Generate ssh key files error: %s", err.Error())
	} else {
		if f, err := OpenSSHKeyFiles(publicKeyFile, privateKeyFile); err != nil {
			t.Errorf("Open ssh key files error: %s", err.Error())
		} else {
			if f != fingerprint {
				t.Errorf("Fingerprint is different")
			}
		}
	}

}

func TestSSHKeyWithoutPublicKeyFile(t *testing.T) {
	// Clean test data
	if utils.IsFileExist("/tmp/ssh/id_rsa") == true {
		os.Remove("/tmp/ssh/id_rsa")
	}
	if utils.IsFileExist("/tmp/ssh/id_rsa.pub") == true {
		os.Remove("/tmp/ssh/id_rsa.pub")
	}

	//
	if publicKeyFile, privateKeyFile, fingerprint, err := GenerateSSHKeyFiles("/tmp"); err != nil {
		t.Errorf("Generate ssh key files error: %s", err.Error())
	} else {
		if err := os.RemoveAll(publicKeyFile); err != nil {
			t.Error(err)
		}

		if f, err := OpenSSHKeyFiles("", privateKeyFile); err != nil {
			t.Errorf("Open ssh key files error: %s", err.Error())
		} else {
			if f != fingerprint {
				t.Errorf("Fingerprint is different")
			}
		}
	}

}
