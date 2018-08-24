package main

import (
	"fmt"
	"github.com/maxim-kuderko/cloud-logger/handlers"
	"github.com/maxim-kuderko/cloud-logger/middlewares"
	"github.com/maxim-kuderko/cloud-logger/registry"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done)

	reg := registry.NewRegistry()
	routes := defineRoutes(reg)
	middleware := middlewares.NewMiddlwares(middlewares.Logging)
	go fasthttp.ListenAndServe(":8000", middleware.Then(routes))

	<-done
	fmt.Println("exiting")
}

func defineRoutes(registry *registry.Registry) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/is_alive":
			handlerWrapper(handlers.IsAlive, registry)(ctx)
		case "/push":
			handlerWrapper(handlers.Push, registry)(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
}

func handlerWrapper(handler func(ctx *fasthttp.RequestCtx, registry *registry.Registry), registry *registry.Registry) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		handler(ctx, registry)
	}
}
