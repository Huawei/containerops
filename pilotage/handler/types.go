package handler

type CommonResp struct {
	OK        bool   `json:"ok"`
	ErrorCode uint32 `json:"error_code,omitempty"`
	Message   string `json:"message",omitempty`
}

type CreateComponentResp struct {
	id uint64 `json:"id"`
	CommonResp
}
