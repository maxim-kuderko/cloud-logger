package handlers

import (
	"github.com/maxim-kuderko/cloud-logger/registry"
	"github.com/valyala/fasthttp"
	"github.com/maxim-kuderko/cloud-logger/models"
	"encoding/json"
	"fmt"
)

func Push(ctx *fasthttp.RequestCtx, registry *registry.Registry) {
	var body []byte
	if ctx.IsPost(){
		body = ctx.Request.Body()
	} else if ctx.IsGet(){
		body = ctx.RequestURI()
	} else{
		ctx.Response.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		return
	}
	topic, err := getTopic(ctx.Request.Header.Peek("topic"), body, ctx.QueryArgs().Peek(`topic`))
	if err != nil{
		r := models.Response{
			Success: false,
			Error: err.Error(),
		}
		b, e := json.Marshal(r)
		if e != nil{
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType(`application/json`)
		ctx.SetBody(b)
		return
	}
	if err := registry.Ds.Write(topic, ctx.Request.Header.Header(), body); err != nil{
		r := models.Response{
			Success: false,
			Error: err.Error(),
		}
		b, e := json.Marshal(r)
		if e != nil{
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType(`application/json`)
		ctx.SetBody(b)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}


func getTopic(header []byte, body []byte, uri []byte) (string, error){
	if len(header) > 0{
		return string(header), nil
	}
	if len(uri) > 0{
		return string(uri), nil
	}
	return "", fmt.Errorf("topic not found in request")
}