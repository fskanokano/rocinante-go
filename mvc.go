package rocinante

import (
	"reflect"
)

type ControllerInterface interface {
	GET(*Context)
	POST(*Context)
	PUT(*Context)
	DELETE(*Context)
	PATCH(*Context)
	OPTIONS(*Context)
	HEAD(*Context)
}

type Controller struct {
}

func (c *Controller) GET(ctx *Context) {
	MethodNotAllowedHandler(ctx)
}

func (c *Controller) POST(ctx *Context) {
	MethodNotAllowedHandler(ctx)
}

func (c *Controller) PUT(ctx *Context) {
	MethodNotAllowedHandler(ctx)
}

func (c *Controller) DELETE(ctx *Context) {
	MethodNotAllowedHandler(ctx)
}

func (c *Controller) PATCH(ctx *Context) {
	MethodNotAllowedHandler(ctx)
}

func (c *Controller) OPTIONS(ctx *Context) {
	OptionsHandler(ctx)
}

func (c *Controller) HEAD(ctx *Context) {
	MethodNotAllowedHandler(ctx)
}

func dispatchHandler(relativePath string, c ControllerInterface, r *RouterGroup) {
	for i := 0; i < len(Methods); i++ {
		method := Methods[i]

		cv := reflect.ValueOf(c)
		handlerV := cv.MethodByName(method)

		handler := func(c *Context) {
			handlerV.Call([]reflect.Value{
				reflect.ValueOf(c),
			})
		}

		r.addRoute(relativePath, method, true, handler)
	}
}

var Methods = []string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"PATCH",
	"OPTIONS",
	"HEAD",
}
