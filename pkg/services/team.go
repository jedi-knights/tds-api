package services

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/models"
	"go.uber.org/zap"
	"net/url"
	"strconv"
)

var teamDivisionToUrlMapping = map[pkg.Division]string{
	pkg.DivisionDI:    "https://www.topdrawersoccer.com/college/teams/?divisionName=di&divisionId=1",
	pkg.DivisionDII:   "https://www.topdrawersoccer.com/college/teams/?divisionName=dii&divisionId=2",
	pkg.DivisionDIII:  "https://www.topdrawersoccer.com/college/teams/?divisionName=diii&divisionId=3",
	pkg.DivisionNAIA:  "https://www.topdrawersoccer.com/college/teams/?divisionName=naia&divisionId=4",
	pkg.DivisionNJCAA: "https://www.topdrawersoccer.com/college/teams/?divisionName=njcaa&divisionId=5",
}

type TeamServicer interface {
	GetConferencesByDivision(division pkg.Division) ([]models.Conference, error)

	GetConferences() ([]models.Conference, error)
	GetConferenceById(id int) (*models.Conference, error)
	GetTeams() ([]models.Team, error)
	GetTeamById(id int) (*models.Team, error)
	GetTeamsByGender(gender pkg.Gender) ([]models.Team, error)
	GetTeamsByGenderAndDivision(gender pkg.Gender, division pkg.Division) ([]models.Team, error)
	GetTeamByNameAndGender(name, string, gender pkg.Gender) (*models.Team, error)
	GetTeamsByConferenceName(name string) ([]models.Team, error)
	GetTeamsByConferenceNameAndGender(name string, gender pkg.Gender) ([]models.Team, error)
	GetTeamsByConferenceId(id int) ([]models.Team, error)
	GetTeamsByConferenceIdAndGender(conferenceId int, gender pkg.Gender) ([]models.Team, error)
	GetUrlByDivision(division pkg.Division) (string, error)
	GetTeamsByDivision(division pkg.Division) ([]models.Team, error)
	GetTeamsByDivisionAndGender(division pkg.Division, gender pkg.Gender) ([]models.Team, error)
}

type TeamService struct {
	logger *zap.Logger
}

func NewTeam() *TeamService {
	return &TeamService{}
}

func (s TeamService) GetTeams() ([]models.Team, error) {
	var (
		err   error
		teams []models.Team
	)

	teams = make([]models.Team, 0)

	for _, division := range pkg.Divisions {
		var currentTeams []models.Team

		if currentTeams, err = s.GetTeamsByDivision(division); err != nil {
			return nil, err
		}

		teams = append(teams, currentTeams...)
	}

	return teams, nil
}

