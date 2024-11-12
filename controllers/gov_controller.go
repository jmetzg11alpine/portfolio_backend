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

type ForeignAidRequest struct {
	Country string `json:"country" binding:"required"`
	Year    string `json:"year" binding:"required"`
}

func GetForeignAidHandler(c *gin.Context) {
	var request ForeignAidRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := gov.GetForeignAidData(request.Country, request.Year)
	response := gin.H{
		"data": data,
	}
	c.JSON(http.StatusOK, response)
}
