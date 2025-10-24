package controllers

import (
	"net/http"
	"time"

	"github.com/ADEMOLA200/hng-country-api/models"
	"github.com/ADEMOLA200/hng-country-api/services"
	"github.com/gin-gonic/gin"
)

type StatusController struct {
	countryService *services.CountryService
}

func NewStatusController() *StatusController {
	return &StatusController{
		countryService: services.NewCountryService(),
	}
}

func (sc *StatusController) GetStatus(c *gin.Context) {
	status, err := sc.countryService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (sc *StatusController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"timestamp": time.Now().Unix(),
		"service":   "HNG Country API",
		"version":   "1.0.0",
	})
}
