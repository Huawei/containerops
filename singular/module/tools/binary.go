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
	"fmt"
	"net/url"
	"os"

	"github.com/Huawei/containerops/common/utils"
	//"github.com/Huawei/containerops/singular/module"
)

const (
	DefaultSSHUser = "root"
	DefaultSSHPort = 22
)

// DownloadComponents is download component binary file to the host.
// If the src is URL, execute curl command in the host.
// If the src is local file, execute scp command upload to the host.
//
func DownloadComponent(src, dest, host, private string) error {
	if _, err := url.Parse(src); err != nil {
		if utils.IsFileExist(src) == true {
			if err := uploadBinary(src, dest, host, private); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Invalid src format, neither URL or local path.")
		}
	} else {
		if err := downloadBinary(src, dest, host, private); err != nil {
			return err
		}
	}

	return nil
}

func downloadBinary(src, dest, host, private string) error {
	cmd := fmt.Sprintf("curl %s -o %s", src, dest)
	if err := utils.SSHCommand(DefaultSSHUser, private, host, DefaultSSHPort, cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func uploadBinary(file, dest, host, private string) error {
	if err := utils.SSHScp(DefaultSSHUser, private, host, DefaultSSHPort, file, dest); err != nil {
		return err
	}

	return nil
}
