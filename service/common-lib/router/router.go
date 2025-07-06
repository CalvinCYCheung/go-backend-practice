package router

import "github.com/gin-gonic/gin"

func InitRouter(initRoute func(router *gin.RouterGroup), middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middlewares...)
	api := router.Group("/api")
	initRoute(api)
	return router
}
