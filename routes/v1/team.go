package v1

import (
	"github.com/jedi-knights/tds-api/controllers"
	"github.com/jedi-knights/tds-api/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HandleGetTeams godoc
// @Summary Get all teams
// @Description Get all teams
// @Tags Teams
// @Accept json
// @Produce json
// @Success 200 {array} models.Team
// @Router /v2/teams [get]
func HandleGetTeams(c echo.Context) error {
	var (
		err   error
		teams []models.Team
	)

	if teams, err = controllers.NewTeam().GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, teams)
}
