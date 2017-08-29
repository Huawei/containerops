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

package objects

import (
	"fmt"
	"io"
	"time"
)

type Logger interface {
	WriteLog(log string, writer io.Writer) error
}

func WriteLog(obj Logger, log string, writer io.Writer, timestamp bool) error {
	if timestamp == true {
		log = fmt.Sprintf("[%d] %s", time.Now().Unix(), log)
	}

	if err := obj.WriteLog(log, writer); err != nil {
		return err
	}

	return nil
}

// Node is
type Node struct {
	ID     int    `json:"id" yaml:"id"`
	IP     string `json:"ip" yaml:"ip"`
	User   string `json:"user" yaml:"user"`
	Distro string `json:"distro" yaml:"distro"`
}

// Service is
type Service struct {
	Provider string `json:"provider" yaml:"provider"`
	Token    string `json:"token" yaml:"token"`
	Region   string `json:"region" yaml:"region"`
	Size     string `json:"size" yaml:"size"`
	Image    string `json:"image" yaml:"image"`
	Nodes    int    `json:"nodes" yaml:"nodes"`
}

// Tools is
type Tools struct {
	SSH SSH `json:"ssh" yaml:"ssh"`
}

// SSH is
type SSH struct {
	Private     string `json:"private" yaml:"private"`
	Public      string `json:"public" yaml:"public"`
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
}

// Component is
type Component struct {
	Binary  string   `json:"binary" yaml:"binary"`
	URL     string   `json:"url" yaml:"url"`
	Package bool     `json:"package" yaml:"package"`
	Systemd string   `json:"systemd" yaml:"systemd"`
	CA      string   `json:"ca" yaml:"ca"`
	Before  string   `json:"before" yaml:"before"`
	After   string   `json:"after" yaml:"after"`
	Logs    []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}

//WriteLog implement Logger interface.
func (c *Component) WriteLog(log string, writer io.Writer) error {
	c.Logs = append(c.Logs, log)

	if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
		return err
	}

	return nil
}
