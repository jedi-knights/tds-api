package services

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/models"
	"go.uber.org/zap"
	"slices"
	"sort"
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

	divisionString := pkg.DivisionToString(division)

	s.logger.Debug("GetConferencesByDivision", zap.String("division", divisionString))

	conferences := []models.Conference{}

	if division == pkg.DivisionAll {
		return s.GetConferences()
	}

	url, ok := divisionToUrlMapping[division]
	if !ok {
		s.logger.Error("unsupported division", zap.String("division", divisionString))
		return nil, fmt.Errorf("the division %s is not supported", divisionString)
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
		s.logger.Error("failed to retrieve conferences by division",
			zap.String("url", r.Request.URL.String()),
			zap.String("division", divisionString),
			zap.Error(err),
		)
	})

	if err = c.Visit(url); err != nil {
		s.logger.Error("failed to retrieve conferences by division",
			zap.String("url", url),
			zap.String("division", divisionString),
			zap.Error(err),
		)

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
			s.logger.Error("failed to retrieve conferences by division",
				zap.String("division", pkg.DivisionToString(division)),
				zap.Error(err),
			)

			return nil, err
		}

		conferences = append(conferences, divisionConferences...)
	}

	return conferences, nil
}

// GetConferenceNamesByDivision returns a slice of strings containing the names of all conferences associated with
// the specified division.
func (s ConferenceService) GetConferenceNamesByDivision(division pkg.Division) ([]string, error) {
	s.logger.Debug("GetConferenceNamesByDivision", zap.String("division", pkg.DivisionToString(division)))

	var (
		ok          bool
		err         error
		divisionUrl string
	)

	divisionString := pkg.DivisionToString(division)

	s.logger.Debug("GetConferencesByDivision", zap.String("division", divisionString))

	conferenceNames := []string{}

	if divisionUrl, ok = divisionToUrlMapping[division]; !ok {
		s.logger.Error("unsupported division", zap.String("division", divisionString))
	}

	c := colly.NewCollector()

	c.OnHTML("table > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			decorator := api.Decorate(el, pkg.Prefix)

			conferenceName := decorator.GetTextFromCell(1)

			if slices.Contains(conferenceNames, conferenceName) {
				return
			}

			conferenceNames = append(conferenceNames, conferenceName)
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		s.logger.Error("failed to retrieve conference names", zap.Error(err))
	})

	if err = c.Visit(divisionUrl); err != nil {
		return nil, err
	}

	return conferenceNames, nil
}

// GetConferenceNames returns a slice of strings containing the names of all conferences.
func (s ConferenceService) GetConferenceNames() ([]string, error) {
	s.logger.Debug("GetConferenceNames")

	var (
		err             error
		divisionNames   []string
		conferenceNames []string
	)

	for _, division := range pkg.Divisions {
		if divisionNames, err = s.GetConferenceNamesByDivision(division); err != nil {
			return nil, err
		}

		for _, name := range divisionNames {
			if slices.Contains(conferenceNames, name) {
				continue
			}

			conferenceNames = append(conferenceNames, name)
		}
	}

	// sort the slice of conference names
	sort.SliceStable(conferenceNames, func(i, j int) bool {
		return conferenceNames[i] < conferenceNames[j]
	})

	return conferenceNames, nil
}

