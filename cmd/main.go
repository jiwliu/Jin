package main

import (
	"jgweb/pkg/Jin"
	"log"
	"net/http"
	"time"
)

func main() {
	engine := Jin.New()
	engine.GET("/", func(context *Jin.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	engine.GET("/hello", func(context *Jin.Context) {
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Query("name"), context.Path)
	})
	engine.GET("/login", func(context *Jin.Context) {
		context.JSON(http.StatusOK, Jin.H{
			"user:": context.PostForm("user"),
			"paas:": context.PostForm("paas"),
		})
	})
	v1 := engine.Group("/v1")
	v1.Use(func(context *Jin.Context) {
		t := time.Now()
		log.Println("before...........")
		context.Next()
		log.Println("after...........")
		log.Printf("[%d] %s in %v for group v2", context.StatusCode, context.Request.RequestURI, time.Since(t))
	})
	v1.GET("/hello", func(context *Jin.Context) {
		log.Println("hello doing....")
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Query("name"), context.Path)
	})
	v1.GET("/panic", func(c *Jin.Context) {
		names := []string{"Jin"}
		c.String(http.StatusOK, names[100])
	})
	//---------------------
	//engine := Jin.DEFAULT()
	//engine.GET("/hello", func(context *Jin.Context) {
	//	context.String(http.StatusOK, "hello %s, you're at %s\n", context.Query("name"), context.Path)
	//})
	//engine.GET("/panic", func(c *Jin.Context) {
	//	names := []string{"Jin"}
	//	c.String(http.StatusOK, names[100])
	//})
	//engine.Run(":9999")

}
