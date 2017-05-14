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

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
)

var (
	dockerClient *client.Client = nil

	ErrorsNoDockerClient = errors.New("No docker client detected")
)

func init() {
	var err error
	dockerClient, err = client.NewClient("unix:///var/run/docker.sock", "", nil, nil)
	if err != nil {
		return
	}

	v, err := dockerClient.ServerVersion(context.Background())
	if err != nil {
		dockerClient = nil
		return
	}

	dockerClient.UpdateClientVersion(v.APIVersion)
}

func IsImageCached(imageName string) (bool, error) {
	if dockerClient == nil {
		return false, ErrorsNoDockerClient
	}

	options := types.ImageListOptions{MatchName: imageName}
	images, err := dockerClient.ImageList(context.Background(), options)
	return len(images) != 0, err
}

func PullImage(imageName string) error {
	if dockerClient == nil {
		return ErrorsNoDockerClient
	}

	body, err := dockerClient.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer body.Close()

	// wait until the image download is finished
	dec := json.NewDecoder(body)
	m := map[string]interface{}{}
	for {
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	// if the final stream object contained an error, return it
	if errMsg, ok := m["error"]; ok {
		return fmt.Errorf("%v", errMsg)
	}

	return nil
}

// StartContainer start and create a container
func StartContainer(config container.Config, hostConfig container.HostConfig, containerName string) error {
	if dockerClient == nil {
		return ErrorsNoDockerClient
	}

	resp, err := dockerClient.ContainerCreate(context.Background(), &config, &hostConfig, nil, containerName)
	if err != nil {
		return err
	}

	return dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
}
