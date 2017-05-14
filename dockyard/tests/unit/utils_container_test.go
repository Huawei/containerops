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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/engine-api/types/container"
	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/utils"
)

func TestUtilsContainer(t *testing.T) {
	//TODO: dockyard dev team should provide small testing containers.
	imageName := "google/nodejs"
	containerName := ""
	cached, err := utils.IsImageCached(imageName)
	if err == utils.ErrorsNoDockerClient {
		fmt.Println("Please start a docker daemon to continue the container operation test")
		return
	}

	assert.Nil(t, err, "Fail to load Image")

	if !cached {
		err := utils.PullImage(imageName)
		assert.Nil(t, err, "Fail to pull image")
	}

	tmpFile, err := ioutil.TempFile("/tmp", "dockyard-test-container-oper")
	assert.Nil(t, err, "System err, fail to create temp file")
	defer os.Remove(tmpFile.Name())

	var config container.Config
	config.Image = imageName
	config.Cmd = []string{"touch", tmpFile.Name()}
	var hostConfig container.HostConfig
	hostConfig.Binds = append(hostConfig.Binds, "/tmp:/tmp")

	utils.StartContainer(config, hostConfig, containerName)
	//TODO: stop, remove the container process

	assert.Equal(t, true, utils.IsFileExist(tmpFile.Name()), "Fail to touch file by using StartContainer")

}
