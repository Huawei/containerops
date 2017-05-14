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
	"crypto/md5"
	"errors"
	"io/ioutil"
)

var (
	snapshotName   = "simpleAppV1"
	snapshotProtos = []string{"appv1"}
)

type UpdateServiceSnapshotAppv1 struct {
	info SnapshotInputInfo
}

func init() {
	RegisterSnapshot(snapshotName, &UpdateServiceSnapshotAppv1{})
}

func (m *UpdateServiceSnapshotAppv1) New(info SnapshotInputInfo) (UpdateServiceSnapshot, error) {
	if info.CallbackID == "" || info.DataURL == "" {
		return nil, errors.New("'CallbackID', 'DataURL' should not be empty")
	}

	m.info = info
	return m, nil
}

func (m *UpdateServiceSnapshotAppv1) Supported(proto string) bool {
	for _, p := range snapshotProtos {
		if p == proto {
			return true
		}
	}

	return false
}

func (m *UpdateServiceSnapshotAppv1) Process() error {
	var data SnapshotOutputInfo

	content, err := ioutil.ReadFile(m.info.DataURL)
	if m.info.CallbackFunc == nil {
		return err
	}

	if err == nil {
		s := md5.Sum(content)
		data.Data = s[:]
	}
	data.Error = err

	return m.info.CallbackFunc(m.info.CallbackID, data)
}

func (m *UpdateServiceSnapshotAppv1) Description() string {
	return "This is a simple snapshot. Scan the appv1 package, return its md5."
}
