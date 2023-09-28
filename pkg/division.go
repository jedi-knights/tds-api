package pkg

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
	if division == "" {
		return DivisionAll
	}

	switch division {
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
