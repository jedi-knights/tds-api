package controllers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jedi-knights/tds-api/models"
	"github.com/jedi-knights/tds-api/pkg"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/labstack/echo/v4"
)

type Conferencer interface {
	GetAll() ([]models.Conference, error)
	GetById(targetId int) (*models.Conference, error)
	GetByName(targetName string) (*models.Conference, error)
	GetByGender(targetGender string) ([]models.Conference, error)
	GetByDivision(targetDivision pkg.Division) ([]models.Conference, error)
	GetByGenderAndDivision(targetGender string, targetDivision pkg.Division) ([]models.Conference, error)
}

type Conference struct {
	context echo.Context
}

func NewConference(context echo.Context) *Conference {
	context.Logger().Debug("NewConference")

	return &Conference{context: context}
}

// GetByDivision returns a slice of conferences in the specified division
func (c Conference) GetByDivision(targetDivision pkg.Division) ([]models.Conference, error) {
	var (
		err         error
		divisionUrl string
		conferences []models.Conference
	)

	c.context.Logger().Debug("GetByDivision")

	if targetDivision == pkg.DivisionUnknown {
		return []models.Conference{}, nil
	}

	if targetDivision == pkg.DivisionAll {
		return c.GetAll()
	}

	if divisionUrl, err = targetDivision.Url(); err != nil {
		return nil, fmt.Errorf("failed to retrieve conferences by targetDivision '%s': %v", targetDivision.String(), err.Error())
	}

	collector := colly.NewCollector()

	collector.OnHTML("table > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			decorator := api.Decorate(el, pkg.Prefix)

			conferenceUrl := decorator.GetLinkFromCell(1)
			conferences = append(conferences, models.Conference{
				Id:       api.GetIdFromUrl(conferenceUrl),
				Name:     decorator.GetTextFromCell(1),
				Url:      conferenceUrl,
				Division: pkg.DivisionToString(targetDivision),
				Gender:   pkg.GenderToString(api.GetGenderFromUrl(conferenceUrl)),
			})
		})
	})

	collector.OnError(func(r *colly.Response, err error) {
		c.context.
			Logger().
			Error(fmt.Sprintf("failed to retrieve conferences by targetDivision '%s' from '%s': %v", targetDivision.String(), r.Request.URL.String(), err.Error()))
	})

	if err = collector.Visit(divisionUrl); err != nil {
		c.context.
			Logger().
			Error(fmt.Sprintf("failed to retrieve conferences by targetDivision '%s' from '%s': %v", targetDivision.String(), divisionUrl, err.Error()))

		return nil, err
	}

	return conferences, nil
}

// GetAll returns all conferences
func (c Conference) GetAll() ([]models.Conference, error) {
	var (
		err                 error
		conferences         []models.Conference
		divisionConferences []models.Conference
	)
	c.context.Logger().Debug("GetAll")

	for _, division := range pkg.Divisions {
		if divisionConferences, err = c.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve conferences by division '%s': %v", division.String(), err.Error())
		}

		conferences = append(conferences, divisionConferences...)
	}

	return conferences, nil
}

// GetById returns a conference by id
func (c Conference) GetById(targetId int) (*models.Conference, error) {
	var (
		err                 error
		divisionConferences []models.Conference
	)

	c.context.Logger().Debug("GetById")

	for _, division := range pkg.Divisions {
		if divisionConferences, err = c.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve conferences by division '%s': %v", division.String(), err.Error())
		}

		for _, conference := range divisionConferences {
			if conference.Id == targetId {
				return &conference, nil
			}
		}
	}

	return nil, fmt.Errorf("conference with targetId '%d' not found", targetId)
}

// GetByName returns a conference by name
func (c Conference) GetByName(targetName string) (*models.Conference, error) {
	var (
		err                 error
		divisionConferences []models.Conference
	)

	c.context.Logger().Debug("GetByName")

	for _, division := range pkg.Divisions {
		if divisionConferences, err = c.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve conferences by division '%s': %v", division.String(), err.Error())
		}

		for _, conference := range divisionConferences {
			if conference.Name == targetName {
				return &conference, nil
			}
		}
	}

	return nil, fmt.Errorf("conference with targetName '%s' not found", targetName)
}

func (c Conference) GetByGender(targetGender pkg.Gender) ([]models.Conference, error) {
	var (
		err                 error
		conferences         []models.Conference
		divisionConferences []models.Conference
	)

	c.context.Logger().Debug("GetByGender")

	for _, division := range pkg.Divisions {
		if divisionConferences, err = c.GetByDivision(division); err != nil {
			return nil, fmt.Errorf("failed to retrieve conferences by division '%s': %v", division.String(), err.Error())
		}

		for _, conference := range divisionConferences {
			if conference.Gender == targetGender.String() {
				conferences = append(conferences, conference)
			}
		}
	}

	return conferences, nil
}

func (c Conference) GetByGenderAndDivision(targetGender pkg.Gender, targetDivision pkg.Division) ([]models.Conference, error) {
	var (
		divisionUrl string
		err         error
		conferences []models.Conference
	)

	c.context.Logger().Debug("GetByGenderAndDivision")

	if targetDivision == pkg.DivisionAll || targetDivision == pkg.DivisionUnknown {
		if targetGender == pkg.GenderBoth || targetGender == pkg.GenderUnknown {
			return c.GetAll()
		} else {
			return c.GetByGender(targetGender)
		}
	} else {
		if targetGender == pkg.GenderBoth || targetGender == pkg.GenderUnknown {
			return c.GetByDivision(targetDivision)
		}
	}

	if divisionUrl, err = targetDivision.Url(); err != nil {
		return nil, fmt.Errorf("failed to retrieve conferences by targetDivision '%s': %v", targetDivision.String(), err.Error())
	}

	collector := colly.NewCollector()

	collector.OnHTML("table > tbody", func(h *colly.HTMLElement) {
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

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err = collector.Visit(divisionUrl); err != nil {
		return nil, err
	}

	return conferences, nil
}
