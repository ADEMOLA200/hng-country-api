package middleware

import (
	"log"
	"net/http"

	"github.com/ADEMOLA200/hng-country-api/models"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			log.Printf("Error: %v | Path: %s | Method: %s", err, c.Request.URL.Path, c.Request.Method)

			var statusCode int
			var errorResponse models.ErrorResponse

			switch err.Err.(type) {
			case *ValidationError:
				statusCode = http.StatusBadRequest
				validationErr := err.Err.(*ValidationError)
				errorResponse = models.ErrorResponse{
					Error:   "Validation failed",
					Details: validationErr.Details,
				}
			default:
				errorMsg := err.Error()
				switch {
				case errorMsg == "country not found":
					statusCode = http.StatusNotFound
					errorResponse = models.ErrorResponse{
						Error: "Country not found",
					}
				case errorMsg == "external data source unavailable":
					statusCode = http.StatusServiceUnavailable
					errorResponse = models.ErrorResponse{
						Error: "External data source unavailable",
					}
				default:
					statusCode = http.StatusInternalServerError
					errorResponse = models.ErrorResponse{
						Error: "Internal server error",
					}
				}
			}

			c.JSON(statusCode, errorResponse)
			c.Abort()
		}
	}
}

type ValidationError struct {
	Details map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(details map[string]string) *ValidationError {
	return &ValidationError{Details: details}
}
