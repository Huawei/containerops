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
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/dockyard/models"
	"github.com/Huawei/dockyard/module/signature"
	"github.com/Huawei/dockyard/utils"
	"github.com/containerops/configure"
)

//AppcDiscoveryV1Handler is
func AppcDiscoveryV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	discovery := ctx.Query("ac-discovery")

	if len(discovery) > 0 && discovery == "1" {
		if t, err := template.ParseGlob("views/aci/discovery.html"); err != nil {
			log.Errorf("[%s] get gpg file template status: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(map[string]string{"Error": "Get GPG File Template Status Error"})
			return http.StatusBadRequest, result
		} else {
			t.Execute(ctx.Resp, map[string]string{
				"Domains":    configure.GetString("deployment.domains"),
				"Namespace":  namespace,
				"Repository": repository,
			})
		}
	}

	return http.StatusOK, []byte("")
}

//AppcGetACIV1Handler is
func AppcGetACIV1Handler(ctx *macaron.Context) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	filename := ctx.Params(":file")

	aci := strings.Trim(filename, ".asc")
	//TODO: Decode aci file name with template.
	version := strings.Split(aci, "-")[1]

	r := new(models.AppcV1)
	i := new(models.ACIv1)

	if err := r.Get(namespace, repository); err != nil {
		log.Errorf("[%s] get AppcV1 repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Get Appc Repository Error."})

		ctx.Resp.WriteHeader(http.StatusBadRequest)
		ctx.Resp.Write(result)
		return
	}

	if err := i.Get(r.ID, version, aci); err != nil {
		log.Errorf("[%s] get ACIV1 data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Get ACI Data Error."})

		ctx.Resp.WriteHeader(http.StatusBadRequest)
		ctx.Resp.Write(result)
		return
	}

	var path string
	if asc := strings.Contains(filename, ".asc"); asc == true {
		path = i.Sign
	} else {
		path = i.Path
	}

	if file, err := os.Open(path); err != nil {
		log.Errorf("[%s] get File(%v) error: %s", ctx.Req.RequestURI, file, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Get ACI or ASC file Error."})

		ctx.Resp.WriteHeader(http.StatusBadRequest)
		ctx.Resp.Write(result)
		return
	} else {
		stat, _ := file.Stat()
		size := strconv.FormatInt(stat.Size(), 10)

		ctx.Resp.Header().Set("Content-Description", "File Transfer")
		ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
		ctx.Resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		ctx.Resp.Header().Set("Content-Length", size)
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

//AppcPUTDetails is
type AppcPUTDetails struct {
	ACIPushVersion string `json:"aci_push_version"`
	Multipart      bool   `json:"multipart"`
	ManifestURL    string `json:"upload_manifest_url"`
	SignatureURL   string `json:"upload_signature_url"`
	ACIURL         string `json:"upload_aci_url"`
	CompletedURL   string `json:"completed_url"`
}

//AppcPostACIV1Handler is
func AppcPostACIV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	aci := ctx.Params(":aci")
	//TODO: Decode aci file name with template.
	version := strings.Split(aci, "-")[1]

	r := new(models.AppcV1)
	if err := r.Put(namespace, repository); err != nil {
		log.Errorf("[%s] put AppcV1 repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Post Appc Repository Error."})
		return http.StatusBadRequest, result
	}

	prefix := fmt.Sprintf("https://%s/appc/%s/%s/push", configure.GetString("deployment.domains"), namespace, repository)
	appc := AppcPUTDetails{
		ACIPushVersion: "0.0.1",
		Multipart:      false,
		ManifestURL:    fmt.Sprintf("%s/%s/manifest/%s", prefix, version, aci),
		SignatureURL:   fmt.Sprintf("%s/%s/asc/%s", prefix, version, aci),
		ACIURL:         fmt.Sprintf("%s/%s/aci/%s", prefix, version, aci),
		CompletedURL:   fmt.Sprintf("%s/%s/complete/%s", prefix, version, aci),
	}

	result, _ := json.Marshal(appc)
	return http.StatusOK, result
}

