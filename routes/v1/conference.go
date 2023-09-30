package v1

import (
	"github.com/jedi-knights/tds-api/controllers"
	"github.com/jedi-knights/tds-api/models"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/routes"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HandleGetConferences godoc
// @Summary Get a list of conferences
// @Description Get a list of conferences
// @Tags Conferences
// @Accept json
// @Produce json
// Success 200 {array} models.Conference
// @Router /v2/conferences [get]
func HandleGetConferences(c echo.Context) error {
	var (
		err         error
		conferences []models.Conference
	)
	c.Logger().Debug("HandleGetConferences")

	if conferences, err = controllers.NewConference(c).GetAll(); err != nil {
		return routes.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	return c.JSON(http.StatusOK, conferences)
}

// HandleGetConferencesByDivision godoc
// @Summary Get a list of conferences
// @Description Get a list of conferences
// @Tags Conferences
// @Accept json
// @Produce json
// @Param division query string false "Specify a division you are interested in" Enums(all,di,dii,diii,naia,njcaa) Default(all)
// @Success 200 {array} models.Conference
// @Router /v2/conferences/:division [get]
func HandleGetConferencesByDivision(c echo.Context) error {
	var (
		err         error
		division    pkg.Division
		conferences []models.Conference
	)
	c.Logger().Debug("HandleGetConferencesByDivision")

	// read query parameter
	division = pkg.StringToDivision(c.QueryParam("division"))

	if conferences, err = controllers.NewConference(c).GetByDivision(division); err != nil {
		return routes.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	return c.JSON(http.StatusOK, conferences)
}
