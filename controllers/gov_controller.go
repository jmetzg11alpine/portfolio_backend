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

type ForeignAidMapRequest struct {
	Year string `json:"year" binding:"required"`
}

func GetForeignAidMapHandler(c *gin.Context) {
	var request ForeignAidMapRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mapData, countries := gov.GetForeignAidMapData(request.Year)
	response := gin.H{
		"mapData":   mapData,
		"countries": countries,
	}
	c.JSON(http.StatusOK, response)
}

type ForeignAidCountryRequest struct {
	Country string `json:"country" binding:"required"`
}

func GetForeignAidBarHandler(c *gin.Context) {
	var request ForeignAidCountryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	barData := gov.GetForeignAidBarData(request.Country)
	response := gin.H{
		"barData": barData,
	}
	c.JSON(http.StatusOK, response)
}
