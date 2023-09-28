package api

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jedi-knights/tds-api/pkg"
	"strconv"
	"strings"
)

type HTMLElementDecorator struct {
	element *colly.HTMLElement
	prefix  string
}

func Decorate(element *colly.HTMLElement, prefix string) *HTMLElementDecorator {
	return &HTMLElementDecorator{element, prefix}
}

func (d *HTMLElementDecorator) GetLinkFromCell(index int) string {
	var link string

	goquerySelector := fmt.Sprintf("td:nth-child(%d) > a", index)

	links := d.element.ChildAttrs(goquerySelector, "href")
	if len(links) > 0 {
		link = d.prefix + links[0]
	} else {
		link = ""
	}

	return link
}

func (d *HTMLElementDecorator) GetTextFromCell(index int) string {
	goquerySelector := fmt.Sprintf("td:nth-child(%d)", index)

	text := d.element.ChildText(goquerySelector)
	text = NormalizeText(text)

	return text
}

func NormalizeText(text string) string {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	text = strings.Trim(text, " ")

	return text
}

func GetIdFromUrl(url string) int {
	index := strings.LastIndex(url, "-")

	if index < 0 {
		return -1
	}

	suffix := url[index+1:]
	id, err := strconv.Atoi(suffix)

	if err != nil {
		return -1
	}

	return id
}

func GetGenderFromUrl(url string) pkg.Gender {
	if strings.Contains(url, "/men/") {
		return pkg.GenderMale
	} else if strings.Contains(url, "/women/") {
		return pkg.GenderFemale
	} else {
		return pkg.GenderBoth
	}
}
