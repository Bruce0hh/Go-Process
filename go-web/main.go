package main

import (
	"goweb/web"
	"log"
	"net/http"
	"time"
)

func main() {
	r := web.Default()

	r.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World!</h1>", nil)
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *web.Context) {
			c.HTML(http.StatusOK, "<h1>Hello World!</h1>", nil)
		})
		v1.GET("/hello", func(c *web.Context) {
			c.String(http.StatusOK, "hello", c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *web.Context) {
			c.String(http.StatusOK, "hello %v, you're at %v\n", c.Query("name"), c.Path)
		})
		v2.POST("/login", func(c *web.Context) {
			c.JSON(http.StatusOK, web.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.GET("/hello", func(c *web.Context) {
		c.String(http.StatusOK, "hello %v, you're at %v\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *web.Context) {
		c.String(http.StatusOK, "hello %v, you're at %v\n", c.Query("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *web.Context) {
		c.JSON(http.StatusOK, web.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	v3 := r.Group("/v3")
	v3.Use(onlyForV3())
	{
		v3.GET("/hello/:name", func(ctx *web.Context) {
			ctx.String(http.StatusOK, "hello %v, you're at %v\n", ctx.Param("name"), ctx.Path)
		})
	}

	r.GET("/panic", func(ctx *web.Context) {
		name := []string{"myName"}
		ctx.String(http.StatusOK, name[10])
	})

	_ = r.Run(":9999")
}

func onlyForV3() web.HandlerFunc {
	return func(ctx *web.Context) {
		t := time.Now()
		ctx.String(500, "Internal Server Error\n")
		log.Printf("[%d] %s in %v for group v3\n", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
