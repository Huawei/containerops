package handler

import (
	"encoding/json"
	"github.com/Huawei/containerops/pilotage/module"
)

type CommonResp struct {
	OK        bool    `json:"ok"`
	ErrorCode errCode `json:"error_code,omitempty"`
	Message   string  `json:"message",omitempty`
}

type ComponentResp struct {
	ComponentReq `json:"component,omitempty"`
	CommonResp   `json:"common"`
}

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ComponentSetup struct {
	ImageName   string           `json:"image_name"`
	ImageTag    string           `json:"image_tag"`
	Timeout     int              `json:"timeout"`
	Type        string           `json:"type"`
	DataFrom    string           `json:"data_from"`
	UseAdvanced bool             `json:"use_advanced"`
	Pod         *json.RawMessage `json:"pod"`
	Service     *json.RawMessage `json:"service"`
}

type ComponentReq struct {
	ID             int64                  `json:"id"`
	Name           string                 `json:"name"`
	Version        string                 `json:"version"`
	Input          map[string]interface{} `json:"input"`
	Output         map[string]interface{} `json:"output"`
	Env            []Env                  `json:"env"`
	ComponentSetup `json:"setup"`
}

type ListComponentsReq struct {
}

type DebugEvent struct {
	Type    module.EventType `json:"type"`
	Content string           `json:"content"`
}

type DebugComponentMessage struct {
	DebugID    int64                  `json:"debug_id"`
	Kubernetes string                 `json:"kubernetes,omitempty"`
	Input      map[string]interface{} `json:"input"`
	Env        []Env                  `json:"env,omitempty"`
	Output     map[string]interface{} `json:"output"`
	Event      DebugEvent             `json:"event"`
	CommonResp `json:"common"`
}
