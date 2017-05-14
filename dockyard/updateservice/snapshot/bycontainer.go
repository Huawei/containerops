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
package snapshot

import (
	"errors"
	"fmt"

	"github.com/docker/engine-api/types/container"

	"github.com/Huawei/dockyard/utils"
)

var (
	snapshotProcess   = "snapshot"
	snapshotMountDir  = "/snapshot-data"
	byContainerName   = "bycontainer"
	byContainerProtos = []string{"appv1", "dockerv1"}
)

type UpdateServiceSnapshotByContainer struct {
	info SnapshotInputInfo
}

func init() {
	RegisterSnapshot(byContainerName, &UpdateServiceSnapshotByContainer{})
}

func (m *UpdateServiceSnapshotByContainer) New(info SnapshotInputInfo) (UpdateServiceSnapshot, error) {
	_, imagename := info.GetName()
	if info.CallbackID == "" || info.DataURL == "" || imagename == "" {
		return nil, errors.New("'ID' , 'URL', 'Image Name' should not be empty")
	}

	if info.CallbackFunc != nil {
		fmt.Println("It is useless to set callback func in 'bycontainer' plugin!")
	}

	m.info = info
	return m, nil
}

func (m *UpdateServiceSnapshotByContainer) Supported(proto string) bool {
	for _, p := range byContainerProtos {
		if p == proto {
			return true
		}
	}

	return false
}

func (m *UpdateServiceSnapshotByContainer) Process() error {
	_, imageName := m.info.GetName()
	cached, err := utils.IsImageCached(imageName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !cached {
		//TODO: remove this image to save server disk
		err := utils.PullImage(imageName)
		if err != nil {
			return err
		}
	}

	var config container.Config
	config.Image = imageName
	config.Cmd = []string{snapshotProcess, m.info.CallbackID, m.info.Host, m.info.DataProto}
	var hostConfig container.HostConfig
	hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s", m.info.DataURL, snapshotMountDir))
	containerName := "scan-" + m.info.CallbackID

	err = utils.StartContainer(config, hostConfig, containerName)

	return err
}

func (m *UpdateServiceSnapshotByContainer) Description() string {
	desc := `Group Snapshot. Scan the data by running container and collect its output.
		"The scan image should have a main process called <%s> and a mounted scan dir <%s>.`
	return fmt.Sprintf(desc, snapshotProcess, snapshotMountDir)
}
