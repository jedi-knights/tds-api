package handlers

import (
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/models"
	"github.com/jedi-knights/tds-api/pkg/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

var logger = api.GetLogger()

func HandleGetConferenceNames(c echo.Context) error {
	return nil
}

func HandleGetConferences(c echo.Context) error {
	var (
		err         error
		conferences []models.Conference
	)

	logger.Debug("HandleGetConferences")

	if conferences, err = services.NewConference().GetConferences(); err != nil {
		return api.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	return c.JSON(http.StatusOK, conferences)
}
