package main

import (
	"goweb/web"
	"net/http"
)

func main() {
	r := web.New()
	r.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World!</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *web.Context) {
			c.HTML(http.StatusOK, "<h1>Hello World!</h1>")
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

	_ = r.Run(":9999")
}
