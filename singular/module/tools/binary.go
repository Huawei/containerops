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
	"io"
	"net/url"
	"os"

	"github.com/Huawei/containerops/common/utils"
)

const (
	DefaultSSHUser = "root"
	DefaultSSHPort = 22
)

//DownloadComponents is download component binary file to the host.
//If the src is URL, execute curl command in the host.
//If the src is local file, execute scp command upload to the host.
func DownloadComponent(files []map[string]string, host, private, user string, stdout io.Writer) error {
	if u, err := url.Parse(files[0]["src"]); err != nil {
		return fmt.Errorf("invalid src format, neither url or local path")
	} else {
		if u.Scheme == "" {
			if utils.IsFileExist(files[0]["src"]) == true {
				if err = uploadBinary(files, host, private, user); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("the file not exist")
			}
		} else {
			if err = downloadBinary(files, host, private, user, stdout); err != nil {
				return err
			}
		}

	}

	return nil
}

//downloadBinary exec curl command download binary in the node.
func downloadBinary(files []map[string]string, host, private, user string, stdout io.Writer) error {
	for _, file := range files {
		cmd := fmt.Sprintf("curl %s -o %s", file["src"], file["dest"])
		if err := utils.SSHCommand(user, private, host, DefaultSSHPort, []string{cmd}, stdout, os.Stderr); err != nil {
			return err
		}
	}

	return nil
}

//uploadBinary exec scp command copy local file to the node.
func uploadBinary(files []map[string]string, host, private, user string) error {
	if err := utils.SSHScp(user, private, host, DefaultSSHPort, files); err != nil {
		return err
	}

	return nil
}
