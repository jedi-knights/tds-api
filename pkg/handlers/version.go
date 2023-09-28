package handlers

import (
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleGetVersion(c echo.Context) error {
	var (
		version string
		err     error
	)

	service := services.NewVersion()

	if version, err = service.GetVersion(); err != nil {
		return api.Error{
			Status: http.StatusInternalServerError,
			Msg:    "failed to retrieve version: " + err.Error(),
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"version": version})
}
