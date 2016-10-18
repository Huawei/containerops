package handler

import (
	"encoding/json"
	"net/http"

	"gopkg.in/macaron.v1"
)

//
func IndexV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "Pilotage Backend REST API Service"})
	return http.StatusOK, result
}
