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
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/containerops/configure"
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"

	"github.com/Huawei/dockyard/models"
	"github.com/Huawei/dockyard/utils"
)

//GetPingV1Handler returns http.StatusOK(200) when Dockyard provide the Docker Registry V1 support.
//TODO: Add a config option for provide Docker Registry V1.
func GetPingV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//GetUsersV1Handler is Docker client login handler functoin, should be integration with [Crew](https://gitub.com/containerops/crew) project.
//TODO: Integration with Crew project.
func GetUsersV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//PostUsersV1Handler In Docker Registry V1, the Docker client will POST /v1/users to create an user.
//If the Dockyard allow create user in the CLI, should be integration with [Crew](https://github.com/containerops/crew).
//If don't, Dockyard returns http.StatusUnauthorized(401) for forbidden.
//TODO: Add a config option for allow/forbidden create user in the CLI, and integrated with [Crew](https://github.com/containerops/crew).
func PostUsersV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusUnauthorized, result
}

//PutTagV1Handler
func PutTagV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.

	//In Docker Registry V1, the repository json data in the body of `PUT /v1/:namespace/:repository`
	if body, err := ctx.Req.Body().String(); err != nil {
		log.Errorf("[%s] get tag from http body error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Tag JSON Error"})
		return http.StatusBadRequest, result
	} else {
		rege, _ := regexp.Compile(`"([[:alnum:]]+)"`)
		imageID := rege.FindStringSubmatch(body)

		tag := ctx.Params(":tag")
		namespace := ctx.Params(":namespace")
		repository := ctx.Params(":repository")

		t := new(models.DockerTagV1)
		if err := t.Put(imageID[1], tag, namespace, repository); err != nil {
			log.Errorf("[%s] put repository tag error: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(map[string]string{"Error": "Put Repository Tag Error"})
			return http.StatusBadRequest, result
		}
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//PutRepositoryImagesV1Handler
func PutRepositoryImagesV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.

	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	r := new(models.DockerV1)
	if err := r.Unlocked(namespace, repository); err != nil {
		log.Errorf("[%s] unlock repository error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Unlock Repository Error"})
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusNoContent, result
}

//GetRepositoryImagesV1Handler will return images json data.
func GetRepositoryImagesV1Handler(ctx *macaron.Context) (int, []byte) {
	var username string
	var err error

	if username, _, err = utils.DecodeBasicAuth(ctx.Req.Header.Get("Authorization")); err != nil {
		log.Errorf("[%s] decode Authorization error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Decode Authorization Error"})
		return http.StatusUnauthorized, result
	}

	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	r := new(models.DockerV1)
	if v1, err := r.Get(namespace, repository); err != nil {
		log.Errorf("[%s] get repository images data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Repository Images Error"})
		return http.StatusBadRequest, result
	} else {

		//If the Docker client use "X-Docker-Token", will return a randon token value.
		if ctx.Req.Header.Get("X-Docker-Token") == "true" {
			token := fmt.Sprintf("Token signature=%v,repository=\"%v/%v\",access=%v",
				utils.MD5(username), namespace, repository, "read")

			ctx.Resp.Header().Set("X-Docker-Token", token)
			ctx.Resp.Header().Set("WWW-Authenticate", token)
		}

		ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(v1.JSON)))

		return http.StatusOK, []byte(v1.JSON)
	}

}

//GetTagV1Handler is
func GetTagV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	r := new(models.DockerV1)
	if tags, err := r.GetTags(namespace, repository); err != nil {
		log.Errorf("[%s] get repository tags data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Repository Tags Error"})
		return http.StatusBadRequest, result
	} else {
		result, _ := json.Marshal(tags)
		ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(result)))

		return http.StatusOK, result
	}
}

//PutRepositoryV1Handler will create or update the repository, it's first step of Docker push.
//TODO: @1 When someone create or update the repository, it will be locked to forbidden others action include pull action.
//TODO: @2 Add a config option for allow/forbidden Docker client pull action when a repository is locked.
//TODO: @3 Intergated with [Crew](https://github.com/containerops/crew).
//TODO: @4 Token will be store in Redis, and link the push action with username@repository.
func PutRepositoryV1Handler(ctx *macaron.Context) (int, []byte) {
	var username, body string
	//var passwd string
	var err error

	if username, _, err = utils.DecodeBasicAuth(ctx.Req.Header.Get("Authorization")); err != nil {
		log.Errorf("[%s] decode Authorization error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Decode Authorization Error"})
		return http.StatusUnauthorized, result
	}

	//When integrated with crew, like this:
	//@1: username, passwd, _ := utils.DecodeBasicAuth(ctx.Req.Header.Get("Authorization"))
	//@2: username, passwd authorizated in Crew.

	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	//When integrated the Crew, should be check the privilage.
	if username != namespace {

	}

	//In Docker Registry V1, the repository json data in the body of `PUT /v1/:namespace/:repository`
	if body, err = ctx.Req.Body().String(); err != nil {
		log.Errorf("[%s] get repository json from http body error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Repository JSON Error"})
		return http.StatusBadRequest, result
	}

	//Create or update the repository.
	r := new(models.DockerV1)
	if e := r.Put(namespace, repository, body, ctx.Req.Header.Get("User-Agent")); e != nil {
		log.Errorf("[%s] put repository error: %s", ctx.Req.RequestURI, e.Error())

		result, _ := json.Marshal(map[string]string{"Error": "PUT Repository Error"})
		return http.StatusBadRequest, result
	}

	//If the Docker client use "X-Docker-Token", will return a randon token value.
	if ctx.Req.Header.Get("X-Docker-Token") == "true" {
		token := fmt.Sprintf("Token signature=%v,repository=\"%v/%v\",access=%v",
			utils.MD5(username), namespace, repository, "write")

		ctx.Resp.Header().Set("X-Docker-Token", token)
		ctx.Resp.Header().Set("WWW-Authenticate", token)
	}

	//TODO: When deploy multi instances of dockyard, the endpoints will schedule comply all instances stauts and arithmetic.
	ctx.Resp.Header().Set("X-Docker-Endpoints", configure.GetString("deployment.domains"))

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//GetImageAncestryV1Handler
func GetImageAncestryV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	imageID := ctx.Params(":image")

	image := new(models.DockerImageV1)
	if i, err := image.Get(imageID); err != nil {
		log.Errorf("[%s] get image ancestry error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Image Ancestry Error"})
		return http.StatusBadRequest, result
	} else {
		ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(i.Ancestry)))

		return http.StatusOK, []byte(i.Ancestry)
	}
}

