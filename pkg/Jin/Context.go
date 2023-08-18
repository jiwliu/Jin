package Jin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Path        string
	Method      string
	StatusCode  int
	middlewares []Handler
	index       int
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{Writer: writer, Request: request, Path: request.URL.Path, Method: request.Method, index: -1}
}

func (context *Context) PostForm(key string) string {
	return context.Request.FormValue(key)
}

func (context Context) Query(key string) string {
	return context.Request.URL.Query().Get(key)
}

func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

func (context *Context) Data(code int, byte []byte) {
	context.Status(code)
	context.Writer.Write(byte)
}

func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	context.Writer.Write([]byte(html))
}

func (context *Context) Next() {
	context.index++
	s := len(context.middlewares)
	for ; context.index < s; context.index++ {
		context.middlewares[context.index](context)
	}
}

func (context *Context) Fail(code int, str string) {
	context.String(code, "%s", str)
}
