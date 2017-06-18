/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

package router

import (
	"gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/dockyard/handler"
)

// SetRouters is setting REST API interface with handler function.
func SetRouters(m *macaron.Macaron) {
	// Create Repository
	m.Group("/v1", func() {
		m.Post("/:namespace/:repository/:type", handler.PostRepositoryV1Handler)
	})

	// Docker Registry V2
	m.Group("/v2", func() {
		m.Get("/", handler.GetPingV2Handler)
		m.Get("/_catalog", handler.GetCatalogV2Handler)

		// User mode: /namespace/repository:tag
		m.Head("/:namespace/:repository/blobs/:digest", handler.HeadBlobsV2Handler)
		m.Post("/:namespace/:repository/blobs/uploads", handler.PostBlobsV2Handler)
		m.Patch("/:namespace/:repository/blobs/uploads/:uuid", handler.PatchBlobsV2Handler)
		m.Put("/:namespace/:repository/blobs/uploads/:uuid", handler.PutBlobsV2Handler)
		m.Get("/:namespace/:repository/blobs/:digest", handler.GetBlobsV2Handler)
		m.Put("/:namespace/:repository/manifests/:tag", handler.PutManifestsV2Handler)
		m.Get("/:namespace/:repository/tags/list", handler.GetTagsListV2Handler)
		m.Get("/:namespace/:repository/manifests/:tag", handler.GetManifestsV2Handler)
		m.Delete("/:namespace/:repository/blobs/:digest", handler.DeleteBlobsV2Handler)
		m.Delete("/:namespace/:repository/:blobs/:uuid", handler.DeleteBlobsUUIDV2Handler)
		m.Delete("/:namespace/:repository/manifests/:reference", handler.DeleteManifestsV2Handler)

		// Library mode: /repository:tag
		m.Get("/:repository/blobs/:digest", handler.GetBlobsV2LibraryHandler)
		m.Get("/:repository/tags/list", handler.GetTagsListV2LibraryHandler)
		m.Get("/:repository/manifests/:tag", handler.GetManifestsV2LibraryHandler)
	})

	// Binary File
	m.Group("/Binary", func() {
		// V1 Version
		m.Group("/v1", func() {
			m.Post("/:namespace/:repository/binary/:binary/:tag", handler.PostBinaryV1Handler)
			m.Get("/:namespace/:repository/binary/:binary/:tag", handler.GetBinaryV1Handler)
			m.Put("/:namespace/:repository/binary/:binary/:tag/:label", handler.PutBinaryLabelV1Handler)
			m.Delete("/:namespace/:repository/binary/:binary/:tag", handler.DeleteBinaryV1Handler)
		})
	})

}
