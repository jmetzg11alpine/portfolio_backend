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
		govRouters.POST("/foreign-aid-map", controllers.GetForeignAidMapHandler)
		govRouters.POST("/foreign-aid-bar", controllers.GetForeignAidBarHandler)
		govRouters.GET("/comparison", controllers.GetComparisonHandler)
	}
	etfRouters := router.Group("/etf")
	{
		etfRouters.GET("/reserves", controllers.GetETFReserveData)
		etfRouters.POST("/stocks", controllers.GetStockData)
	}
}
