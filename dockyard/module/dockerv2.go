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

package module

import (
	"encoding/json"
	"strconv"
	"strings"
)

var ErrorDescription = make(map[string]string)

const (
	//UNKNOWN is
	UNKNOWN = "UNKNOWN"
	//DIGEST_INVALID is
	DIGEST_INVALID = "DIGEST_INVALID"
	//NAME_INVALID is
	NAME_INVALID = "NAME_INVALID"
	//TAG_INVALID is
	TAG_INVALID = "TAG_INVALID"
	//NAME_UNKNOWN is
	NAME_UNKNOWN = "NAME_UNKNOWN"
	//MANIFEST_UNKNOWN is
	MANIFEST_UNKNOWN = "MANIFEST_UNKNOWN"
	//MANIFEST_INVALID is
	MANIFEST_INVALID = "MANIFEST_INVALID"
	//MANIFEST_UNVERIFIED is
	MANIFEST_UNVERIFIED = "MANIFEST_UNVERIFIED"
	//MANIFEST_BLOB_UNKNOWN is
	MANIFEST_BLOB_UNKNOWN = "MANIFEST_BLOB_UNKNOWN"
	//BLOB_UNKNOWN is
	BLOB_UNKNOWN = "BLOB_UNKNOWN"
	//BLOB_UPLOAD_UNKNOWN is
	BLOB_UPLOAD_UNKNOWN = "BLOB_UPLOAD_UNKNOWN"
	//BLOB_UPLOAD_INVALID is
	BLOB_UPLOAD_INVALID = "BLOB_UPLOAD_INVALID"
)

func init() {
	ErrorDescription[UNKNOWN] = "unknown error"
	ErrorDescription[DIGEST_INVALID] = "provided digest did not match uploaded content"
	ErrorDescription[NAME_INVALID] = "invalid repository name"
	ErrorDescription[TAG_INVALID] = "manifest tag did not match URI"
	ErrorDescription[NAME_UNKNOWN] = "repository name not known to registry"
	ErrorDescription[MANIFEST_UNKNOWN] = "manifest unknown"
	ErrorDescription[MANIFEST_INVALID] = "manifest invalid"
	ErrorDescription[MANIFEST_UNVERIFIED] = "manifest failed signature verification"
	ErrorDescription[MANIFEST_BLOB_UNKNOWN] = "blob unknown to registry"
	ErrorDescription[BLOB_UNKNOWN] = "blob unknown to registry"
	ErrorDescription[BLOB_UPLOAD_UNKNOWN] = "blob upload unknown to registry"
	ErrorDescription[BLOB_UPLOAD_INVALID] = "blob upload invalid"
}

//Errors is
type Errors struct {
	Errors []ErrorUnit `json:"errors"`
}

//ErrorUnit is
type ErrorUnit struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

//EncodingError is
func EncodingError(code string, detail interface{}) ([]byte, error) {
	var errs = Errors{}

	item := ErrorUnit{
		Code:    code,
		Message: ErrorDescription[code],
		Detail:  detail,
	}

	errs.Errors = append(errs.Errors, item)

	return json.Marshal(errs)
}

//CheckDockerVersion19 is
func CheckDockerVersion19(headers string) (bool, error) {
	agents := map[string]string{}
	for _, v := range strings.Split(headers, " ") {
		if len(strings.Split(v, "/")) > 1 {
			agents[strings.Split(v, "/")[0]] = strings.Split(v, "/")[1]
		}
	}

	versions := strings.Split(agents["docker"], ".")
	major, _ := strconv.ParseInt(versions[0], 10, 64)
	version, _ := strconv.ParseInt(versions[1], 10, 64)

	if major > 1 {
		return true, nil
	} else if major == 1 {
		if version > 9 {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, nil
}

//GetTarsumlist is
func GetTarsumlist(data []byte) ([]string, string, int64, error) {
	var tarsumlist []string
	var imageID string
	var layers = []string{"", "fsLayers", "layers"}
	var tarsums = []string{"", "blobSum", "digest"}

	var manifest map[string]interface{}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return []string{}, "", 0, err
	}

	schemaVersion := int64(manifest["schemaVersion"].(float64))

	if schemaVersion == 2 {
		confblobsum := manifest["config"].(map[string]interface{})["digest"].(string)
		imageID = strings.Split(manifest["config"].(map[string]interface{})["digest"].(string), ":")[1]
		tarsum := strings.Split(confblobsum, ":")[1]
		tarsumlist = append(tarsumlist, tarsum)
	}

	section := layers[schemaVersion]
	item := tarsums[schemaVersion]
	for i := len(manifest[section].([]interface{})) - 1; i >= 0; i-- {
		blobsum := manifest[section].([]interface{})[i].(map[string]interface{})[item].(string)
		tarsum := strings.Split(blobsum, ":")[1]
		tarsumlist = append(tarsumlist, tarsum)
	}

	return tarsumlist, imageID, schemaVersion, nil
}
