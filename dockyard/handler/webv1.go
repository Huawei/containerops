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
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/configure"
)

//GetIndexPageV1Handler is the index page of Dockyard web.
//When with params `/?ac-discovery=1` means access from `rkt trust --prefix={domain}` possibility.
//Other access maybe come from web browser, and will generate HTML page then return.
func GetIndexPageV1Handler(ctx *macaron.Context) {
	var tFile string
	var t *template.Template
	var err error

	discovery := ctx.Query("ac-discovery")

	//Generate GPG html response from template.
	//TODO: Use the setting or environment parameter with GPG html template.
	//TODO: Use the const parameter instead of `1`, wow the ac-discovery value is only 1.
	if len(discovery) > 0 && discovery == "1" {
		if t, err = template.ParseGlob("views/aci/gpg.html"); err != nil {
			log.Errorf("[%s] get gpg file template status: %s", ctx.Req.RequestURI, err.Error())

			result, _ := json.Marshal(map[string]string{"Error": "Get GPG File Template Status Error"})

			ctx.Resp.Write(result)
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx.Resp.WriteHeader(http.StatusOK)
		t.Execute(ctx.Resp, map[string]string{"Domains": configure.GetString("deployment.domains")})

		return
	}

	if configure.GetString("runmode") == "dev" {
		tFile = "views/index.html"
	} else {
		tFile = "views/coming.html"
	}

	if t, err = template.ParseGlob(tFile); err != nil {
		log.Errorf("[%s] get gpg file template status: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get Index Template Status Error"})

		ctx.Resp.Write(result)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Resp.WriteHeader(http.StatusOK)
	t.Execute(ctx.Resp, map[string]string{
		"Title":   configure.GetString("site.domain"),
		"Salt":    configure.GetString("site.salt"),
		"RunMode": configure.GetString("runmode"),
	})

	return
}

//GetGPGFileV1Handler is downloading `dockyard.sh`'s GPG file.
func GetGPGFileV1Handler(ctx *macaron.Context) (int, []byte) {
	var file []byte
	var err error

	//TODO: Use the setting or environment paramete with GPG file path.
	if _, err := os.Stat("external/signs/pubkeys.gpg"); err != nil {
		log.Errorf("[%s] get gpg file status: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get GPG File Status Error"})
		return http.StatusBadRequest, result
	}

	if file, err = ioutil.ReadFile("external/signs/pubkeys.gpg"); err != nil {
		log.Errorf("[%s] get gpg file data: %s", ctx.Req.RequestURI, err.Error())

		result, _ := json.Marshal(map[string]string{"Error": "Get GPG File Data Error"})
		return http.StatusBadRequest, result
	}

	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(file)))

	return http.StatusOK, file
}

//GetNamespacePageV1Handler is the namespace page for all types.
func GetNamespacePageV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//GetRepositoryPageV1Handler is the repository page.
func GetRepositoryPageV1Handler(ctx *macaron.Context) (int, []byte) {
	//type := ctx.Params(":type") // -> appc/app/dockerv1/dockerv2
	//namespace := ctx.Params(":namespace")
	//repository := ctx.Params(":repository")

	return http.StatusOK, []byte("")
}

//PostRepositoryRESTV1Handler is REST API handler function for create repository from web page.
func PostRepositoryRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//GetRepositoryRESTV1Handler is
func GetRepositoryRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")

}

//PutRepositoryRESTV1Handler is
func PutRepositoryRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//DeleteRepositoryRESTV1Handler is
func DeleteRepositoryRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//GetPackagePageV1Handler is the package page.
//Docker V1 & V2 -> tag
//Appc -> {name}-{version}-{os}-{arch}.{ext}
//App -> {name}
func GetPackagePageV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//PostPackageRESTV1Handler is
func PostPackageRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//GetPackageRESTV1Hanfdler is
func GetPackageRESTV1Hanfdler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//PutPackageRESTV1Handler is
func PutPackageRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//DeletePacakgeRESTV1Handler is
func DeletePacakgeRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//GetManifestPageV1Handler is the manifest of package page.
//Docker V1 -> none
//Docker V2 -> tag manifest
//Appc -> appc manifest
//App -> manifest
func GetManifestPageV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//PostManifestRESTV1Handler is
func PostManifestRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//GetManifestRESTV1Handler is
func GetManifestRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//PutManifestRESTV1Handler is
func PutManifestRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

//DeleteManifestRESTV1Handler is
func DeleteManifestRESTV1Handler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