//GetImageJSONV1Handler is getting image json data function.
//When docker client push an image, dockyard return http status code '400' or '404' if haven't it. Then the docker client will push the json data and layer file.
//If dockyard has the image and return 200, the docker client will ignore it and push another iamge.
func GetImageJSONV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	imageID := ctx.Params(":image")

	image := new(models.DockerImageV1)
	if i, err := image.Get(imageID); err != nil && err == gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"image": i.ImageID,
		}).Info("Image Not Found.")

		result, _ := json.Marshal(map[string]string{})
		return http.StatusNotFound, result
	} else if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("[%s] get image error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Image Error"})
		return http.StatusBadRequest, result
	} else {
		ctx.Resp.Header().Set("X-Docker-Checksum-Payload", i.Checksum)
		ctx.Resp.Header().Set("X-Docker-Size", fmt.Sprint(i.Size))
		ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(i.JSON)))

		return http.StatusOK, []byte(i.JSON)
	}
}

//GetImageLayerV1Handler
func GetImageLayerV1Handler(ctx *macaron.Context) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	imageID := ctx.Params(":image")

	image := new(models.DockerImageV1)
	if i, err := image.Get(imageID); err != nil {
		log.Errorf("[%s] get image ancestry error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Image Layer Error"})
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		ctx.Resp.Write(result)
		return
	} else {
		if file, err := os.Open(i.Path); err != nil {
			log.Errorf("[%s] get image layer file status: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(map[string]string{"Error": "Get Image Layer File Status Error"})
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			ctx.Resp.Write(result)
			return
		} else {
			size := strconv.FormatInt(i.Size, 10)

			ctx.Resp.Header().Set("Content-Description", "File Transfer")
			ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
			ctx.Resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", i.ImageID))
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
}

//PutImageJSONV1Handler is
func PutImageJSONV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.

	if body, err := ctx.Req.Body().String(); err != nil {
		log.Errorf("[%s] get image json from http body error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Image JSON Error"})
		return http.StatusBadRequest, result
	} else if err == nil {
		imageID := ctx.Params(":image")
		image := new(models.DockerImageV1)

		if err := image.PutJSON(imageID, body); err != nil {
			log.Errorf("[%s] put image json error: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(map[string]string{"Error": "Put Image JSON Error"})
			return http.StatusBadRequest, result
		}
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//PutImageLayerV1Handler is save image layer file in the server.
func PutImageLayerV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	imageID := ctx.Params(":image")

	basePath := configure.GetString("dockerv1.storage")
	imagePath := fmt.Sprintf("%s/images/%s", basePath, imageID)
	layerfile := fmt.Sprintf("%s/images/%s/%s", basePath, imageID, imageID)

	if !utils.IsDirExist(imagePath) {
		os.MkdirAll(imagePath, os.ModePerm)
	}

	if _, err := os.Stat(layerfile); err == nil {
		os.Remove(layerfile)
	}

	if file, err := os.Create(layerfile); err != nil {
		log.Errorf("[%s] Create image file error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Write Image File Error."})
		return http.StatusBadRequest, result
	} else {
		io.Copy(file, ctx.Req.Request.Body)
	}

	size, _ := utils.GetFileSize(layerfile)

	image := new(models.DockerImageV1)
	if err := image.PutLayer(imageID, layerfile, size); err != nil {
		log.Errorf("[%s] Failed to save image layer data error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Put Image Layer Data Error"})
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

//PutImageChecksumV1Handler is put image checksum and payload value in the database.
func PutImageChecksumV1Handler(ctx *macaron.Context) (int, []byte) {
	//TODO: If standalone == true, Dockyard will check HEADER Authorization; if standalone == false, Dockyard will check HEADER TOEKN.
	imageID := ctx.Params(":image")

	checksum := ctx.Req.Header.Get("X-Docker-Checksum")
	payload := ctx.Req.Header.Get("X-Docker-Checksum-Payload")

	image := new(models.DockerImageV1)
	if err := image.PutChecksum(imageID, checksum, payload); err != nil {
		log.Errorf("[%s] Failed to set image checksum and payload error: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"message": "Put Image Checksum And Payload Data Error"})
		return http.StatusBadRequest, result
	}

	//TODO: Verify the file's checksum.

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}
