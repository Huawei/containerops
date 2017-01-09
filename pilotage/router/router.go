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
	"github.com/Huawei/containerops/pilotage/handler"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

//SetRouters is pilotage router's definition function.
func SetRouters(m *macaron.Macaron) {
	m.Group("/v2", func() {
		m.Get("/", handler.IndexV1Handler)

		m.Group("/events", func() {
			m.Post("/", handler.CreateEvent)
		})

		m.Group("/components", func() {
			m.Get("/", handler.ListComponents)

			m.Post("/", handler.CreateComponent)
			m.Get("/:component", handler.GetComponent)
			m.Put("/:component", handler.UpdateComponent)
			m.Delete("/:component", handler.DeleteComponent)

			m.Get("/:component/debug", handler.DebugComponentJson(), handler.DebugComponentLog)
		})

		m.Group("/images", func() {
			m.Post("/build", handler.BuildImage)
		})

		m.Group("/:namespace", func() {

			m.Group("/service", func() {
				m.Get("/list", handler.GetServiceDefinitionListV1Handler)
				m.Post("/", binding.Bind(handler.PostServiceDefinitionForm{}), handler.PostServiceDefinitionV1Handler)

				m.Get("/:service", handler.GetServiceDefinitionV1Handler)
				m.Put("/:service", handler.PutServiceDefinitionV1Handler)
				m.Delete("/:service", handler.DeleteServiceDefinitionV1Handler)
			})

			m.Group("/:repository", func() {

				m.Get("/system/v1/setting", handler.GetSettingV1Handler)
				m.Put("/system/v1/setting", handler.PutSettingV1Handler)

				m.Group("/workflow", func() {
					m.Group("/v1", func() {

						m.Group("/define", func() {
							m.Get("/list", handler.GetWorkflowListV1Handler)
							m.Post("/", handler.PostWorkflowV1Handler)

							m.Get("/event/:site/:event", handler.GetEventDefineJsonV1Handler)

							m.Get("/:workflow", handler.GetWorkflowV1Handler)
							m.Put("/:workflow", handler.PutWorkflowV1Handler)
							m.Delete("/:workflow", handler.DeleteWorkflowV1Handler)

							m.Get("/:workflow/token", handler.GetWorkflowTokenV1Handler)

							m.Get("/:workflow/env", handler.GetWorkflowEnvV1Handler)
							m.Put("/:workflow/env", handler.PutWorkflowEnvV1Handler)

							m.Get("/:workflow/var", handler.GetWorkflowVarV1Handler)
							m.Put("/:workflow/var", handler.PutWorkflowVarV1Handler)

							m.Put("/:workflow/state", handler.PutWorkflowStateV1Handler)
						})

						m.Post("/exec/:workflow", handler.ExecuteWorkflowV1Handler)

						m.Group("/runtime", func() {
							//m.Post("/event/:workflow/register", handler.PostActionRegisterV1Handler)
							//m.Post("/event/:workflow/:event", handler.PostActionEventV1Handler)

							m.Post("/var/:workflow", handler.PostActionSetVarV1Handler)

							m.Post("/linkstart/:workflow/:target", handler.PostActionLinkStartV1Handler)
						})

						m.Group("/history", func() {
							m.Get("/workflow/list", handler.GetWorkflowHistoriesV1Handler)
							m.Get("/workflow/:workflow/version/list", handler.GetWorkflowVersionHistoriesV1Handler)
							m.Get("/workflow/:workflow/version/:version/list", handler.GetWorkflowSequenceHistoriesV1Handler)
							m.Get("/workflow/:workflow/version/:version/sequence/:sequence/action/:action/linkstart/list", handler.GetActionLinkstartListV1Handler)

							m.Get("/:workflow/:version", handler.GetWorkflowHistoryDefineV1Handler)
							m.Get("/:workflow/:version/:sequence/stage/:stage", handler.GetStageHistoryInfoV1Handler)
							m.Get("/:workflow/:version/:sequence/stage/:stage/action/:action", handler.GetActionHistoryInfoV1Handler)
							m.Get("/:workflow/:version/:sequence/stage/:stage/action/:action/console/log", handler.GetActionConsoleLogV1Handler)
							m.Get("/:workflow/:version/:sequence/:relation", handler.GetSequenceLineHistoryV1Handler)
						})
					})
				})
			})
		})
	})
}
