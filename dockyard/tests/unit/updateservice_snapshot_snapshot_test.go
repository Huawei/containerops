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

	"github.com/Huawei/dockyard/updateservice/snapshot"
)

type UpdateServiceSnapshotMock struct {
}

func (m *UpdateServiceSnapshotMock) New(info snapshot.SnapshotInputInfo) (snapshot.UpdateServiceSnapshot, error) {
	return m, nil
}

func (m *UpdateServiceSnapshotMock) Supported(proto string) bool {
	return proto == "mock"
}

func (m *UpdateServiceSnapshotMock) Process() error {
	return nil
}

func (m *UpdateServiceSnapshotMock) Description() string {
	return "mock description"
}

// expunge all the registed implementaions
func preTest() {
	cases := []struct {
		name string
		f    snapshot.UpdateServiceSnapshot
	}{
		{"mname0", &UpdateServiceSnapshotMock{}},
		{"mname1", &UpdateServiceSnapshotMock{}},
		{"aname0", &snapshot.UpdateServiceSnapshotAppv1{}},
		{"aname1", &snapshot.UpdateServiceSnapshotAppv1{}},
	}

	snapshot.UnregisterAllSnapshot()

	for _, c := range cases {
		snapshot.RegisterSnapshot(c.name, c.f)
	}
}

func TestRegisterSnapshot(t *testing.T) {
	preTest()

	cases := []struct {
		name     string
		f        snapshot.UpdateServiceSnapshot
		expected bool
	}{
		{"", &UpdateServiceSnapshotMock{}, false},
		{"testsname", nil, false},
		{"testsname", &UpdateServiceSnapshotMock{}, true},
		{"testsname", &UpdateServiceSnapshotMock{}, false},
	}

	for _, c := range cases {
		err := snapshot.RegisterSnapshot(c.name, c.f)
		assert.Equal(t, c.expected, err == nil, "Fail to register snapshot")
	}
}

func TestListSnapshot(t *testing.T) {
	preTest()

	strs := snapshot.ListSnapshotByProto("mock")
	assert.Equal(t, 2, len(strs), "Fail to get correct snapshot list")
}

func TestNewUpdateServiceSnapshot(t *testing.T) {
	preTest()

	cases := []struct {
		name     string
		expected bool
	}{
		{"mname0", true},
		{"invalidname", false},
	}

	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{Name: c.name}
		_, err := snapshot.NewUpdateServiceSnapshot(info)
		assert.Equal(t, c.expected, err == nil, "Fail to create new snapshot")
	}
}

func TestIsSnapshotSupported(t *testing.T) {
	preTest()

	cases := []struct {
		p        string
		n        string
		expected bool
	}{
		{"mock", "mname0", true},
		{"mock", "invalid", false},
		{"invalid", "mname0", false},
	}

	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{DataProto: c.p, Name: c.n}
		ok, _ := snapshot.IsSnapshotSupported(info)
		assert.Equal(t, c.expected, ok, "Fail to get supported result")
	}
}

func TestInfoGetName(t *testing.T) {
	cases := []struct {
		name           string
		expectedPlugin string
		expectedImage  string
	}{
		{"appv1", "appv1", ""},
		{"bycontainer/ospaf/scan", "bycontainer", "ospaf/scan"},
	}

	for _, c := range cases {
		info := snapshot.SnapshotInputInfo{Name: c.name}
		p, i := info.GetName()
		assert.Equal(t, c.expectedPlugin, p, "Fail to get plugin name")
		assert.Equal(t, c.expectedImage, i, "Fail to get image name")
	}
}
