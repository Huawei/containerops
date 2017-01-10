package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/macaron.v1"
)

func PreChecker(ctx *macaron.Context) {
	if ctx.Query("forbidden") == "true" {
		errMsg, _ := json.Marshal(map[string]string{"errMsg": "request is forbidden"})
		ctx.Resp.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(ctx.Resp, string(errMsg))
		return
	}
}
