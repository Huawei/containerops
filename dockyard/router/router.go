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

package router

import (
	"gopkg.in/macaron.v1"

	"github.com/Huawei/dockyard/handler"
)

//SetRouters is setting REST API interface with handler function.
func SetRouters(m *macaron.Macaron) {
	//Web API
	m.Get("/", handler.GetIndexPageV1Handler)
	m.Get("/pubkeys", handler.GetGPGFileV1Handler)

	//REST API For Web Operations
	m.Group("/web", func() {

		m.Get("/:namespace", handler.GetNamespacePageV1Handler)
		m.Get("/:type/:namespace/:repository", handler.GetRepositoryPageV1Handler)
		m.Get("/:type/:namespace/:repository/:package", handler.GetPackagePageV1Handler)
		m.Get("/:type/:namespace/:repository/:package/manifest", handler.GetManifestPageV1Handler)

		m.Group("/v1", func() {
			m.Post("/:type/:namespace/:repository", handler.PostRepositoryRESTV1Handler)
			m.Get("/:type/:namespace/:repository", handler.GetRepositoryRESTV1Handler)
			m.Put("/:type/:namespace/:repository", handler.PutRepositoryRESTV1Handler)
			m.Delete("/:type/namespace/:repository", handler.DeleteRepositoryRESTV1Handler)

			m.Post("/:type/:namesapce/:repository/:package", handler.PostPackageRESTV1Handler)
			m.Get("/:type/:namespace/:repository/:package", handler.GetPackageRESTV1Hanfdler)
			m.Put("/:type/:namespace/:repository/:package", handler.PutPackageRESTV1Handler)
			m.Delete("/:type/:namespace/:repository/:pacakge", handler.DeletePacakgeRESTV1Handler)

			m.Post("/:type/:namespace/:repository/:package/manifest", handler.PostManifestRESTV1Handler)
			m.Get("/:type/:namespace/:repository/:package/manifest", handler.GetManifestRESTV1Handler)
			m.Put("/:type/:namespace/:repository/:pacakge/manifest", handler.PutManifestRESTV1Handler)
			m.Delete("/:type/:namespace/:repository/:pacakge/manifest", handler.DeleteManifestRESTV1Handler)
		})
	})

	//Docker Registry V1
	m.Group("/v1", func() {
		m.Get("/_ping", handler.GetPingV1Handler)

		m.Get("/users", handler.GetUsersV1Handler)
		m.Post("/users", handler.PostUsersV1Handler)

		m.Group("/repositories", func() {
			m.Put("/:namespace/:repository/tags/:tag", handler.PutTagV1Handler)
			m.Put("/:namespace/:repository/images", handler.PutRepositoryImagesV1Handler)
			m.Get("/:namespace/:repository/images", handler.GetRepositoryImagesV1Handler)
			m.Get("/:namespace/:repository/tags", handler.GetTagV1Handler)
			m.Put("/:namespace/:repository", handler.PutRepositoryV1Handler)
		})

		m.Group("/images", func() {
			m.Get("/:image/ancestry", handler.GetImageAncestryV1Handler)
			m.Get("/:image/json", handler.GetImageJSONV1Handler)
			m.Get("/:image/layer", handler.GetImageLayerV1Handler)
			m.Put("/:image/json", handler.PutImageJSONV1Handler)
			m.Put("/:image/layer", handler.PutImageLayerV1Handler)
			m.Put("/:image/checksum", handler.PutImageChecksumV1Handler)
		})
	})

	//Docker Registry V2
	m.Group("/v2", func() {
		m.Get("/", handler.GetPingV2Handler)
		m.Get("/_catalog", handler.GetCatalogV2Handler)

		//user mode: /namespace/repository:tag
		m.Head("/:namespace/:repository/blobs/:digest", handler.HeadBlobsV2Handler)
		m.Post("/:namespace/:repository/blobs/uploads", handler.PostBlobsV2Handler)
		m.Patch("/:namespace/:repository/blobs/uploads/:uuid", handler.PatchBlobsV2Handler)
		m.Put("/:namespace/:repository/blobs/uploads/:uuid", handler.PutBlobsV2Handler)
		m.Get("/:namespace/:repository/blobs/:digest", handler.GetBlobsV2Handler)
		m.Put("/:namespace/:repository/manifests/:tag", handler.PutManifestsV2Handler)
		m.Get("/:namespace/:repository/tags/list", handler.GetTagsListV2Handler)
		m.Get("/:namespace/:repository/manifests/:tag", handler.GetManifestsV2Handler)
		m.Delete("/:namespace/:repository/blobs/:digest", handler.DeleteBlobsV2Handler)
		m.Delete("/:namespace/:repository/:blobs/:uuid", handler.DeleteBlobsUUUIDV2Handler)
		m.Delete("/:namespace/:repository/manifests/:reference", handler.DeleteManifestsV2Handler)

		//library mode: /repository:tag
		m.Get("/:repository/blobs/:digest", handler.GetBlobsV2LibraryHandler)
		m.Get("/:repository/tags/list", handler.GetTagsListV2LibraryHandler)
		m.Get("/:repository/manifests/:tag", handler.GetManifestsV2LibraryHandler)
	})

	//App Discovery
	m.Group("/app", func() {
		m.Group("/v1", func() {
			//Global Search
			m.Get("/search", handler.AppGlobalSearchV1Handler)

			//Get public key
			m.Get("/:namespace/pubkey", handler.AppGetPublicKeyV1Handler)

			m.Group("/:namespace/:repository", func() {
				//Discovery
				m.Get("/?app-discovery=1", handler.AppDiscoveryV1Handler)

				//Scoped Search
				m.Get("/search", handler.AppScopedSearchV1Handler)
				m.Get("/list", handler.AppGetListAppV1Handler)

				//Pull
				m.Get("/meta", handler.AppGetMetaV1Handler)
				m.Get("/metasign", handler.AppGetMetaSignV1Handler)
				m.Get("/:os/:arch/:type/:app/?:tag", handler.AppGetFileV1Handler)
				m.Get("/:os/:arch/:type/:app/manifests/?:tag", handler.AppGetManifestsV1Handler)

				//Push
				m.Post("/", handler.AppPostFileV1Handler)
				m.Put("/:os/:arch/:type/:app/?:tag", handler.AppPutFileV1Handler)
				m.Put("/:os/:arch/:type/:app/manifests/?:tag", handler.AppPutManifestV1Handler)
				m.Patch("/:os/:arch/:type/:app/:status/?:tag", handler.AppPatchFileV1Handler)
				m.Delete("/:os/:arch/:type/:app/?:tag", handler.AppDeleteFileV1Handler)

				//Content Scan
				m.Post("/shook", handler.AppRegistScanHooksV1Handler)
				m.Post("/shook/:callbackID", handler.AppCallbackScanHooksV1Handler)
				m.Post("/:os/:arch/:app/shook/?:tag", handler.AppActiveScanHooksTaskV1Handler)
			})
		})
	})

	//Appc Discovery
	m.Group("/appc", func() {
		m.Group("/:namespace/:repository", func() {
			//Discovery
			m.Get("/?ac-discovery=1", handler.AppcDiscoveryV1Handler)

			//Pull
			m.Get("/fetch/:file", handler.AppcGetACIV1Handler)

			//Push
			m.Post("/push/:aci", handler.AppcPostACIV1Handler)
			m.Put("/push/:version/manifest/:aci", handler.AppcPutManifestV1Handler)
			m.Put("/push/:version/asc/:aci", handler.AppcPutASCV1Handler)
			m.Put("/push/:version/aci/:aci", handler.AppcPutACIV1Handler)
			m.Post("/push/:version/complete/:aci", handler.AppcPostCompleteV1Handler)
		})
	})

	//VM Image Discovery
	m.Group("/image", func() {
		m.Group("/v1", func() {
			//Global Search
			m.Get("/search", handler.ImageGlobalSearchV1Handler)

			m.Group("/:namespace/:repository", func() {
				//Discovery
				m.Get("/?image-discovery=1", handler.ImageDiscoveryV1Handler)

				//Scoped Search
				m.Get("/search", handler.ImageScopedSearchV1Handler)
				m.Get("/list", handler.ImageGetListV1Handler)

				//Pull
				m.Get("/:os/:arch/:image/?:tag", handler.ImageGetFileV1Handler)
				m.Get("/:os/:arch/:image/manifests/?:tag", handler.ImageGetManifestsV1Handler)

				//Push
				m.Post("/", handler.ImagePostV1Handler)
				m.Put("/:os/:arch/:image/?:tag", handler.ImagePutFileV1Handler)
				m.Put("/:os/:arch/:image/manifests/?:tag", handler.ImagePutManifestV1Handler)
				m.Patch("/:os/:arch/:image/:status/?:tag", handler.ImagePatchFileV1Handler)
				m.Delete("/:os/:arch/:image/?:tag", handler.ImageDeleteFileV1Handler)
			})
		})
	})

	//Sync APIS
	m.Group("/sync", func() {
		m.Group("/v1", func() {
			//Server Ping
			m.Get("/ping", handler.SyncGetPingV1Handler)

			m.Group("/master", func() {
				//Server Sync Of Master
				m.Post("/registry", handler.SyncMasterPostRegistryV1Handler)
				m.Delete("/registry", handler.SyncMasterDeleteRegistryV1Handler)

				m.Put("/mode", handler.SyncMasterPutModeRegistryV1Handler)
			})

			m.Group("/slave", func() {
				//Server Sync Of Slaver
				m.Post("/registry", handler.SyncSlavePostRegistryV1Handler)
				m.Put("/registry", handler.SyncSlavePutRegistryV1Handler)
				m.Delete("/registry", handler.SyncSlaveDeleteRegistryV1Handler)

				m.Put("/mode", handler.SyncSlavePutModeRegistryV1Handler)

				//Data Sync
				m.Get("/list", handler.SyncSlaveListDataV1Handler)

				//File Sync
				m.Put("/:namespace/:repository/manifests", handler.SyncSlavePutManifestsV1Handler)
				m.Put("/:namespace/:repository/file", handler.SyncSlavePutFileV1Handler)
				m.Put("/:namespace/:repository/:status", handler.SyncSlavePutStatusV1Handler)
			})
		})
	})

	//Admin APIs
	m.Group("/admin", func() {
		m.Group("/v1", func() {
			//Server Status
			m.Get("/stats/:type", handler.AdminGetStatusV1Handler)

			//Server Config
			m.Get("/config", handler.AdminGetConfigV1Handler)
			m.Put("/config", handler.AdminSetConfigV1Handler)

			//Maintenance
			m.Post("/maintenance", handler.AdminPostMaintenance)
		})
	})
}
