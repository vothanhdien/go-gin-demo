package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type PingRequest struct {
	Name string `json:"name"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/ping", func(c *gin.Context) {
		var req = &PingRequest{}
		if err := c.MustBindWith(req,binding.Default(c.Request.Method, c.ContentType())); err != nil {
			_ = c.AbortWithError(http.StatusOK, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "hello " + req.Name,
		})
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
