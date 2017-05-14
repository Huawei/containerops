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

package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/containerops/configure"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/dockyard/models"
	"github.com/Huawei/dockyard/module"
	"github.com/Huawei/dockyard/module/signature"
	"github.com/Huawei/dockyard/utils"
)

//GetPingV2Handler is https://github.com/docker/distribution/blob/master/docs/spec/api.md#api-version-check
func GetPingV2Handler(ctx *macaron.Context) (int, []byte) {
	//if len(ctx.Req.Header.Get("Authorization")) == 0 {
	//	ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")

	//TODO Bearer Token
	//realm -> describe the authorization endpoint
	//service -> describe the name of the service that hold the resources.
	//scope -> which describe the resources that needed to be accessed,and the operation requested by the client (pulled/pushed).
	//account -> an optional attribute which describes the account used for authentication.
	//	ctx.Resp.Header().Set("WWW-Authenticate",
	//		fmt.Sprintf("Bearer realm=\"https://dockayrd.sh/authorizate\" service=\"Dockyard Service\" scope=\"repository:genedna/dockyard\" account=\"genedna\""))

	//	result, _ := json.Marshal(map[string]string{})
	//	return http.StatusUnauthorized, result
	//}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//GetCatalogV2Handler is
func GetCatalogV2Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//HeadBlobsV2Handler is
func HeadBlobsV2Handler(ctx *macaron.Context) (int, []byte) {
	digest := ctx.Params(":digest")
	tarsum := strings.Split(digest, ":")[1]

	i := new(models.DockerImageV2)
	if err := i.Get(tarsum); err != nil && err == gorm.ErrRecordNotFound {
		log.Info("Not found blob: %s", tarsum)

		result, _ := module.EncodingError(module.BLOB_UNKNOWN, digest)
		return http.StatusNotFound, result
	} else if err != nil && err != gorm.ErrRecordNotFound {
		log.Info("Failed to get blob %s: %s", tarsum, err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, err.Error())
		return http.StatusBadRequest, result
	}

	ctx.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Docker-Content-Digest", digest)
	ctx.Resp.Header().Set("Content-Length", fmt.Sprint(i.Size))

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//PostBlobsV2Handler is
//Initiate a resumable blob upload. If successful, an upload location will be provided to complete the upload.
//Optionally, if the digest parameter is present, the request body will be used to complete the upload in a single request.
func PostBlobsV2Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	r := new(models.DockerV2)

	if err := r.Put(namespace, repository); err != nil {
		log.Errorf("Put or search repository error: %s", err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, map[string]string{"namespace": namespace, "repository": repository})
		return http.StatusBadRequest, result
	}

	uuid := utils.MD5(uuid.NewV4().String())
	state := utils.MD5(fmt.Sprintf("%s/%s/%d", namespace, repository, time.Now().UnixNano()/int64(time.Millisecond)))
	random := fmt.Sprintf("https://%s/v2/%s/%s/blobs/uploads/%s?_state=%s",
		configure.GetString("deployment.domains"), namespace, repository, uuid, state)

	ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Resp.Header().Set("Docker-Upload-Uuid", uuid)
	ctx.Resp.Header().Set("Location", random)
	ctx.Resp.Header().Set("Range", "0-0")

	result, _ := json.Marshal(map[string]string{})
	return http.StatusAccepted, result
}

//PatchBlobsV2Handler is
//Upload a chunk of data for the specified upload.
//Docker 1.9.x above version saves layer in PATCH methord
//Docker 1.9.x below version saves layer in PUT methord
func PatchBlobsV2Handler(ctx *macaron.Context) (int, []byte) {
	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")

	desc := ctx.Params(":uuid")
	uuid := strings.Split(desc, "?")[0]

	if upload, err := module.CheckDockerVersion19(ctx.Req.Header.Get("User-Agent")); err != nil {
		log.Errorf("Decode docker version error: %s", err.Error())

		result, _ := module.EncodingError(module.BLOB_UPLOAD_UNKNOWN, map[string]string{"namespace": namespace, "repository": repository})
		return http.StatusBadRequest, result
	} else if upload == true {
		//It's run above docker 1.9.0
		basePath := configure.GetString("dockerv2.storage")
		uuidPath := fmt.Sprintf("%s/uuid/%s", basePath, uuid)
		uuidFile := fmt.Sprintf("%s/uuid/%s/%s", basePath, uuid, uuid)

		if !utils.IsDirExist(uuidPath) {
			os.MkdirAll(uuidPath, os.ModePerm)
		}

		if _, err := os.Stat(uuidFile); err == nil {
			os.Remove(uuidFile)
		}

		if file, err := os.Create(uuidFile); err != nil {
			log.Errorf("[%s] Create UUID file error: %s", ctx.Req.RequestURI, err.Error())

			result, _ := module.EncodingError(module.BLOB_UPLOAD_UNKNOWN, map[string]string{"namespace": namespace, "repository": repository})
			return http.StatusBadRequest, result
		} else {
			io.Copy(file, ctx.Req.Request.Body)

			size, _ := utils.GetFileSize(uuidFile)

			ctx.Resp.Header().Set("Range", fmt.Sprintf("0-%v", size-1))
		}
	}

	state := utils.MD5(fmt.Sprintf("%s/%v", fmt.Sprintf("%s/%s", namespace, repository), time.Now().UnixNano()/int64(time.Millisecond)))
	random := fmt.Sprintf("https://%s/v2/%s/%s/blobs/uploads/%s?_state=%s",
		configure.GetString("deployment.domains"), namespace, repository, uuid, state)

	ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Resp.Header().Set("Docker-Upload-Uuid", uuid)
	ctx.Resp.Header().Set("Location", random)

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//PutBlobsV2Handler is
//Complete the upload specified by uuid, optionally appending the body as the final chunk.
func PutBlobsV2Handler(ctx *macaron.Context) (int, []byte) {
	var size int64

	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")

	desc := ctx.Params(":uuid")
	uuid := strings.Split(desc, "?")[0]

	digest := ctx.Query("digest")
	tarsum := strings.Split(digest, ":")[1]

	basePath := configure.GetString("dockerv2.storage")
	imagePath := fmt.Sprintf("%s/image/%s", basePath, tarsum)
	imageFile := fmt.Sprintf("%s/image/%s/%s", basePath, tarsum, tarsum)

	//Save from uuid path save too image path
	if !utils.IsDirExist(imagePath) {
		os.MkdirAll(imagePath, os.ModePerm)
	}

	if _, err := os.Stat(imageFile); err == nil {
		os.Remove(imageFile)
	}

	if upload, err := module.CheckDockerVersion19(ctx.Req.Header.Get("User-Agent")); err != nil {
		log.Errorf("Decode docker version error: %s", err.Error())

		result, _ := module.EncodingError(module.BLOB_UPLOAD_INVALID, map[string]string{"namespace": namespace, "repository": repository})
		return http.StatusBadRequest, result
	} else if upload == true {
		//Docker 1.9.x above version saves layer in PATCH method, in PUT method move from uuid to image:sha256
		uuidPath := fmt.Sprintf("%s/uuid/%s", basePath, uuid)
		uuidFile := fmt.Sprintf("%s/uuid/%s/%s", basePath, uuid, uuid)

		if _, err := os.Stat(uuidFile); err == nil {
			if err := os.Rename(uuidFile, imageFile); err != nil {
				log.Errorf("Move the temp file to image folder %s error: %s", imageFile, err.Error())

				result, _ := module.EncodingError(module.BLOB_UPLOAD_INVALID, map[string]string{"namespace": namespace, "repository": repository})
				return http.StatusBadRequest, result
			}

			size, _ = utils.GetFileSize(imagePath)

			os.RemoveAll(uuidFile)
			os.RemoveAll(uuidPath)
		}
	} else if upload == false {
		//Docker 1.9.x below version saves layer in PUT methord, save data to file directly.
		if file, err := os.Create(imageFile); err != nil {
			log.Errorf("Save the file %s error: %s", imageFile, err.Error())

			result, _ := module.EncodingError(module.BLOB_UPLOAD_INVALID, map[string]string{"namespace": namespace, "repository": repository})
			return http.StatusBadRequest, result
		} else {
			io.Copy(file, ctx.Req.Request.Body)
			size, _ = utils.GetFileSize(imagePath)
		}
	}

	i := new(models.DockerImageV2)
	if err := i.Put(tarsum, imageFile, size); err != nil {
		log.Errorf("Save the iamge data %s error: %s", tarsum, err.Error())

		result, _ := module.EncodingError(module.BLOB_UPLOAD_INVALID, map[string]string{"namespace": namespace, "repository": repository})
		return http.StatusBadRequest, result
	}

	state := utils.MD5(fmt.Sprintf("%s/%v", fmt.Sprintf("%s/%s", namespace, repository), time.Now().UnixNano()/int64(time.Millisecond)))
	random := fmt.Sprintf("https://%s/v2/%s/%s/blobs/uploads/%s?_state=%s",
		configure.GetString("deployment.domains"), namespace, repository, uuid, state)

	ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Resp.Header().Set("Docker-Content-Digest", digest)
	ctx.Resp.Header().Set("Location", random)

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//GetBlobsV2Handler is
func GetBlobsV2Handler(ctx *macaron.Context) {
	digest := ctx.Params(":digest")
	tarsum := strings.Split(digest, ":")[1]

	i := new(models.DockerImageV2)
	if err := i.Get(tarsum); err != nil && err == gorm.ErrRecordNotFound {
		log.Info("Not found blob: %s", tarsum)

		result, _ := module.EncodingError(module.BLOB_UNKNOWN, digest)
		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		log.Info("Failed to get blob %s: %s", tarsum, err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, err.Error())
		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if file, err := os.Open(i.Path); err != nil {
		log.Info("Failed to get blob %s: %s", tarsum, err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, err.Error())
		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		stat, _ := file.Stat()
		size := strconv.FormatInt(stat.Size(), 10)

		ctx.Resp.Header().Set("Content-Description", "File Transfer")
		ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
		ctx.Resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", i.ImageID))
		ctx.Resp.Header().Set("Content-Length", size)
		ctx.Resp.Header().Set("Docker-Content-Digest", digest)
		ctx.Resp.Header().Set("Expires", "0")
		ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
		ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
		ctx.Resp.Header().Set("Pragma", "public")

		file.Seek(0, 0)
		defer file.Close()

		io.Copy(ctx.Resp, file)
		ctx.Resp.WriteHeader(http.StatusOK)
		return
	}
}

//PutManifestsV2Handler is
func PutManifestsV2Handler(ctx *macaron.Context) (int, []byte) {
	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")
	agent := ctx.Req.Header.Get("User-Agent")
	tag := ctx.Params(":tag")

	if data, err := ctx.Req.Body().String(); err != nil {
		log.Errorf("Get the manifest data error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{})
		return http.StatusBadRequest, result
	} else {
		_, imageID, version, _ := module.GetTarsumlist([]byte(data))
		digest, _ := signature.DigestManifest([]byte(data))

		r := new(models.DockerV2)
		if err := r.PutAgent(namespace, repository, agent, strconv.FormatInt(version, 10)); err != nil {
			log.Errorf("Put the manifest data error: %s", err.Error())

			result, _ := json.Marshal(map[string]string{})
			return http.StatusBadRequest, result
		}

		t := new(models.DockerTagV2)
		if err := t.Put(namespace, repository, tag, imageID, data, strconv.FormatInt(version, 10)); err != nil {
			log.Errorf("Put the manifest data error: %s", err.Error())

			result, _ := json.Marshal(map[string]string{})
			return http.StatusBadRequest, result
		}

		random := fmt.Sprintf("https://%s/v2/%s/%s/manifests/%s",
			configure.GetString("deployment.domains"), namespace, repository, digest)

		ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		ctx.Resp.Header().Set("Docker-Content-Digest", digest)
		ctx.Resp.Header().Set("Location", random)

		status := []int{http.StatusBadRequest, http.StatusAccepted, http.StatusCreated}

		result, _ := json.Marshal(map[string]string{})
		return status[version], result
	}
}

//GetTagsListV2Handler is
func GetTagsListV2Handler(ctx *macaron.Context) (int, []byte) {
	var err error
	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")

	data := map[string]interface{}{}
	data["name"] = fmt.Sprintf("%s/%s", namespace, repository)

	r := new(models.DockerV2)

	if data["tags"], err = r.GetTags(namespace, repository); err != nil && err == gorm.ErrRecordNotFound {
		log.Info("Not found repository in getting tags list: %s/%s", namespace, repository)

		result, _ := module.EncodingError(module.BLOB_UNKNOWN, fmt.Sprintf("%s/%s", namespace, repository))
		return http.StatusNotFound, result
	} else if err != nil && err != gorm.ErrRecordNotFound {
		log.Info("Failed found repository in getting tags list %s/%s: %s", namespace, repository, err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, err.Error())
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(data)
	return http.StatusOK, result
}

//GetManifestsV2Handler is
func GetManifestsV2Handler(ctx *macaron.Context) (int, []byte) {
	repository := ctx.Params(":repository")
	namespace := ctx.Params(":namespace")
	tag := ctx.Params(":tag")

	t := new(models.DockerTagV2)
	if _, err := t.Get(namespace, repository, tag); err != nil && err == gorm.ErrRecordNotFound {
		log.Info("Not found repository in tetting tag manifest: %s/%s:%s", namespace, repository, tag)

		result, _ := module.EncodingError(module.BLOB_UNKNOWN, fmt.Sprintf("%s/%s", namespace, repository))
		return http.StatusNotFound, result
	} else if err != nil && err != gorm.ErrRecordNotFound {
		log.Info("Failed to get tag manifest %s/%s:%s : ", namespace, repository, tag, err.Error())

		result, _ := module.EncodingError(module.UNKNOWN, err.Error())
		return http.StatusBadRequest, result
	}

	digest, _ := signature.DigestManifest([]byte(t.Manifest))

	ctx.Resp.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
	ctx.Resp.Header().Set("Docker-Content-Digest", digest)
	ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(t.Manifest)))

	return http.StatusOK, []byte(t.Manifest)
}

//DeleteBlobsV2Handler is
func DeleteBlobsV2Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//DeleteBlobsUUUIDV2Handler is
func DeleteBlobsUUUIDV2Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//DeleteManifestsV2Handler is
func DeleteManifestsV2Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}
