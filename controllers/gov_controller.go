package controllers

import (
	"backend/models/gov"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAgencyHandler(c *gin.Context) {
	data, err := gov.GetAgencyData()
	if err != nil {
		log.Printf("error calling GetAgencyHandler: %v\n", err)
	}
	c.JSON(http.StatusOK, data)
}
