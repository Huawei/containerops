package handler

import "github.com/Huawei/containerops/pilotage/models"

type errCode int64

const (
	componentErrCode errCode = 1 << (4 + iota)
	workflowErrCode
	stageErrCode
	actionErrCode
)

type CommonResp struct {
	OK        bool   `json:"ok"`
	ErrorCode uint32 `json:"error_code,omitempty"`
	Message   string `json:"message",omitempty`
}

type ComponentResp struct {
	models.Component `json:"component,omitempty"`
	CommonResp       `json:"common"`
}

type DebugComponentReq struct {
	Kubernetes  string `json:"kubernetes"`
	Input       string `json:"input"`
	Environment string `json:"environment"`
}

type DebugComponentResp struct {
	LogID      uint64 `json:"log_id"`
	CommonResp `json:"common"`
}
