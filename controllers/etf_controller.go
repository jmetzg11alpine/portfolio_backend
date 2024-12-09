package controllers

import (
	"backend/models/etf"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetETFReserveData(c *gin.Context) {
	reserves := etf.GetETFReserves()
	response := gin.H{
		"reserves": reserves,
	}
	c.JSON(http.StatusOK, response)
}

type TimePeriodRequest struct {
	TimePeriod string `json:"timePeriod"`
}

func GetStockData(c *gin.Context) {
	var req TimePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("time period: %s\n", req.TimePeriod)
	etf.GetStocks()
	response := gin.H{
		"message": "help",
	}
	c.JSON(http.StatusOK, response)
}
