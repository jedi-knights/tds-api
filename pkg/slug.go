package pkg

import (
	"regexp"
	"strings"
)

func CreateSlug(input string) string {
	// Convert the string to lowercase
	input = strings.ToLower(input)

	// Replace spaces with hyphens
	input = strings.ReplaceAll(input, " ", "-")

	// Use a regular expression to remove non-alphanumeric characters and hyphens
	reg, _ := regexp.Compile("[^a-zA-Z0-9-]+")
	input = reg.ReplaceAllString(input, "")

	// Trim any leading or trailing hyphens
	input = strings.Trim(input, "-")

	return input
}
