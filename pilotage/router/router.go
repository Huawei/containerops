package router

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/containerops/pilotage/handler"
)

//SetRouters is pilotage router's definition fucntion.
func SetRouters(m *macaron.Macaron) {
	// Web API
	m.Get("/", handler.IndexV1Handler)

	m.Group("/pipeline", func() {
		m.Group("/v1", func() {
			//Definie the supported service.
			m.Group("/service", func() {
				m.Post("/", binding.Bind(handler.PostServiceDefinitionForm{}), handler.PostServiceDefinitionV1Handler)
				m.Get("/list", handler.GetServiceDefinitionListV1Handler)

				m.Group("/:service", func() {
					m.Put("/", handler.PutServiceDefinitionV1Handler)
					m.Get("/", handler.GetServiceDefinitionV1Handler)
					m.Delete("/", handler.DeleteServiceDefinitionV1Handler)
				})
			})

			//Registry the component in the system.
			m.Group("/:namespace/component", func() {
				m.Post("/", handler.PostComponentV1Handler)

				m.Group("/:component", func() {
					m.Put("/", handler.PutComponentV1Handler)
					m.Get("/", handler.GetComponentV1Handler)
					m.Delete("/", handler.DeleteComponentv1Handler)

					//Define the events of component
					m.Group("/event", func() {
						m.Post("/", handler.PostEventV1Handler)

						m.Group("/:evnet", func() {
							m.Get("/", handler.GetEventV1Handler)
							m.Put("/", handler.PutEventV1Handler)
							m.Delete("/", handler.DeleteEventV1Handler)
						})
					})
				})
			})

			//CRUD of pipeline.
			m.Group("/:namespace/:repository", func() {
				//Define pipeline
				m.Get("/", handler.GetPipelineListV1Handler)
				m.Post("/", handler.PostPipelineV1Handler)
				m.Post("/json", handler.PostPipelineJSONV1Handler)

				m.Group("/:pipeline", func() {
					//Get/Put/Delete Pipeline
					m.Get("/:format", handler.GetPipelineV1Handler)
					m.Put("/", handler.PutPipelineV1Handler)
					m.Delete("/", handler.DeletePipelineV1Handler)

					//Definie the stage
					m.Group("/stage", func() {
						m.Post("/", handler.PostStageV1Handler)

						m.Group("/:stage", func() {
							m.Get("/", handler.GetStageV1Handler)
							m.Put("/", handler.PutStageV1Handler)
							m.Delete("/", handler.DeleteStageV1Handler)

							m.Post("/action", handler.PostActionV1Handler)
							m.Group("/:action", func() {
								m.Get("/", handler.GetActionV1Handler)
								m.Put("/", handler.PutActionV1Handler)
								m.Delete("/", handler.DeleteActionV1Handler)

								//Binding the service supported with User/Organization
								m.Group("/service", func() {
									m.Post("/", handler.PostServiceV1Handler)

									//When call service with ?sequence=xxx param
									m.Group("/:service", func() {
										m.Put("/", handler.PutServiceV1Handler)
										m.Get("/", handler.GetServiceV1Handler)
										m.Delete("/", handler.DeleteServiceV1Handler)
										m.Any("/callback", handler.AnyServiceCallbackV1Handler) //The callback must have ?sequence=xxx

										//Define the events of service
										m.Group("/event", func() {
											m.Post("/", handler.PostEventV1Handler)

											m.Group("/:evnet", func() {
												m.Get("/", handler.GetEventV1Handler)
												m.Put("/", handler.PutEventV1Handler)
												m.Delete("/", handler.DeleteEventV1Handler)
											})
										})
									})
								})
							})
						})
					})

					//Run a pipeline with sequence id
					m.Post("/", handler.ExecutePipelineV1Handler)

					// Callback of all action
					m.Put("/event", handler.PutActionEventV1Handler)

					// all action register here
					// pipeline will push data to the url which is send by action on register
					m.Put("/register", handler.PutActionRegisterV1Handler)

					m.Group("/:sequence", func() {
						m.Get("/outcome/list", handler.GetOutcomeListV1Handler)
						m.Get("/:outcome", handler.GetOutcomeV1Handler)

						//CRUD environment in the running pipeline.
						m.Post("/env", handler.PostEnvironmentV1Handler)
						m.Get("/env/list", handler.GetEnvironmentListV1Hander)

						m.Group("/:env", func() {
							m.Get("/", handler.GetEnvironmentV1Handler)
							m.Put("/", handler.PutEnvironmentV1Handler)
							m.Delete("/", handler.DeleteEnvironmentV1Handler)
						})

						// //Callbacks of action for component.
						// //Calllback URL initlization when the action run with sequence param.
						// m.Group("/:stage/:action", func() {
						// 	m.Put("/start", handler.PutStartActionV1Handler)
						// 	m.Put("/execute", handler.PutExecuteActionV1Handler)
						// 	m.Put("/status", handler.PutStatusActionV1Handler)
						// 	m.Put("/result/:result", handler.PutResultActionV1Handler)
						// 	m.Put("/delete", handler.PutDeleteActionV1Handle)
						// })
					})
				})
			})
		})
	})
}
