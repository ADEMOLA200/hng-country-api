package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/ADEMOLA200/hng-country-api/models"
	"github.com/ADEMOLA200/hng-country-api/services"
	"github.com/gin-gonic/gin"
)

type CountryController struct {
	countryService *services.CountryService
	imageService   *services.ImageService
}

func NewCountryController() *CountryController {
	return &CountryController{
		countryService: services.NewCountryService(),
		imageService:   services.NewImageService(),
	}
}

func (cc *CountryController) RefreshCountries(c *gin.Context) {
	response, err := cc.countryService.RefreshCountries()
	if err != nil {
		if strings.Contains(err.Error(), "external data source unavailable") {
			c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
				Error: "External data source unavailable",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (cc *CountryController) GetAllCountries(c *gin.Context) {
	filters := make(map[string]string)

	if region := c.Query("region"); region != "" {
		filters["region"] = region
	}
	if currency := c.Query("currency"); currency != "" {
		filters["currency"] = currency
	}

	sort := c.Query("sort")

	countries, err := cc.countryService.GetAllCountries(filters, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	var response []models.CountryResponse
	for _, country := range countries {
		response = append(response, models.CountryResponse{
			ID:              country.ID,
			Name:            country.Name,
			Capital:         country.Capital,
			Region:          country.Region,
			Population:      country.Population,
			CurrencyCode:    country.CurrencyCode,
			ExchangeRate:    country.ExchangeRate,
			EstimatedGDP:    country.EstimatedGDP,
			FlagURL:         country.FlagURL,
			LastRefreshedAt: country.LastRefreshedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (cc *CountryController) GetCountryByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Country name is required",
		})
		return
	}

	country, err := cc.countryService.GetCountryByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Country not found",
		})
		return
	}

	response := models.CountryResponse{
		ID:              country.ID,
		Name:            country.Name,
		Capital:         country.Capital,
		Region:          country.Region,
		Population:      country.Population,
		CurrencyCode:    country.CurrencyCode,
		ExchangeRate:    country.ExchangeRate,
		EstimatedGDP:    country.EstimatedGDP,
		FlagURL:         country.FlagURL,
		LastRefreshedAt: country.LastRefreshedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (cc *CountryController) DeleteCountry(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Country name is required",
		})
		return
	}

	err := cc.countryService.DeleteCountryByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Country not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Country deleted successfully",
	})
}

func (cc *CountryController) GetCountriesImage(c *gin.Context) {
	imagePath := "cache/summary.png"

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Summary image not found",
		})
		return
	}

	c.File(imagePath)
}
