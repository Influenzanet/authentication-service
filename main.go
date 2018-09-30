package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	log.Println("Hello World")

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginParticipantHandl)
		v1.GET("/validate", validateTokenHandl)
	}

	log.Fatal(router.Run(":3100"))
}
