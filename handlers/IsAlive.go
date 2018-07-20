package handlers

import (
	"github.com/maxim-kuderko/cloud-logger/registry"
	"github.com/valyala/fasthttp"
)

func IsAlive(ctx *fasthttp.RequestCtx, registry *registry.Registry) {
	ctx.SetBody([]byte("ok"))
}
