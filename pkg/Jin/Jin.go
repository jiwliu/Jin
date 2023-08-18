package Jin

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type Router struct {
	handlers map[string]Handler
}

type RouterGroup struct {
	prefix      string
	middlewares []Handler
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func DEFAULT() *Engine {
	engine := New()
	engine.Use(Logging(), Recovery())
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (routeGroup *RouterGroup) addRoute(method string, pattern string, handler Handler) {
	routeGroup.engine.router.addRoute(method, routeGroup.prefix+pattern, handler)
}

func (routeGroup *RouterGroup) GET(pattern string, handler Handler) {
	routeGroup.addRoute("GET", pattern, handler)
}

func (routeGroup *RouterGroup) POST(pattern string, handler Handler) {
	routeGroup.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (routeGroup *RouterGroup) Use(handler ...Handler) {
	routeGroup.middlewares = append(routeGroup.middlewares, handler...)
}

// 中间件支持
// 返回结果，json&html等各种类型支持，header封装
// 异常处理
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []Handler
	context := newContext(writer, request)
	for _, group := range engine.groups {
		if strings.HasPrefix(context.Path, group.prefix) {
			middlewares = append(context.middlewares, group.middlewares...)
		}
	}
	context.middlewares = middlewares
	engine.router.handle(context)
}

func Logging() Handler {
	return func(context *Context) {
		t := time.Now()
		context.Next()
		log.Printf("[%d] %s cost %v", context.StatusCode, context.Request.RequestURI, time.Since(t))
	}
}

func Recovery() Handler {
	log.Println("Recovery...............................")
	return func(context *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				context.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		context.Next()
	}

}

func trace(message string) any {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
