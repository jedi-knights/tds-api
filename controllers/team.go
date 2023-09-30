package controllers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jedi-knights/tds-api/models"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"net/url"
	"strconv"
)

type Teamer interface {
	GetAll() ([]models.Team, error)
	GetById(targetId int) (*models.Team, error)
	GetByGender(targetGender pkg.Gender) ([]models.Team, error)
	GetByNameAndGender(targetName string, targetGender pkg.Gender) (*models.Team, error)
	GetByDivision(targetDivision pkg.Division) ([]models.Team, error)
	GetByGenderAndDivision(targetGender pkg.Gender, targetDivision pkg.Division) ([]models.Team, error)
	GetByConference(targetConference models.Conference) ([]models.Team, error)
	GetByConferenceName(targetConferenceName string) ([]models.Team, error)
	GetByConferenceNameAndGender(targetConferenceName string, targetGender pkg.Gender) ([]models.Team, error)
	GetByConferenceIdAndGender(targetConferenceId int, targetGender pkg.Gender) ([]models.Team, error)
}

type Team struct{}

func NewTeam() *Team {
	return &Team{}
}

func (t Team) GetAll() ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	for _, division := range pkg.Divisions {
		if divisionTeams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}
		teams = append(teams, divisionTeams...)
	}

	return teams, nil
}

func (t Team) GetById(targetId int) (*models.Team, error) {
	var (
		err   error
		teams []models.Team
	)

	for _, division := range pkg.Divisions {
		if teams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}

		for _, team := range teams {
			if team.Id == targetId {
				return &team, nil
			}
		}
	}

	return nil, fmt.Errorf("team with targetId '%d' not found", targetId)
}

func (t Team) GetByGender(targetGender pkg.Gender) ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	targetGenderString := targetGender.String()

	for _, division := range pkg.Divisions {
		if divisionTeams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}

		for _, team := range divisionTeams {
			if team.Gender == targetGenderString {
				teams = append(teams, team)
			}
		}
	}

	return teams, nil
}

func (t Team) GetByNameAndGender(targetName string, targetGender pkg.Gender) (*models.Team, error) {
	var (
		err   error
		teams []models.Team
	)

	targetGenderString := targetGender.String()

	for _, division := range pkg.Divisions {
		if teams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}

		for _, team := range teams {
			if team.Name == targetName && team.Gender == targetGenderString {
				return &team, nil
			}
		}
	}

	return nil, fmt.Errorf("team with targetName '%s' and targetGender '%s' not found", targetName, targetGender.String())
}

func (t Team) GetByDivision(targetDivision pkg.Division) ([]models.Team, error) {
	var (
		err         error
		divisionUrl string
		teams       []models.Team
	)

	if divisionUrl, err = targetDivision.Url(); err != nil {
		return nil, fmt.Errorf("failed to retrieve url for division '%s': %v", targetDivision.String(), err.Error())
	}

	c := colly.NewCollector()

	c.OnHTML("table.table-striped", func(t *colly.HTMLElement) {
		var currentTeams []models.Team

		if currentTeams, err = readTeamsFromTable(pkg.GenderBoth, targetDivision, t); err != nil {
			fmt.Println(err)
			return
		}

		teams = append(teams, currentTeams...)
	})

	if err = c.Visit(divisionUrl); err != nil {
		return nil, fmt.Errorf("failed to visit url '%s': %v", divisionUrl, err.Error())
	}

	return teams, nil
}

func (t Team) GetByGenderAndDivision(targetGender pkg.Gender, targetDivision pkg.Division) ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	if targetDivision == pkg.DivisionAll || targetDivision == pkg.DivisionUnknown {
		return t.GetByGender(targetGender)
	}

	if targetGender == pkg.GenderBoth || targetGender == pkg.GenderUnknown {
		return t.GetByDivision(targetDivision)
	}

	if divisionTeams, err = t.GetByDivision(targetDivision); err != nil {
		return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", targetDivision.String(), err.Error())
	}

	for _, team := range divisionTeams {
		if team.Gender == targetGender.String() {
			teams = append(teams, team)
		}
	}

	return teams, nil
}

func (t Team) GetByConference(targetConference models.Conference) ([]models.Team, error) {
	return t.GetByConferenceName(targetConference.Name)
}

func (t Team) GetByConferenceName(targetConferenceName string) ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	for _, division := range pkg.Divisions {
		if divisionTeams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}

		for _, team := range divisionTeams {
			if team.ConferenceName == targetConferenceName {
				teams = append(teams, team)
			}
		}
	}

	return teams, nil
}

