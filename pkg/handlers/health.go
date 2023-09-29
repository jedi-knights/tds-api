package handlers

import (
	"github.com/jedi-knights/tds-api/pkg/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HandleHealthCheck godoc
// @Summary Health Check
// @Description Check if the API is up and running
// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200 {object} models.HealthCheckResponse
// @Router /health [get]
func HandleHealthCheck(c echo.Context) error {
	var response = services.HeathService{}.GetHealthCheck()

	return c.JSON(http.StatusOK, response)
}