func (s TeamService) GetTeamsByDivision(division pkg.Division) ([]models.Team, error) {
	var (
		ok          bool
		err         error
		divisionUrl string
		teams       []models.Team
	)

	teams = make([]models.Team, 0)

	if divisionUrl, ok = teamDivisionToUrlMapping[division]; !ok {
		s.logger.Error("GetTeamsByDivision",
			zap.String("division", pkg.DivisionToString(division)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("the division %s is not supported", pkg.DivisionToString(division))
	}

	c := colly.NewCollector()

	c.OnHTML("table.table-striped", func(t *colly.HTMLElement) {
		var currentTeams []models.Team

		if currentTeams, err = readTeamsFromTable(pkg.GenderBoth, division, t); err != nil {
			s.logger.Debug("Skipping a table of teams", zap.Error(err))
			return
		}

		teams = append(teams, currentTeams...)
	})

	c.OnError(func(r *colly.Response, err error) {
		s.logger.Debug("GetTeamsByDivision",
			zap.String("Request URL", r.Request.URL.String()),
			zap.String("division", pkg.DivisionToString(division)),
			zap.Error(err),
		)
	})

	if err = c.Visit(divisionUrl); err != nil {
		return nil, err
	}

	return teams, nil
}

func (s TeamService) GetTeamsByConference(conference string) ([]models.Team, error) {
	var (
		err           error
		teams         []models.Team
		divisionTeams []models.Team
	)

	teams = make([]models.Team, 0)

	for _, division := range pkg.Divisions {
		if divisionTeams, err = s.GetTeamsByDivisionAndConference(division, conference); err != nil {
			return nil, err
		}

		teams = append(teams, divisionTeams...)
	}

	return teams, nil
}

func (s TeamService) GetTeamsByDivisionAndConference(division pkg.Division, targetConferenceName string) ([]models.Team, error) {
	var (
		ok          bool
		err         error
		divisionUrl string
		teams       []models.Team
	)

	teams = make([]models.Team, 0)

	if divisionUrl, ok = teamDivisionToUrlMapping[division]; !ok {
		s.logger.Error("GetTeamsByDivisionAndConference",
			zap.String("division", pkg.DivisionToString(division)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("the division %s is not supported", pkg.DivisionToString(division))
	}

	c := colly.NewCollector()

	c.OnHTML("table.table-striped", func(t *colly.HTMLElement) {
		// Get the conference name from the table caption
		conference, err := getConferenceFromTableCaption(t)
		if err != nil {
			s.logger.Debug("Skipping a table of teams", zap.Error(err))
			return
		}

		if conference.Name != targetConferenceName {
			return
		}

		var currentTeams []models.Team

		if currentTeams, err = readTeamsFromTable(pkg.GenderBoth, division, t); err != nil {
			s.logger.Debug("Skipping a table of teams", zap.Error(err))
			return
		}

		teams = append(teams, currentTeams...)
	})

	c.OnError(func(r *colly.Response, err error) {
		s.logger.Debug("GetTeamsByDivisionAndConference",
			zap.String("Request URL", r.Request.URL.String()),
			zap.String("division", pkg.DivisionToString(division)),
			zap.Error(err),
		)
	})

	if err = c.Visit(divisionUrl); err != nil {
		return nil, err
	}

	return teams, nil
}

// GetConferencesByDivision returns all conferences associated with the specified division
func (s TeamService) GetConferencesByDivision(division pkg.Division) ([]models.Conference, error) {
	var (
		ok          bool
		err         error
		divisionUrl string
		conferences []models.Conference
	)

	conferences = make([]models.Conference, 0)

	if divisionUrl, ok = teamDivisionToUrlMapping[division]; !ok {
		s.logger.Error("GetConferencesByDivision",
			zap.String("division", pkg.DivisionToString(division)),
			zap.Error(err),
		)

		return nil, fmt.Errorf("the division %s is not supported", pkg.DivisionToString(division))
	}

	c := colly.NewCollector()

	c.OnHTML("table.table-striped", func(t *colly.HTMLElement) {
		var currentConference models.Conference

		if currentConference, err = getConferenceFromTableCaption(t); err != nil {
			s.logger.Debug("Skipping a table of teams", zap.Error(err))
			return
		}

		currentConference.Division = pkg.DivisionToString(division)
		currentConference.Gender = pkg.GenderToString(pkg.GenderBoth)

		conferences = append(conferences, currentConference)
	})

	c.OnError(func(r *colly.Response, err error) {
		s.logger.Debug("GetConferencesByDivision",
			zap.String("Request URL", r.Request.URL.String()),
			zap.String("division", pkg.DivisionToString(division)),
			zap.Error(err),
		)
	})

	if err = c.Visit(divisionUrl); err != nil {
		return nil, err
	}

	return conferences, nil
}

// GetConferences returns all conferences
func (s TeamService) GetConferences() ([]models.Conference, error) {
	conferences := make([]models.Conference, 0)

	for _, division := range pkg.Divisions {
		var (
			err                 error
			divisionConferences []models.Conference
		)

		if divisionConferences, err = s.GetConferencesByDivision(division); err != nil {
			s.logger.Error("GetConferences",
				zap.String("division", pkg.DivisionToString(division)),
				zap.Error(err),
			)

			return nil, err
		}

		conferences = append(conferences, divisionConferences...)
	}

	return conferences, nil
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