func (t Team) GetByConferenceNameAndGender(targetConferenceName string, targetGender pkg.Gender) ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	targetGenderString := targetGender.String()

	for _, division := range pkg.Divisions {
		if divisionTeams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}

		for _, team := range divisionTeams {
			if team.ConferenceName == targetConferenceName && team.Gender == targetGenderString {
				teams = append(teams, team)
			}
		}
	}

	return teams, nil
}

func (t Team) GetByConferenceIdAndGender(targetConferenceId int, targetGender pkg.Gender) ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	targetGenderString := targetGender.String()

	for _, division := range pkg.Divisions {
		if divisionTeams, err = t.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve teams by division '%s': %v", division.String(), err.Error())
		}

		for _, team := range divisionTeams {
			if team.Id == targetConferenceId && team.Gender == targetGenderString {
				teams = append(teams, team)
			}
		}
	}

	return teams, nil
}

func getConferenceIdFromUrl(targetUrl string) (int, error) {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		return 0, err
	}

	params, err := url.ParseQuery(myUrl.RawQuery)
	if err != nil {
		return 0, err
	}

	conferenceIdString := params.Get("conferenceId")

	conferenceId, err := strconv.Atoi(conferenceIdString)
	if err != nil {
		return 0, err
	}

	return conferenceId, nil
}

func getConferenceFromTableCaption(t *colly.HTMLElement) (models.Conference, error) {
	var conferenceId int
	var conferenceUrl string
	var err error
	var conference models.Conference

	href := t.ChildAttr("caption > a", "href")

	conferenceUrl = href

	if conferenceUrl, err = url.JoinPath(pkg.Prefix, conferenceUrl); err != nil {
		return conference, err
	}

	if conferenceUrl, err = url.QueryUnescape(conferenceUrl); err != nil {
		return conference, err
	}

	if conferenceUrl, err = url.QueryUnescape(conferenceUrl); err != nil {
		return conference, err
	}

	if conferenceId, err = getConferenceIdFromUrl(conferenceUrl); err != nil {
		return conference, err
	}

	name := t.ChildText("caption")
	name = api.NormalizeText(name)

	return models.Conference{
		Id:     conferenceId,
		Name:   name,
		Url:    conferenceUrl,
		Gender: "",
	}, nil
}

func getTeamNameFromRow(el *colly.HTMLElement) (string, error) {
	text := el.ChildText("td:nth-child(1)")
	text = api.NormalizeText(text)

	if text == "" {
		return "", fmt.Errorf("could not find team name")
	}

	return text, nil
}

func getTeamLinkFromRow(el *colly.HTMLElement) (string, error) {
	var link string
	var err error

	if link = el.ChildAttr("td:nth-child(1) > a", "href"); link == "" {
		return "", fmt.Errorf("could not find team link")
	}

	if link, err = url.JoinPath(pkg.Prefix, link); err != nil {
		return "", err
	}

	return link, nil
}

func getTeamDetailsFromRow(el *colly.HTMLElement) (string, string, error) {
	var text, link string
	var err error

	if text, err = getTeamNameFromRow(el); err != nil {
		return "", "", err
	}

	if link, err = getTeamLinkFromRow(el); err != nil {
		return "", "", err
	}

	return text, link, nil
}

func readTeamsFromTable(targetGender pkg.Gender, division pkg.Division, t *colly.HTMLElement) ([]models.Team, error) {
	var (
		err   error
		teams []models.Team
	)

	conference, err := getConferenceFromTableCaption(t)
	if err != nil {
		return nil, err
	}

	// Scan each of the rows in the table
	t.ForEach("tr", func(_ int, el *colly.HTMLElement) {
		var text, link string

		if text, link, err = getTeamDetailsFromRow(el); err != nil {
			return
		}

		currentGender := api.GetGenderFromUrl(link)

		if targetGender != pkg.GenderUnknown && targetGender != pkg.GenderBoth {
			if currentGender != targetGender {
				return
			}
		}

		teams = append(teams, models.Team{
			Id:             api.GetIdFromUrl(link),
			Name:           text,
			Url:            link,
			Gender:         pkg.GenderToString(currentGender),
			Division:       pkg.DivisionToString(division),
			ConferenceId:   conference.Id,
			ConferenceName: conference.Name,
			ConferenceUrl:  conference.Url,
		})
	})

	return teams, nil
}
