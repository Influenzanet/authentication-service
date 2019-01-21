package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var userManagementServer = "http://user-management:3200/v1"

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	log.Println("Hello World")

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		userHandles := v1.Group("/user")
		userHandles.Use(RequirePayload())
		userHandles.POST("/login", loginHandl)
		userHandles.POST("/signup", signupHandl)

		tokenHandles := v1.Group("/token")
		tokenHandles.Use(ExtractToken())
		tokenHandles.GET("/validate", validateTokenHandl)
		tokenHandles.GET("/renew", renewTokenHandl)
	}

	log.Fatal(router.Run(":3100"))
}
