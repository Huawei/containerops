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
	"fmt"
	"github.com/fsouza/go-dockerclient"
	log "github.com/Sirupsen/logrus"
	"bytes"
)

type dockerImage struct {
	*docker.Client
	Name string
	Tag  string
}

func NewDockerImage(name, tag string) (*dockerImage, error) {
	client, err := dockerClient("localhost", 2376)
	if err != nil {
		return nil, fmt.Errorf("New docker client error: %s", err)
	}
	return &dockerImage{client, name, tag}, nil
}

func dockerClient(hostname string, port uint64) (*docker.Client, error) {
	endpoint := fmt.Sprintf("tcp://%s:%d", hostname, port)
	log.Debugf("Trying to Contect to dockerd[%s]\n", endpoint)
	return docker.NewClient(endpoint)
}

func (image *dockerImage) imageName() string {
	if image.Tag == "" {
		return image.Name
	} else {
		return fmt.Sprintf("%s:%s", image.Name, image.Tag)
	}
}

func (image *dockerImage) Build(directory string) (string, error) {
	var buf bytes.Buffer
	imageName := image.imageName()
	labels := make(map[string]string)
	labels["name"] = imageName
	labels["build"] = "workflow"
	opts := docker.BuildImageOptions{
		Name:           imageName,
		SuppressOutput: true,
		Pull:           true,
		OutputStream:   &buf,
		RmTmpContainer: true,
		ContextDir:     directory,
		Labels:         labels,
	}

	if err := image.Client.BuildImage(opts); err != nil {
		return "", fmt.Errorf("Build docker image error: %s", err)
	}
	log.Debugf("Build docker image[%s] finished\n", imageName)
	return buf.String(), nil
}

func (image *dockerImage) PushToRegistry(registry, username, password, email string) (string, error) {
	imageName := image.imageName()
	repo := fmt.Sprintf("%s/%s", registry, image.Name)
	err := image.Client.TagImage(imageName, docker.TagImageOptions{
		Repo: repo,
		Tag:  image.Tag,
	})
	log.Debugf("Tag docker image, repo[%s] tag[%s]\n", repo, image.Tag)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	opts := docker.PushImageOptions{
		Name:         repo,
		Tag:          image.Tag,
		Registry:     registry,
		OutputStream: &buf,
	}
	authConf := docker.AuthConfiguration{
		Username: username,
		Password: password,
		Email:    email,
	}
	log.Debugf("push docker image[%s][%s] to registry[%s]\n", registry, repo, image.Tag)
	if err := image.Client.PushImage(opts, authConf); err != nil {
		return "", fmt.Errorf("Push docker image to registry error: %s", err)
	}
	return buf.String(), nil
}
