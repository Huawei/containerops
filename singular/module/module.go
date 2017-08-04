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

package module

// Deployment is
type Deployment struct {
	URI       string                 `json:"uri" yaml:"uri"`
	Title     string                 `json:"title" yaml:"title"`
	Version   int64                  `json:"version" yaml:"version"`
	Tag       string                 `json:"tag" yaml:"tag"`
	Nodes     int                    `json:"nodes" yaml:"nodes"`
	Service   Service                `json:"service" yaml:"service"`
	Tools     Tools                  `json:"tools" yaml:"tools"`
	Infras    []Infra                `json:"infras" yaml:"infras"`
	Logs      []string               `json:"logs,omitempty" yaml:"logs,omitempty"`
	Config    string                 `json:"-" yaml:"-"`
	Verbose   bool                   `json:"-" yaml:"-"`
	Timestamp bool                   `json:"-" yaml:"-"`
	Outputs   map[string]interface{} `json:"-" yaml:"-"`
}

// Service is
type Service struct {
	Provider string `json:"provider" yaml:"provider"`
	Token    string `json:"token" yaml:"token"`
	Region   string `json:"region" yaml:"region"`
	Size     string `json:"size" yaml:"size"`
	Image    string `json:"image" yaml:"image"`
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

// Infra is
type Infra struct {
	Name       string      `json:"name" yaml:"name"`
	Version    string      `json:"version" yaml:"version" `
	Nodes      Nodes       `json:"nodes" yaml:"nodes"`
	Components []Component `json:"components" yaml:"components"`
}

// Nodes is
type Nodes struct {
	Master int `json:"master" yaml:"master"`
	Node   int `json:"node" yaml:"node"`
}

// Component is
type Component struct {
	Binary  string `json:"binary" yaml:"binary"`
	URL     string `json:"url" yaml:"url"`
	Package bool   `json:"package" yaml:"package"`
	Systemd string `json:"systemd" yaml:"systemd"`
	CA      string `json:"ca" yaml:"ca"`
	Before  string `json:"before" yaml:"before"`
	After   string `json:"after" yaml:"after"`
}
