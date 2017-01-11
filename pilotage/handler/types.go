package handler

import (
	"encoding/json"
	"github.com/Huawei/containerops/pilotage/module"
)

type CommonResp struct {
	OK        bool    `json:"ok"`
	ErrorCode ErrCode `json:"error_code,omitempty"`
	Message   string  `json:"message",omitempty`
}

type ComponentResp struct {
	*ComponentReq `json:"component,omitempty"`
	CommonResp    `json:"common"`
}

type ComponentItem struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ListComponentsResp struct {
	Components []ComponentItem `json:"components,omitempty"`
	CommonResp `json:"common"`
}

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ImageSetting struct {
	From            ImageInfo `json:"from"`
	Build           ImageInfo `json:"build"`
	ComponentStart  string `json:"component_start"`
	ComponentResult string `json:"component_result"`
	ComponentStop   string `json:"component_stop"`
}

type ImageInfo struct {
	ImageName string `json:"image_name"`
	ImageTag  string `json:"image_tag"`
}

type ComponentSetup struct {
	ImageName    string `json:"image_name"`
	ImageTag     string `json:"image_tag"`
	ImageSetting `json:"image_setting"`
	Timeout      int              `json:"timeout"`
	Type         string           `json:"type"`
	UseAdvanced  bool             `json:"use_advanced"`
	Pod          *json.RawMessage `json:"pod"`
	Service      *json.RawMessage `json:"service"`
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

type DebugEvent struct {
	Type    module.EventType `json:"type"`
	Content string           `json:"content"`
}

type DebugComponentMessage struct {
	DebugID    int64                  `json:"debug_id"`
	Kubernetes string                 `json:"kubernetes,omitempty"`
	Input      map[string]interface{} `json:"input,omitempty"`
	Env        []Env                  `json:"env,omitempty"`
	Output     map[string]interface{} `json:"output,omitempty"`
	Event      *DebugEvent            `json:"event,omitempty"`
	CommonResp `json:"common"`
}

type EventReq struct {
	EventID   int64               `json:"event_id"`
	EventType module.EventType    `json:"event_type"`
	RunID     string              `json:"run_id"`
	Info      module.EventReqInfo `json:"info"`
}
