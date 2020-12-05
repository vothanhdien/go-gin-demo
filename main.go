package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-gin-demo/modal"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/ping", func(c *gin.Context) {
		var r = &modal.PingRequest{}
		if err := c.MustBindWith(r,binding.Default(c.Request.Method, c.ContentType())); err != nil {
			_ = c.Error(err)
			_ = c.AbortWithError(http.StatusOK, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "hello " + r.Name,
		})
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
