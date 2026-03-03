package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main(){
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	r.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})

	_ = r.Run(":8080")
}