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
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/updateservice/snapshot"
	"github.com/Huawei/dockyard/utils"
)

func TestByContainerNew(t *testing.T) {
	cases := []struct {
		id       string
		url      string
		name     string
		expected bool
	}{
		{"id", "url", "bycontainer/busybox", true},
		{"", "url", "bycontainer/busybox", false},
		{"id", "", "bycontainer/busybox", false},
		{"id", "url", "", false},
	}

	var bc snapshot.UpdateServiceSnapshotByContainer
	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{CallbackID: c.id, DataURL: c.url, Name: c.name}
		_, err := bc.New(info)
		assert.Equal(t, c.expected, err == nil, "Fail to create new bycontainer snapshot")
	}
}

func TestByContainerSupported(t *testing.T) {
	cases := []struct {
		proto    string
		expected bool
	}{
		{"appv1", true},
		{"dockerv1", true},
		{"invalid", false},
	}

	var bycontainer snapshot.UpdateServiceSnapshotByContainer
	for _, c := range cases {
		assert.Equal(t, c.expected, bycontainer.Supported(c.proto), "Fail to get supported status")
	}
}

func TestByContainerProcess(t *testing.T) {
	for n, _ := range cbMap {
		delete(cbMap, n)
	}

	cases := []struct {
		id       string
		url      string
		name     string
		expected bool
	}{
		//We need to add a real scan image to make the test work, but that will be integrate test
		{"1", "snapshot", "bycontainer/ospaf/notexist", false},
		{"2", "snapshot", "bycontainer/busybox", true},
	}
	var bycontainer snapshot.UpdateServiceSnapshotByContainer
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(path), "testdata")

	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{CallbackID: c.id, DataURL: filepath.Join(dir, c.url), Name: c.name}
		a, _ := bycontainer.New(info)
		err := a.Process()
		if err == utils.ErrorsNoDockerClient {
			fmt.Println("No Docker client detected")
			return
		}
		assert.Nil(t, err, "Fail to process")
	}
}
