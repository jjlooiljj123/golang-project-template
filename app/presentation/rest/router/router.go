package router

import (
	"boilerplate/app/infrastructure/config"
	restcontroller "boilerplate/app/presentation/rest/album"
	"boilerplate/app/presentation/rest/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, controller *restcontroller.Controller, cfg *config.AppConfig) {

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.LatencyLogger())
	router.Use(middleware.TimeoutMiddleware(cfg))
	router.Use(middleware.CommonHeadersMiddleware())

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		v1.GET("/albums", controller.GetAlbumsHandler)
		v1.POST("/albums", controller.CreateAlbumHandler)
		v1.GET("/albums/:id", controller.GetAlbumByIDHandler)

		// Apply auth middleware to these routes
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/jsonposts", controller.GetJsonPostHandler)
		}
	}
	{
		v2 := api.Group("/v2")
		v2.GET("/albums", controller.GetAlbumsHandler)
	}
}
