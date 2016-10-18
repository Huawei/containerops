package web

import (
	"gopkg.in/macaron.v1"

	"github.com/containerops/pilotage/middleware"
	"github.com/containerops/pilotage/router"
)

//
func SetPilotagedMacaron(m *macaron.Macaron) {
	//Setting Middleware
	middleware.SetMiddlewares(m)

	//Setting Router
	router.SetRouters(m)
}
