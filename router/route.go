package router

import (
	"demo/controller"
	"demo/middlerware"

	"github.com/gin-gonic/gin"
)

// API ...
func API(eng *gin.Engine) {

	router := eng.Group("")

	// middleware
	{
		router.Use(middlerware.Session())
	}

	// 登录相关

	{
		router.POST("/login", controller.Login)
		router.GET("/logout", controller.Logout)
	}

	// 业务相关

	login := router.Group("")
	login.Use(middlerware.SessionCheck)
	{
		login.GET("/user", controller.User)
	}
}
