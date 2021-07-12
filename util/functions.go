package util

import (
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var Logger *log.Logger

// ExpandHomeDir replaces input tilde characters with the absolute path to the current
// user's home directory.
func ExpandHomeDir(path string) string {
	return strings.Replace(path, "~", System.User.HomeDir, -1)
}

// UNIXtoRFC1123 transforms an input UNIX timestamp into the form specified by RFC1123,
// using the GMT time zone. This function is used to output Last-Modified headers via HTTP.
func UNIXtoRFC1123(unix int64) string {
	return strings.Replace(time.Unix(unix, 0).UTC().Format(time.RFC1123), "UTC", "GMT", 1)
}

// NormalizeString removes accent marks from strings.
func NormalizeString(s string) string 	{
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		Logger.Printf("DB: Normalize String: %s with value %s", err, s)
		return s
	}

	return result
}

func InitializeLogger() {
	Logger = log.Default()
	Logger.SetOutput(os.Stdout)
}