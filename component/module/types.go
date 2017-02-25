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

package module

import (
	"encoding/json"
)

// ComponentBaseData is component's base data to describe a component
type ComponentBaseData struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ComponentData is get component detail's unit struct
type ComponentData struct {
	ComponentBaseData
	Input        interface{} `json:"input"`
	Output       interface{} `json:"output"`
	Env          []Env       `json:"env"`
	ImageName    string      `json:"imageName"`
	ImageTag     string      `json:"imageTag"`
	Timeout      int         `json:"timeout"`
	Type         string      `json:"type"`
	UseAdvanced  bool        `json:"useAdvanced"`
	Pod          interface{} `json:"pod"`
	Service      interface{} `json:"service"`
	ImageSetting *BuildInfo  `json:"imageSetting"`
}

// Env is the environment define info
type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CommonPodcfg is config that if a component not use advanced kube setting
type CommonPodcfg struct {
	Spec struct {
		Containers []struct {
			Resources struct {
				Limits struct {
					Cpu    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"limits"`
				Requests struct {
					Cpu    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"requests"`
			} `json:"resources"`
		} `json:"containers"`
	} `json:"spec"`
}

// CommonServicecfg is config that if a component not use advanced kube setting
type CommonServicecfg struct {
	Spec struct {
		Ports []struct {
			NodePort   string `json:"nodePort"`
			Port       string `json:"port"`
			TargetPort string `json:"targetPort"`
		} `json:"ports"`
		Type string `json:"type"`
	} `json:"spec"`
}

// BuildInfo is the build info when build an img
type BuildInfo struct {
	ID         int64      `json:"id,omitempty"`
	From       *ImageInfo `json:"from"`
	Push       *ImageInfo `json:"push"`
	EventShell struct {
		ComponentStart  string `json:"componentStart"`
		ComponentResult string `json:"componentResult"`
		ComponentStop   string `json:"componentStop"`
	} `json:"events"`
}

// DebugInfo is the info that debug an compoent need
type DebugInfo struct {
	ImageInfo
	Input       *json.RawMessage `json:"input"`
	Output      *json.RawMessage `json:"output"`
	Env         []Env            `json:"env"`
	Timeout     int              `json:"timeout"`
	UseAdvanced bool             `json:"useAdvanced"`
	Pod         *json.RawMessage `json:"pod"`
	Service     *json.RawMessage `json:"service"`
}

// ImageInfo is the image info
type ImageInfo struct {
	ImageName string `json:"imageName"`
	ImageTag  string `json:"imageTag"`
	AuthType  string `json:"authType,omitempty"`
	Username  string `json:"username,omitempty"`
	Pwd       string `json:"pwd,omitempty"`
}
