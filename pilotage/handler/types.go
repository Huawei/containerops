package handler

import (
	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/module"
)

type errCode uint64

const (
	componentErrCode errCode = 00010000
	workflowErrCode  errCode = 00100000
	stageErrCode     errCode = 01000000
	actionErrCode    errCode = 10000000
)

type CommonResp struct {
	OK        bool    `json:"ok"`
	ErrorCode errCode `json:"error_code,omitempty"`
	Message   string  `json:"message",omitempty`
}

type ComponentResp struct {
	*models.Component `json:"component,omitempty"`
	CommonResp        `json:"common"`
}

type DebugEvent struct {
	Type    module.EventType `json:"type"`
	Content string           `json:"content"`
}

type DebugComponentMessage struct {
	DebugID     int64      `json:"debug_id"`
	Kubernetes  string     `json:"kubernetes,omitempty"`
	Input       string     `json:"input,omitempty"`
	Environment string     `json:"environment,omitempty"`
	Output      string     `json:"output,omitempty"`
	Event       DebugEvent `json:"event,omitempty"`
	CommonResp  `json:"common"`
}
