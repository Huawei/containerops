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

package module

import (
	"encoding/json"
)

var ErrorDescription = make(map[string]string)

const (
	// This const parameters inherits from Docker V2 protocol.
	UNKNOWN               = "UNKNOWN"
	DIGEST_INVALID        = "DIGEST_INVALID"
	NAME_INVALID          = "NAME_INVALID"
	TAG_INVALID           = "TAG_INVALID"
	NAME_UNKNOWN          = "NAME_UNKNOWN"
	MANIFEST_UNKNOWN      = "MANIFEST_UNKNOWN"
	MANIFEST_INVALID      = "MANIFEST_INVALID"
	MANIFEST_UNVERIFIED   = "MANIFEST_UNVERIFIED"
	MANIFEST_BLOB_UNKNOWN = "MANIFEST_BLOB_UNKNOWN"
	BLOB_UNKNOWN          = "BLOB_UNKNOWN"
	BLOB_UPLOAD_UNKNOWN   = "BLOB_UPLOAD_UNKNOWN"
	BLOB_UPLOAD_INVALID   = "BLOB_UPLOAD_INVALID"

	// This const parameters added by ContainerOps team.
	REPOSITORY_CREATE_FAILED       = "REPOSITORY_CREATE_FAILED"
	REPOSITORY_CREATE_REDUPLICATED = "REPOSITORY_CREATE_REDUPLICATED"
	REPOSITORY_NONE                = "REPOSITORY_NONE"
	AUTHENTICATION_FAILED          = "AUTHENTICATION_FAILED"
	PARAMETER_UNKNOWN              = "PARAMETER_UNKNOWN"
)

func init() {
	// This error messages inherits from Docker V2 protocol.
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

	// This error messages added by ContainerOps Team
	ErrorDescription[REPOSITORY_CREATE_FAILED] = "repository created faild"
	ErrorDescription[REPOSITORY_CREATE_REDUPLICATED] = "repository created reduplicated"
	ErrorDescription[REPOSITORY_NONE] = "no repository"
	ErrorDescription[AUTHENTICATION_FAILED] = "authentication failed"
	ErrorDescription[PARAMETER_UNKNOWN] = "parameter unknown"
}

// Errors is
type Errors struct {
	Errors []ErrorUnit `json:"errors"`
}

// ErrorUnit is
type ErrorUnit struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

// EncodingError is
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
