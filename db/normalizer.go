package db

import (
	"log"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// normalizeString removes accent marks from strings.
func normalizeString(s string) string 	{
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		log.Printf("DB: Normalize String: %s with value %s", err, s)
		return s
	}

	return result
}