//AppcPutManifestV1Handler is
func AppcPutManifestV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	version := ctx.Params(":version")
	aci := ctx.Params(":aci")

	r := new(models.AppcV1)
	if err := r.Get(namespace, repository); err != nil {
		log.Errorf("[%s] get AppcV1 repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Get Appc Repository Error."})
		return http.StatusBadRequest, result
	}

	data, _ := ctx.Req.Body().Bytes()
	i := new(models.ACIv1)
	if err := i.PutManifest(r.ID, version, aci, string(data)); err != nil {
		log.Errorf("[%s] put ACIV1 manifest error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Put ACI Manifest Error."})
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//AppcPutASCV1Handler is
func AppcPutASCV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	version := ctx.Params(":version")
	aci := ctx.Params(":aci")

	r := new(models.AppcV1)
	if err := r.Get(namespace, repository); err != nil {
		log.Errorf("[%s] get AppcV1 repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Get Appc Repository Error."})
		return http.StatusBadRequest, result
	}

	basePath := configure.GetString("appc.storage")
	imagePath := fmt.Sprintf("%s/%s/%s/%s", basePath, namespace, repository, version)
	filePath := fmt.Sprintf("%s/%s.asc", imagePath, aci)

	if !utils.IsDirExist(imagePath) {
		os.MkdirAll(imagePath, os.ModePerm)
	}

	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	if file, err := os.Create(filePath); err != nil {
		log.Errorf("[%s] Create asc file error:%s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Create .asc File Error."})
		return http.StatusBadRequest, result
	} else {
		io.Copy(file, ctx.Req.Request.Body)
	}

	i := new(models.ACIv1)
	if err := i.PutSign(r.ID, version, aci, filePath); err != nil {
		log.Errorf("[%s] write sign{.asc} data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Write .asc data Error."})
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//AppcPutACIV1Handler is
func AppcPutACIV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	version := ctx.Params(":version")
	aci := ctx.Params(":aci")

	r := new(models.AppcV1)
	if err := r.Get(namespace, repository); err != nil {
		log.Errorf("[%s] get AppcV1 repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Get Appc Repository Error."})
		return http.StatusBadRequest, result
	}

	basePath := configure.GetString("appc.storage")
	imagePath := fmt.Sprintf("%s/%s/%s/%s", basePath, namespace, repository, version)
	filePath := fmt.Sprintf("%s/%s", imagePath, aci)

	if !utils.IsDirExist(imagePath) {
		os.MkdirAll(imagePath, os.ModePerm)
	}

	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	if file, err := os.Create(filePath); err != nil {
		log.Errorf("[%s] Create aci file error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Create .aci File Error."})
		return http.StatusBadRequest, result
	} else {
		io.Copy(file, ctx.Req.Request.Body)
	}

	size, _ := utils.GetFileSize(filePath)

	i := new(models.ACIv1)
	if err := i.PutACI(r.ID, size, version, aci, filePath); err != nil {
		log.Errorf("[%s] write aci data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Write .aci Data Error."})
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//CompleteMsg is
type CompleteMsg struct {
	Success      bool   `json:"success"`
	Reason       string `json:"reason,omitempty"`
	ServerReason string `json:"serverreason,omitempty"`
}

//AppcPostCompleteV1Handler is
func AppcPostCompleteV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	version := ctx.Params(":version")
	aci := ctx.Params(":aci")

	complete := new(CompleteMsg)
	i := new(models.ACIv1)
	r := new(models.AppcV1)

	if err := r.Get(namespace, repository); err != nil {
		log.Errorf("[%s] get AppcV1 repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(CompleteMsg{
			Success:      false,
			Reason:       "",
			ServerReason: "Get Appc Repository Error.",
		})
		return http.StatusBadRequest, result
	}

	if err := i.Get(r.ID, version, aci); err != nil {
		log.Errorf("[%s] get ACIV1 data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(CompleteMsg{
			Success:      false,
			Reason:       "",
			ServerReason: "Get ACI Data Error.",
		})
		return http.StatusBadRequest, result
	}

	data, _ := ctx.Req.Body().Bytes()
	if err := json.Unmarshal(data, &complete); err != nil {
		log.Errorf("[%s] decode complete json dataerror: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(CompleteMsg{
			Success:      false,
			Reason:       "",
			ServerReason: "Decode Complete JSON Data Error.",
		})
		return http.StatusBadRequest, result
	}

	if complete.Success == true {
		root, _ := os.Getwd()
		sign := fmt.Sprintf("%s/%s", root, "external/signs")

		if err := signature.VerifyACISignature(i.Path, i.Sign, sign); err != nil {
			log.Errorf("[%s] Verify ACI signature error: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(CompleteMsg{
				Success:      false,
				Reason:       "",
				ServerReason: "Verify ACI Signature Error.",
			})
			return http.StatusBadRequest, result
		}

		if err := i.Unlocked(r.ID, version, aci); err != nil {
			log.Errorf("[%s] Unlocked ACI file error: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(CompleteMsg{
				Success:      false,
				Reason:       "",
				ServerReason: "Unlocked ACI File Error.",
			})

			return http.StatusBadRequest, result
		}

		result, _ := json.Marshal(CompleteMsg{
			Success:      true,
			Reason:       "",
			ServerReason: "",
		})
		return http.StatusOK, result
	}

	result, _ := json.Marshal(CompleteMsg{
		Success:      false,
		Reason:       complete.Reason,
		ServerReason: "",
	})

	return http.StatusBadRequest, result
}
