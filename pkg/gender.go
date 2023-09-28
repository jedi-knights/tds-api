package pkg

type Gender int

const (
	GenderBoth = iota
	GenderFemale
	GenderMale
	GenderUnknown
)

func GenderToString(gender Gender) string {
	switch gender {
	case GenderBoth:
		return "both"
	case GenderFemale:
		return "female"
	case GenderMale:
		return "male"
	default:
		return "both"
	}
}

func StringToGender(gender string) Gender {
	if gender == "" {
		return GenderBoth
	}
	
	switch gender {
	case "both":
		return GenderBoth
	case "female":
		return GenderFemale
	case "male":
		return GenderMale
	default:
		return GenderUnknown
	}
}
