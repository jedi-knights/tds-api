package pkg

import "fmt"

type Division int

const (
	DivisionAll = iota
	DivisionDI
	DivisionDII
	DivisionDIII
	DivisionNAIA
	DivisionNJCAA
	DivisionUnknown
)

var Divisions = []Division{
	DivisionDI,
	DivisionDII,
	DivisionDIII,
	DivisionNAIA,
	DivisionNJCAA,
}

var divisionToUrlMapping = map[Division]string{
	DivisionDI:    "https://www.topdrawersoccer.com/college-soccer/college-conferences/di/divisionid-1",
	DivisionDII:   "https://www.topdrawersoccer.com/college-soccer/college-conferences/dii/divisionid-2",
	DivisionDIII:  "https://www.topdrawersoccer.com/college-soccer/college-conferences/diii/divisionid-3",
	DivisionNAIA:  "https://www.topdrawersoccer.com/college-soccer/college-conferences/naia/divisionid-4",
	DivisionNJCAA: "https://www.topdrawersoccer.com/college-soccer/college-conferences/njcaa/divisionid-5",
}

func (d Division) String() string {
	return DivisionToString(d)
}

func (d Division) Url() (string, error) {
	var (
		ok          bool
		divisionUrl string
	)

	if d == DivisionAll {
		return "", fmt.Errorf("cannot get url for DivisionAll")
	}

	if divisionUrl, ok = divisionToUrlMapping[d]; !ok {
		return "", fmt.Errorf("unsupported division for url mapping: %v", d)
	}

	return divisionUrl, nil
}

func DivisionToString(division Division) string {
	switch division {
	case DivisionDI:
		return "di"
	case DivisionDII:
		return "dii"
	case DivisionDIII:
		return "diii"
	case DivisionNAIA:
		return "naia"
	case DivisionNJCAA:
		return "njcaa"
	default:
		return "all"
	}
}

func StringToDivision(division string) Division {
	switch division {
	case "":
		return DivisionAll
	case "all":
		return DivisionAll
	case "di":
		return DivisionDI
	case "dii":
		return DivisionDII
	case "diii":
		return DivisionDIII
	case "naia":
		return DivisionNAIA
	case "njcaa":
		return DivisionNJCAA
	default:
		return DivisionUnknown
	}
}
