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

package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	logs "github.com/Huawei/containerops/component/log"
	"github.com/containerops/configure"

	"gopkg.in/macaron.v1"
)

//  log is log pkg instance
var log *logs.Logger

// init is init log pkg instance
func init() {
	log = logs.New()
}

// SetMiddlewares is
func SetMiddlewares(m *macaron.Macaron) {
	//Set static file directory,static file access without log output
	m.Use(macaron.Static("external", macaron.StaticOptions{
		Expires: func() string { return "max-age=0" },
	}))

	//Set recovery handler to returns a middleware that recovers from any panics
	m.Use(macaron.Recovery())

	m.Use(func(ctx *macaron.Context) {
		ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Resp.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE")
	})

	m.Use(func(ctx *macaron.Context) {
		if ctx.Req.Method == "OPTIONS" {
			ctx.Resp.Write([]byte("success"))
			ctx.Resp.Flush()
		}
	})

	if configure.GetBool("log.showReqInfo") {
		m.Use(func(ctx *macaron.Context) {
			headerStr := ""
			for key, value := range ctx.Req.Header {
				headerStr += "\t" + key + ":" + strings.Join(value, " ") + "\n"
			}

			body := gotBody(ctx, 10*1024*1024)
			bodyStr := "\t" + strings.Replace(string(body), "\n", "\n\t", -1)
			log.Debugf("got request from:%9s ;url:%s\nheader is:\n%sbody is:\n%s", ctx.Req.Host, ctx.Req.RequestURI, headerStr, bodyStr)
		})
	}
}

func gotBody(ctx *macaron.Context, MaxMemory int64) []byte {
	if ctx.Req.Request.Body == nil {
		return []byte{}
	}
	safe := &io.LimitedReader{R: ctx.Req.Request.Body, N: MaxMemory}
	requestbody, _ := ioutil.ReadAll(safe)
	ctx.Req.Request.Body.Close()
	bf := bytes.NewBuffer(requestbody)
	ctx.Req.Request.Body = ioutil.NopCloser(bf)
	return requestbody
}
