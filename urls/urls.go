package urls

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	homeRoutes := router.Group("/")
	{
		homeRoutes.GET("", controllers.GetHome)
	}
	govRouters := router.Group("/gov")
	{
		govRouters.GET("/agency", controllers.GetAgencyHandler)
		govRouters.GET("/foreign-aid", controllers.GetForeignAidHandler)
	}
}
