package controllers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {
	message := models.GetHomeMessage()
	c.JSON(http.StatusOK, message)
}
