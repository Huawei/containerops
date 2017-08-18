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

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {

	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	target, url, key, err := parseEnv(data)

	err = update(target, url, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to update service %s: %s\n", target, err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT=true\n")
}

func parseEnv(env string) (target, url, sshKey string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		err = fmt.Errorf("CO_DATA value is null\n")
		return
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "target":
			target = value
		case "url":
			url = value
		case "key":
			sshKey = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}
	return
}

func update(target, url, key string) error {
	cmd := fmt.Sprintf("/var/containerops/scripts/%s/deploy.sh '%s'", target, url)
	return sshCommand("root", key, "hub.opshub.sh", 22, cmd, os.Stdout, os.Stderr)
}

func sshCommand(user, privateKey, host string, port int, command string, stdout, stderr io.Writer) error {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKeyFile(privateKey),
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

	if session, err = client.NewSession(); err != nil {
		return err
	}
	defer session.Close()

	if err != nil {
		return err
	}

	session.Stdout = stdout
	session.Stderr = stderr

	err = session.Run(command)
	if err != nil {
		return err
	}

	return nil
}

func publicKeyFile(file string) ssh.AuthMethod {
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
