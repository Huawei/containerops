package web

import (
	"github.com/Huawei/containerops/assembling/handler"
	macaron "gopkg.in/macaron.v1"
)

func SetAssemblingMacaron(m *macaron.Macaron, config string) {
	//Setting Config
	// common.SetConfig(config)

	//Setting Middleware
	// middleware.SetMiddlewares(m)

	//Setting Router
	// router.SetRouters(m)

	m.Group("/assembling", func() {
		m.Post("/build", handler.BuildImageHandler)
	})
}
