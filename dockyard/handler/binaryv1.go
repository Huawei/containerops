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

package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/dockyard/model"
	"github.com/Huawei/containerops/dockyard/module"
)

// PostBinaryV1Handler is
func PostBinaryV1Handler(ctx *macaron.Context) (int, []byte) {
	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")
	binary := ctx.Params(":binary")
	tag := ctx.Params(":tag")

	b := new(model.BinaryV1)
	if err := b.Get(namespace, repository); err != nil {
		log.Errorf("get repository error: %s", err.Error())

		result, _ := module.EncodingError(module.REPOSITORY_NONE, map[string]string{"namespace": namespace, "repository": repository})
		return http.StatusBadRequest, result
	}

	f := new(model.BinaryFileV1)
	if err := f.Get(b.ID, binary, tag); err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("get file error: %s", err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, map[string]string{"namespace": namespace, "repository": repository, "file": binary, "tag": tag})
		return http.StatusBadRequest, result
	}

	if f.ID == 0 {
		// Storage pattern namespace/repository/tag/file
		basePath := common.Storage.BinaryV1
		tagPath := fmt.Sprintf("%s/%s/%s/%s", basePath, namespace, repository, tag)
		binaryPath := fmt.Sprintf("%s/%s/%s/%s/%s", basePath, namespace, repository, tag, binary)

		if !utils.IsDirExist(tagPath) {
			os.MkdirAll(tagPath, os.ModePerm)
		}

		if _, err := os.Stat(binaryPath); err == nil {
			os.Remove(binaryPath)
		}

		if file, err := os.Create(binaryPath); err != nil {
			log.Errorf("[%s] Create binary file error: %s", ctx.Req.RequestURI, err.Error())

			result, _ := module.EncodingError(module.BLOB_UPLOAD_UNKNOWN, map[string]string{"namespace": namespace, "repository": repository, "file": binary, "tag": tag})
			return http.StatusBadRequest, result
		} else {
			io.Copy(file, ctx.Req.Request.Body)

			size, _ := utils.GetFileSize(binaryPath)
			sha512, _ := utils.GetFileSHA512(binaryPath)

			if err := f.Put(b.ID, size, binary, tag, sha512, binaryPath); err != nil {
				result, _ := module.EncodingError(module.BLOB_UPLOAD_UNKNOWN, map[string]string{"namespace": namespace, "repository": repository, "file": binary, "tag": tag})
				return http.StatusBadRequest, result
			}
		}
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

// GetBinaryV1Handler is
func GetBinaryV1Handler(ctx *macaron.Context) {
	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")
	binary := ctx.Params(":binary")
	tag := ctx.Params(":tag")

	b := new(model.BinaryV1)
	if err := b.Get(namespace, repository); err != nil {
		log.Errorf("get repository error: %s", err.Error())

		result, _ := module.EncodingError(module.REPOSITORY_NONE, map[string]string{"namespace": namespace, "repository": repository})
		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	f := new(model.BinaryFileV1)
	if err := f.Get(b.ID, binary, tag); err != nil {
		log.Errorf("get file error: %s", err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, map[string]string{"namespace": namespace, "repository": repository, "file": binary, "tag": tag})
		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if file, err := os.Open(f.Path); err != nil {
		result, _ := module.EncodingError(module.UNKNOWN, err.Error())
		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		stat, _ := file.Stat()
		size := strconv.FormatInt(stat.Size(), 10)

		ctx.Resp.Header().Set("Content-Description", "File Transfer")
		ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
		ctx.Resp.Header().Set("Content-Length", size)
		ctx.Resp.Header().Set("sha512", f.SHA512)
		ctx.Resp.Header().Set("Expires", "0")
		ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
		ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
		ctx.Resp.Header().Set("Pragma", "public")

		file.Seek(0, 0)
		defer file.Close()

		io.Copy(ctx.Resp, file)

		return
	}
}

// DeleteBinaryV1Handler is
func DeleteBinaryV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

// PutBinaryLabelV1Handler is
func PutBinaryLabelV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}
