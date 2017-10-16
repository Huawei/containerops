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
	"io"
	"os"
	"strings"

	"github.com/Huawei/containerops/common/utils"
)

const (
	DistroUbuntu = "ubuntu"
	DistroCentOS = "centos"
)

//InitializationEnvironment init the environment of node.
func InitializationEnvironment(key, ip, user, distro string, stdout io.Writer) ([]string, error) {
	var commands []string

	if cmd, err := initFolders(key, ip, user, stdout); err != nil {
		return commands, err
	} else {
		commands = append(commands, strings.Join(cmd, " && "))
	}

	if cmd, err := initEnvironment(key, ip, user, distro, stdout); err != nil {
		return commands, err
	} else {
		commands = append(commands, strings.Join(cmd, " && "))
	}

	return commands, nil
}

//initFolders mkdir etc and data folder in the node.
func initFolders(key, ip, user string, stdout io.Writer) ([]string, error) {
	initCmd := []string{
		"mkdir -p /etc/kubernetes/ssl",
		"mkdir -p /etc/etcd/ssl",
		"mkdir -p /var/lib/etcd",
	}

	if err := utils.SSHCommand(user, key, ip, DefaultSSHPort, initCmd, stdout, os.Stderr); err != nil {
		return initCmd, err
	}

	return initCmd, nil
}

//initEnvironment init the environment of node for deployment Cloud Native stack.
func initEnvironment(key, ip, user, distro string, stdout io.Writer) ([]string, error) {
	initCmd := map[string][]string{
		DistroUbuntu: []string{
			"systemctl stop ufw",
			"/lib/systemd/systemd-sysv-install disable ufw",
		},
		DistroCentOS: []string{},
	}

	if err := utils.SSHCommand(user, key, ip, DefaultSSHPort, initCmd[distro], stdout, os.Stderr); err != nil {
		return initCmd[distro], err
	}

	return initCmd[distro], nil
}
