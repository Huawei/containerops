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
package unittest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/updateservice/client"
)

// TestMCRAppV1New
func TestMCRAppV1New(t *testing.T) {
	var appV1 client.UpdateClientAppV1Repo

	invalidURL := "containerops.me/containerops/official"
	_, err := appV1.New(invalidURL)
	assert.Equal(t, err, client.ErrorsUCRepoInvalid, "Fail to parse invalid url")

	invalidURL2 := "http://containerops.me/containerops"
	_, err = appV1.New(invalidURL2)
	assert.Equal(t, err, client.ErrorsUCRepoInvalid, "Fail to parse invalid url")

	validURL := "http://containerops.me/containerops/official"
	f, err := appV1.New(validURL)
	assert.Nil(t, err, "Fail to setup a valid repo")
	assert.Equal(t, f.String(), validURL, "Fail to parse url")
	assert.Equal(t, f.NRString(), "containerops/official", "Fail to parse url")
}