// GetConferenceByDivision returns a conference associated with the specified name and division.
func (s ConferenceService) GetConferenceByDivision(name string, division pkg.Division) (*models.Conference, error) {
	var (
		ok          bool
		err         error
		divisionUrl string
		conference  *models.Conference
	)

	s.logger.Debug("GetConferenceByDivision",
		zap.String("name", name),
		zap.String("division", pkg.DivisionToString(division)),
	)

	divisionString := pkg.DivisionToString(division)

	s.logger.Debug("GetConferenceByDivision", zap.String("division", divisionString))

	if divisionUrl, ok = divisionToUrlMapping[division]; !ok {
		s.logger.Error("unsupported division", zap.String("division", divisionString))
	}

	c := colly.NewCollector()

	c.OnHTML("table > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			decorator := api.Decorate(el, pkg.Prefix)

			conferenceName := decorator.GetTextFromCell(1)

			if conferenceName != name {
				return
			}

			conferenceUrl := decorator.GetLinkFromCell(1)

			gender := api.GetGenderFromUrl(conferenceUrl)

			conference = &models.Conference{
				Id:       api.GetIdFromUrl(conferenceUrl),
				Name:     decorator.GetTextFromCell(1),
				Url:      conferenceUrl,
				Division: pkg.DivisionToString(division),
				Gender:   pkg.GenderToString(gender),
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		s.logger.Error("failed to retrieve conference names", zap.Error(err))
	})

	if err = c.Visit(divisionUrl); err != nil {
		return nil, err
	}

	if conference == nil {
		s.logger.Error("conference not found", zap.String("name", name))

		return nil, fmt.Errorf("conference not found '%s'", name)
	} else {
		s.logger.Debug("conference found", zap.String("name", name))

		return conference, nil
	}
}

// GetConference returns a conference associated with the specified name.
func (s ConferenceService) GetConference(name string) (*models.Conference, error) {
	s.logger.Debug("GetConference", zap.String("name", name))

	return nil, fmt.Errorf("not implemented")
}

// GetConferencesByGender returns a slice of conferences associated with the specified gender.
func (s ConferenceService) GetConferencesByGender(gender pkg.Gender) ([]models.Conference, error) {
	var (
		err                 error
		conferences         []models.Conference
		divisionConferences []models.Conference
	)

	s.logger.Debug("GetConferencesByGender", zap.String("gender", pkg.GenderToString(gender)))

	for _, division := range pkg.Divisions {
		if divisionConferences, err = s.GetConferencesByGenderAndDivision(gender, division); err != nil {
			return nil, err
		}

		conferences = append(conferences, divisionConferences...)
	}

	return conferences, nil
}

func (s ConferenceService) GetConferencesByGenderAndDivision(targetGender pkg.Gender, targetDivision pkg.Division) ([]models.Conference, error) {
	var (
		err error
	)

	if targetDivision == pkg.DivisionAll || targetDivision == pkg.DivisionUnknown {
		if targetGender == pkg.GenderBoth || targetGender == pkg.GenderUnknown {
			return s.GetConferences()
		} else {
			return s.GetConferencesByGender(targetGender)
		}
	} else {
		if targetGender == pkg.GenderBoth || targetGender == pkg.GenderUnknown {
			return s.GetConferencesByDivision(targetDivision)
		}
	}

	targetGenderString := pkg.GenderToString(targetGender)
	targetDivisionString := pkg.DivisionToString(targetDivision)

	s.logger.Debug("GetConferencesByGenderAndDivision",
		zap.String("targetGender", targetGenderString),
		zap.String("targetDivision", targetDivisionString),
	)

	conferences := []models.Conference{}

	url, ok := divisionToUrlMapping[targetDivision]
	if !ok {
		return nil, fmt.Errorf("the targetDivision %s is not supported", targetDivisionString)
	}

	c := colly.NewCollector()

	c.OnHTML("table > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			var (
				currentGender pkg.Gender
				decorator     *api.HTMLElementDecorator
			)

			decorator = api.Decorate(el, pkg.Prefix)

			conferenceUrl := decorator.GetLinkFromCell(1)

			if currentGender = api.GetGenderFromUrl(conferenceUrl); currentGender != targetGender {
				return
			}

			conferences = append(conferences, models.Conference{
				Id:       api.GetIdFromUrl(conferenceUrl),
				Name:     decorator.GetTextFromCell(1),
				Url:      conferenceUrl,
				Division: pkg.DivisionToString(targetDivision),
				Gender:   pkg.GenderToString(currentGender),
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
