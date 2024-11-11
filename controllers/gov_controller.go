package controllers

import (
	"backend/models/gov"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAgencyHandler(c *gin.Context) {
	mainData, otherData, remainingData, err := gov.GetAgencyData()
	if err != nil {
		log.Printf("error calling GetAgencyHandler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get agency data"})
		return
	}
	response := gin.H{
		"mainData":      mainData,
		"otherData":     otherData,
		"remainingData": remainingData,
	}
	c.JSON(http.StatusOK, response)
}
