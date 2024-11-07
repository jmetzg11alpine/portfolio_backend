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
}

// make different beging endpoints:
// routes/url.go
// package routes

// import (
//     "my-backend/controllers"
//     "github.com/gin-gonic/gin"
// )

// func InitializeRoutes(router *gin.Engine) {
//     itemRoutes := router.Group("/items")
//     {
//         itemRoutes.GET("/", controllers.GetItems)
//         itemRoutes.POST("/", controllers.CreateItem)
//     }
// }
