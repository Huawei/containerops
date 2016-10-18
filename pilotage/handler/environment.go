package handler

import (
	"encoding/json"
	"net/http"

	"gopkg.in/macaron.v1"
)

//PostEnvironmentV1Handler is
func PostEnvironmentV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetEnvironmentListV1Hander is
func GetEnvironmentListV1Hander(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetEnvironmentV1Handler is
func GetEnvironmentV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutEnvironmentV1Handler is
func PutEnvironmentV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//DeleteEnvironmentV1Handler is
func DeleteEnvironmentV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}
