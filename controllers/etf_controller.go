package controllers

import (
	"backend/models/etf"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetETFData(c *gin.Context) {
	reserves := etf.GetETFReserves()
	transactions := etf.GetETFTransactions()
	response := gin.H{
		"reserves":     reserves,
		"transactions": transactions,
	}
	c.JSON(http.StatusOK, response)
}

type TimePeriodRequest struct {
	TimePeriod string `json:"timePeriod"`
}
