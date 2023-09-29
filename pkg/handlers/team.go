package handlers

import (
	"fmt"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/models"
	"github.com/jedi-knights/tds-api/pkg/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

// HandleGetTeams returns teams
//func HandleGetTeams(c echo.Context) error {
//	var (
//		err   error
//		teams []models.Team
//	)
//
//	api.GetLogger().Debug("HandleGetTeams")
//
//	genderString := c.QueryParam("gender")
//	gender := pkg.StringToGender(genderString)
//
//	if gender == pkg.GenderUnknown {
//		return api.BadRequestError("expected gender to be (male|female|both)")
//	}
//
//	divisionString := c.QueryParam("division")
//	division := pkg.StringToDivision(divisionString)
//
//	if division == pkg.DivisionUnknown {
//		return api.BadRequestError("expected division to be (di|dii|diii|naia|njcaa)")
//	}
//
//	api.GetLogger().Debug("gender", zap.String("gender", genderString))
//	api.GetLogger().Debug("division", zap.String("division", divisionString))
//
//	svc := services.NewTeam()
//
//
//}

func HandleGetTeamsConferencesByDivision(c echo.Context) error {
	var (
		err         error
		conferences []models.Conference
	)

	// Conferences from path  `teams/conferences/:division`
	divisionString := c.Param("division")

	api.GetLogger().Debug("HandleGetTeamsConferencesByDivision division path parameter", zap.String("division", divisionString))

	division := pkg.StringToDivision(divisionString)

	if division == pkg.DivisionUnknown {
		return api.BadRequestError(fmt.Sprintf("expected division '%s' to be (di|dii|diii|naia|njcaa)", divisionString))
	}

	if conferences, err = services.NewTeam().GetConferencesByDivision(division); err != nil {
		return api.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(conferences)))

	return c.JSON(http.StatusOK, conferences)
}

func HandleGetTeamsConferences(c echo.Context) error {
	var (
		err         error
		conferences []models.Conference
	)

	// Conferences from path  `teams/conferences`
	api.GetLogger().Debug("HandleGetTeamsConferences")

	if conferences, err = services.NewTeam().GetConferences(); err != nil {
		return api.InternalServerError("failed to retrieve conferences: " + err.Error())
	}

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(conferences)))

	return c.JSON(http.StatusOK, conferences)
}

// HandleGetTeams godoc
// @Summary Get a list of teams
// @Description Get a list of teams
// @Tags Teams
// @Accept  json
// @Produce  json
// @Param gender query string false "Specify a gender" Enums(both,male,female) Default(both)
// @Param division query string false "Specify a division you are interested in" Enums(all,di,dii,diii,naia,njcaa) Default(all)
// @Param conference query string false "Specify a conference you are interested in"
// @Param name query string false "The name of the entity you are looking for"
// @Param nameLike query string false "A partial name of the entity you are looking for"
// @Param id query int false "Specify a target id"
// @Success 200 {array} models.Team
// @Failure 400 {object} api.Error
// @Router /teams [get]
func HandleGetTeams(c echo.Context) error {
	var (
		err   error
		teams []models.Team
	)

	api.GetLogger().Debug("HandleGetTeams")

	// Get query parameters
	gender := c.QueryParam("gender")
	division := c.QueryParam("division")
	conference := c.QueryParam("conference")
	name := c.QueryParam("name")
	nameLike := c.QueryParam("nameLike")
	id := c.QueryParam("id")

	if division != "" {
		if conference != "" {
			if teams, err = services.NewTeam().GetTeamsByDivisionAndConference(pkg.StringToDivision(division), conference); err != nil {
				return api.InternalServerError("failed to retrieve teams: " + err.Error())
			}
		} else {
			if teams, err = services.NewTeam().GetTeamsByDivision(pkg.StringToDivision(division)); err != nil {
				return api.InternalServerError("failed to retrieve teams: " + err.Error())
			}
		}
	} else {
		if conference != "" {
			if teams, err = services.NewTeam().GetTeamsByConference(conference); err != nil {
				return api.InternalServerError("failed to retrieve teams: " + err.Error())
			}
		} else {
			if teams, err = services.NewTeam().GetTeams(); err != nil {
				return api.InternalServerError("failed to retrieve teams: " + err.Error())
			}
		}
	}

	// Now the teams list is either empty or is populated with teams from the division
	// and conference specified in the query parameters.

	// Apply filters based on query parameters
	filteredTeams := filterTeams(teams, gender, name, nameLike, id)

	// Set the X-Total-Count header
	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(filteredTeams)))

	return c.JSON(http.StatusOK, filteredTeams)
}

// filterTeams filters a slice of teams based on query parameters.
func filterTeams(teams []models.Team, gender, name, nameLike, id string) []models.Team {
	filteredTeams := teams

	// Apply filters based on query parameters
	if gender != "" && gender != "both" {
		filteredTeams = filterByGender(filteredTeams, gender)
	}
	if name != "" {
		filteredTeams = filterByName(filteredTeams, name)
	}
	if nameLike != "" {
		filteredTeams = filterByNameLike(filteredTeams, nameLike)
	}
	if id != "" {
		filteredTeams = filterByID(filteredTeams, id)
	}

	return filteredTeams
}

// Implement filter functions for each query parameter if needed.
// Example filter functions:
func filterByGender(teams []models.Team, gender string) []models.Team {
	filtered := []models.Team{}
	for _, team := range teams {
		if team.Gender == gender {
			filtered = append(filtered, team)
		}
	}
	return filtered
}

func filterByName(teams []models.Team, name string) []models.Team {
	filtered := []models.Team{}
	for _, team := range teams {
		if team.Name == name {
			filtered = append(filtered, team)
		}
	}
	return filtered
}

func filterByNameLike(teams []models.Team, nameLike string) []models.Team {
	filtered := []models.Team{}
	for _, team := range teams {
		if strings.Contains(team.Name, nameLike) {
			filtered = append(filtered, team)
		}
	}
	return filtered
}

func filterByID(teams []models.Team, id string) []models.Team {
	filtered := []models.Team{}
	for _, team := range teams {
		// Convert team ID to string for comparison
		teamID := strconv.Itoa(team.Id)
		if teamID == id {
			filtered = append(filtered, team)
		}
	}
	return filtered
}
