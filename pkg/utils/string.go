package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// NormalizeCity standardizes city name to Title Case
// Trims whitespace and applies proper capitalization
func NormalizeCity(city string) string {
	if city == "" {
		return ""
	}
	trimmed := strings.TrimSpace(city)
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(trimmed))
}
