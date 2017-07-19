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
	URI     string   `json:"uri" yaml:"uri"`
	Title   string   `json:"title" yaml:"title"`
	Version int64    `json:"version" yaml:"version"`
	Tag     string   `json:"tag" yaml:"tag"`
	Service Service  `json:"service" yaml:"service"`
	Infra   Infra    `json:"Infra" yaml:"yaml"`
	Logs    []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}

// Service is
type Service struct {
	Provider string `json:"provider" yaml:"provider"`
	Token    string `json:"token" yaml:"token"`
}

// Infra is
type Infra struct {
	Etcd       Etcd       `json:"etcd" yaml:"etcd"`
	Kubernetes Kubernetes `json:"kubernetes" yaml:"kubernetes"`
}

// Etcd is
type Etcd struct {
}

// Kubernetes is
type Kubernetes struct {
}
