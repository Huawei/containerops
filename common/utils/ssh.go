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

package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"path"
)

// PublicKeyFile parse private key file, and return ssh.AuthMethod.
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

//SSHCommand execute command from SSH connect in the remote host.
func SSHCommand(user, privateKey, host string, port int, commands []string, stdout, stderr io.Writer) error {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
		cmd          string
	)

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(privateKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 0,
	}

	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	defer client.Close()

	if session, err = client.NewSession(); err != nil {
		return err
	}
	defer session.Close()

	if err != nil {
		return err
	}

	session.Stdout = stdout
	session.Stderr = stderr

	if user != "root" {
		c := []string{}

		for _, command := range commands {
			fmt.Println(command)
			c = append(c, fmt.Sprintf("sudo %s", command))
		}

		cmd = strings.Join(c, " && ")
	} else {
		cmd = strings.Join(commands, " && ")
	}

	err = session.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}

//SSHScp copy local file to remote dest using scp command.
func SSHScp(user, privateKey, host string, port int, files []map[string]string) error {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		c            *ssh.Client
		err          error
	)

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(privateKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 0,
	}

	addr = fmt.Sprintf("%s:%d", host, port)

	if c, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	defer c.Close()

	for _, file := range files {
		client, err := sftp.NewClient(c)
		if err != nil {
			return err
		}

		localFile, err := os.Open(file["src"])
		if err != nil {
			return err
		}

		if user == "root" {
			remoteFile, err := client.Create(file["dest"])
			if err != nil {
				return err
			}

			_, err = io.Copy(remoteFile, localFile)
		} else {
			tmp := path.Join("/tmp", path.Base(file["dest"]))
			remoteFile, err := client.Create(tmp)
			if err != nil {
				return err
			}

			_, err = io.Copy(remoteFile, localFile)

			mv := []string{
				fmt.Sprintf("mv %s %s", tmp, file["dest"]),
			}

			if err := SSHCommand(user, privateKey, host, port, mv, os.Stdout, os.Stderr); err != nil {
				return err
			}
		}
	}

	return nil
}
