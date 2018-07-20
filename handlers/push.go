package handlers

import (
	"github.com/maxim-kuderko/cloud-logger/registry"
	"github.com/valyala/fasthttp"
)

func PushPost(ctx *fasthttp.RequestCtx, registry *registry.Registry) {
	ctx.SetBody([]byte("ok"))
	registry.Ds.Write()
}

func PushGet(ctx *fasthttp.RequestCtx, registry *registry.Registry) {
	ctx.SetBody([]byte("ok"))
	registry.Ds.Write()
}
