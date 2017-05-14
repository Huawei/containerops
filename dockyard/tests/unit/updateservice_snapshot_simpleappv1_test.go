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
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Huawei/dockyard/updateservice/snapshot"
)

func TestAppv1New(t *testing.T) {
	cases := []struct {
		id       string
		url      string
		expected bool
	}{
		{"id", "url", true},
		{"", "url", false},
		{"id", "", false},
	}

	var appv1 snapshot.UpdateServiceSnapshotAppv1
	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{CallbackID: c.id, DataURL: c.url}
		_, err := appv1.New(info)
		assert.Equal(t, c.expected, err == nil, "Fail to create new snapshot appv1")
	}
}

func TestAppv1Supported(t *testing.T) {
	cases := []struct {
		proto    string
		expected bool
	}{
		{"appv1", true},
		{"invalid", false},
	}

	var appv1 snapshot.UpdateServiceSnapshotAppv1
	for _, c := range cases {
		assert.Equal(t, c.expected, appv1.Supported(c.proto), "Fail to get supported status")
	}
}

var (
	cbMap = make(map[string]snapshot.SnapshotOutputInfo)
)

func testCB(id string, data snapshot.SnapshotOutputInfo) error {
	if id != "1" && id != "2" {
		return errors.New("invalid id")
	}

	cbMap[id] = data
	return nil
}

func TestAppv1Process(t *testing.T) {
	for n, _ := range cbMap {
		delete(cbMap, n)
	}

	cases := []struct {
		id         string
		url        string
		cb         snapshot.Callback
		pExpected  bool
		idExpected bool
		md5        string
	}{
		{"1", "snapshot/testmd5", testCB, true, true, "ffe7c736f2aa54531ac6430e3cbf2545"},
		{"2", "snapshot/invalid", testCB, true, true, ""},
		{"3", "snapshot/testmd5", testCB, false, false, ""},
		{"4", "snapshot/testmd5", nil, true, false, ""},
	}

	var appv1 snapshot.UpdateServiceSnapshotAppv1
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(path), "testdata")

	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{CallbackID: c.id, DataURL: filepath.Join(dir, c.url), CallbackFunc: c.cb}
		a, _ := appv1.New(info)
		err := a.Process()
		assert.Equal(t, c.pExpected, err == nil, "Fail to get correct process output")

		data, ok := cbMap[c.id]
		assert.Equal(t, c.idExpected, ok, "Fail to call cb")
		if ok {
			assert.Equal(t, c.md5, fmt.Sprintf("%x", data.Data), "Fail to call cb md5")
		}
	}
}
