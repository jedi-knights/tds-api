package handlers

import (
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/models"
	"github.com/jedi-knights/tds-api/pkg/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var logger = api.GetLogger()

func HandleGetConferenceNames(c echo.Context) error {
	logger.Debug("HandleGetConferenceNames")

	var (
		err             error
		conferenceNames []string
	)

	svc := services.NewConference()

	if conferenceNames, err = svc.GetConferenceNames(); err != nil {
		return api.InternalServerError("failed to retrieve conference names: " + err.Error())
	}

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(conferenceNames)))

	return c.JSON(http.StatusOK, conferenceNames)
}

func HandleGetConferences(c echo.Context) error {
	var (
		err         error
		conferences []models.Conference
	)

	logger.Debug("HandleGetConferences")

	genderString := c.QueryParam("gender")
	gender := pkg.StringToGender(genderString)

	if gender == pkg.GenderUnknown {
		return api.BadRequestError("expected gender to be (male|female|both)")
	}

	divisionString := c.QueryParam("division")
	division := pkg.StringToDivision(divisionString)

	if division == pkg.DivisionUnknown {
		return api.BadRequestError("expected division to be (di|dii|diii|naia|njcaa)")
	}

	logger.Debug("gender", zap.String("gender", genderString))
	logger.Debug("division", zap.String("division", divisionString))

	svc := services.NewConference()

	res := c.Response()
	header := res.Header()

	if conferences, err = svc.GetConferencesByGenderAndDivision(gender, division); err != nil {
		header.Set("X-Total-Count", "0")

		return api.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	header.Set("X-Total-Count", strconv.Itoa(len(conferences)))

	return c.JSON(http.StatusOK, conferences)
}
