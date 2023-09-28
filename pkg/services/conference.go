package services

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/models"
	"go.uber.org/zap"
)

var divisionToUrlMapping = map[pkg.Division]string{
	pkg.DivisionDI:    "https://www.topdrawersoccer.com/college-soccer/college-conferences/di/divisionid-1",
	pkg.DivisionDII:   "https://www.topdrawersoccer.com/college-soccer/college-conferences/dii/divisionid-2",
	pkg.DivisionDIII:  "https://www.topdrawersoccer.com/college-soccer/college-conferences/diii/divisionid-3",
	pkg.DivisionNAIA:  "https://www.topdrawersoccer.com/college-soccer/college-conferences/naia/divisionid-4",
	pkg.DivisionNJCAA: "https://www.topdrawersoccer.com/college-soccer/college-conferences/njcaa/divisionid-5",
}

/*
For better performance the Conference service will cache the conferences in memory.
*/

// var Conferences []models.Conference

type ConferenceService struct {
	logger *zap.Logger
}

func NewConference() *ConferenceService {
	logger := api.GetLogger()

	logger.Debug("NewConference")

	return &ConferenceService{logger: logger}
}

// GetConferencesByDivision returns all conferences associated with the specified division
func (s ConferenceService) GetConferencesByDivision(division pkg.Division) ([]models.Conference, error) {
	var err error

	s.logger.Debug("GetConferencesByDivision", zap.String("division", pkg.DivisionToString(division)))

	conferences := []models.Conference{}

	url, ok := divisionToUrlMapping[division]
	if !ok {
		return nil, fmt.Errorf("the division %s is not supported", division)
	}

	c := colly.NewCollector()

	c.OnHTML("table > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			decorator := api.Decorate(el, pkg.Prefix)

			conferenceUrl := decorator.GetLinkFromCell(1)
			conferences = append(conferences, models.Conference{
				Id:       api.GetIdFromUrl(conferenceUrl),
				Name:     decorator.GetTextFromCell(1),
				Url:      conferenceUrl,
				Division: pkg.DivisionToString(division),
				Gender:   pkg.GenderToString(api.GetGenderFromUrl(conferenceUrl)),
			})
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err = c.Visit(url); err != nil {
		return nil, err
	}

	return conferences, nil
}

// GetConferences returns a slice of all conferences.
func (s ConferenceService) GetConferences() ([]models.Conference, error) {
	var (
		err                 error
		conferences         []models.Conference
		divisionConferences []models.Conference
	)

	s.logger.Debug("GetConferences")

	for _, division := range pkg.Divisions {
		if divisionConferences, err = s.GetConferencesByDivision(division); err != nil {
			return nil, err
		}

		conferences = append(conferences, divisionConferences...)
	}

	return conferences, nil
}

// GetConferenceNames returns a slice of strings containing the names of all conferences.
func (s ConferenceService) GetConferenceNames() ([]string, error) {
	s.logger.Debug("GetConferenceNames")

	return nil, nil
}

// GetConference returns a conference associated with the specified name.
func (s ConferenceService) GetConference(name string) (*models.Conference, error) {
	s.logger.Debug("GetConference", zap.String("name", name))

	return nil, fmt.Errorf("not implemented")
}

// GetConferencesByGender returns a slice of conferences associated with the specified gender.
func (s ConferenceService) GetConferencesByGender(gender pkg.Gender) ([]models.Conference, error) {
	s.logger.Debug("GetConferencesByGender", zap.String("gender", pkg.GenderToString(gender)))

	return nil, nil
}
