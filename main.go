package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var userManagementServer = "http//test:3250"

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	log.Println("Hello World")

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		loginHandles := v1.Group("/user")
		loginHandles.POST("/login", loginHandl)

		signupHandles := v1.Group("/signup")
		signupHandles.POST("/participant", signupParticipantHandl)
		signupHandles.POST("/researcher", signupResearcherHandl)
		signupHandles.POST("/admin", signupAdminHandl)

		tokenHandles := v1.Group("/token")
		tokenHandles.GET("/validate", validateTokenHandl)
		tokenHandles.GET("/renew", renewTokenHandl)
	}

	log.Fatal(router.Run(":3100"))
}
