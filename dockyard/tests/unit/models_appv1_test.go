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

	"github.com/Huawei/dockyard/models"
)

// TestNewAppV1
func TestNewAppV1(t *testing.T) {
	namespace := "containerops"
	repository := "official"
	_, err := models.NewAppV1(namespace, repository)
	assert.Nil(t, err, "Fail to create/query a repository")
}

// TestAppV1Put
func TestAppV1Put(t *testing.T) {
	namespace := "containerops"
	repository := "official"
	a := models.ArtifactV1{}

	r, _ := models.NewAppV1(namespace, repository)
	err := r.Put(a)
	assert.Nil(t, err, "Fail to add an artifact to a repository")

	r.Locked = true
	err = r.Put(a)
	assert.NotNil(t, err, "Should not add an artifact to a locked repository")
}

// TestArtifactV1GetName
func TestArtifactV1GetName(t *testing.T) {
	cases := []struct {
		a        models.ArtifactV1
		expected string
	}{
		{a: models.ArtifactV1{OS: "os", Arch: "arch", App: "app", Tag: "tag"}, expected: "os/arch/app:tag"},
		{a: models.ArtifactV1{OS: "os", Arch: "arch", App: "app", Tag: ""}, expected: "os/arch/app"},
		{a: models.ArtifactV1{OS: "os", Arch: "arch", App: "", Tag: "tag"}, expected: ""},
		{a: models.ArtifactV1{OS: "os", Arch: "", App: "app", Tag: "tag"}, expected: ""},
		{a: models.ArtifactV1{OS: "", Arch: "arch", App: "app", Tag: "tag"}, expected: ""},
	}

	for _, c := range cases {
		assert.Equal(t, c.a.GetName(), c.expected, "Fail to get artifact name")
	}
}
