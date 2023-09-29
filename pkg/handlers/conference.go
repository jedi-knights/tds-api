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

// HandleGetConferences godoc
// @Summary Get a list of conferences
// @Description Get a list of conferences
// @Tags Conferences
// @Accept  json
// @Produce  json
// @Param division query string false "Specify a division you are interested in" Enums(all,di,dii,diii,naia,njcaa) Default(all)
// @Success 200 {array} models.Conference
// @Router /conferences/:division [get]
func HandleGetConferences(c echo.Context) error {
	var (
		err         error
		conferences []models.Conference
	)

	api.GetLogger().Debug("HandleGetConferences")

	divisionString := c.QueryParam("division")
	division := pkg.StringToDivision(divisionString)

	if division == pkg.DivisionUnknown {
		return api.BadRequestError("expected division to be (di|dii|diii|naia|njcaa)")
	}

	api.GetLogger().Debug("division", zap.String("division", divisionString))

	svc := services.NewConference()

	res := c.Response()
	header := res.Header()

	if conferences, err = svc.GetConferencesByDivision(division); err != nil {
		header.Set("X-Total-Count", "0")

		return api.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	header.Set("X-Total-Count", strconv.Itoa(len(conferences)))

	return c.JSON(http.StatusOK, conferences)
}
