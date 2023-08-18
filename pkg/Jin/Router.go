package Jin

import (
	"log"
	"net/http"
)

type Handler func(context *Context)

func newRouter() *Router {
	return &Router{handlers: make(map[string]Handler)}
}

func (router *Router) addRoute(method string, pattern string, handler Handler) {
	log.Printf("router %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *Router) handle(context *Context) {
	key := context.Method + "-" + context.Path
	if handler, ok := router.handlers[key]; ok {
		// 当前handler放在最尾端被调用
		context.middlewares = append(context.middlewares, handler)
	} else {
		context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
	}
	context.Next()
}
