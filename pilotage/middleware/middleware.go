package middleware

import (
	"gopkg.in/macaron.v1"
)

func SetMiddlewares(m *macaron.Macaron) {
	//Set static file directory,static file access without log output
	m.Use(macaron.Static("external", macaron.StaticOptions{
		Expires: func() string { return "max-age=0" },
	}))

	//Set recovery handler to returns a middleware that recovers from any panics
	m.Use(macaron.Recovery())

	m.Use(func(ctx *macaron.Context) {
		ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*")
	})
}
