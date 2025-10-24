package routes

import (
	"github.com/ADEMOLA200/hng-country-api/controllers"
	"github.com/ADEMOLA200/hng-country-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ErrorHandler())
	router.Use(middleware.Logging())

	countryController := controllers.NewCountryController()
	statusController := controllers.NewStatusController()

	api := router.Group("/api")
	{
		countries := api.Group("/countries")
		{
			countries.POST("/refresh", countryController.RefreshCountries)
			countries.GET("", countryController.GetAllCountries)
			countries.GET("/image", countryController.GetCountriesImage)
			countries.GET("/:name", countryController.GetCountryByName)
			countries.DELETE("/:name", countryController.DeleteCountry)
		}

		api.GET("/status", statusController.GetStatus)
	}

	router.GET("/health", statusController.HealthCheck)

	return router
}
