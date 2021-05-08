package integrated

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/divan/num2words"
	"github.com/irlndts/go-discogs"
	"github.com/pascoej/gomusicbrainz"
)

// Object representing album data sent to methods to find their third party ids.
type sendAlbum struct {
	ID         int
	MetadataID int
	Artist     string
	Title      string
}

type sendArtist struct {}

type mbArtistResponse struct {}

// Object that holds a musicbrainz response. Either a Album or Artist type.
type musicbrainzResponse struct {
	AlbumResponse  *gomusicbrainz.Release
	ArtistResponse *mbArtistResponse
}

type discogsArtistResponse struct {}

// Object that holds a discogs response. Either a Album or Artist type.
type discogsResponse struct {
	AlbumResponse  *discogs.Master
	ArtistResponse *discogsArtistResponse
}

type recvAlbum struct {
	AlbumID       int
	Album         string
	Artist        string
	MetadataID    int
	MBRelease     *gomusicbrainz.Release
	DiscogsMaster *discogs.Master
}

const (
	maxRoutines       = 5
	maxRequestRetries = 20
	maxWorkers        = 5
	requestRetryTimer = 3
	discogsPerPage    = 25
	mbLimit           = 10
	mbOffset          = -1
)

func replaceSpecials(s string) string {
	var r string = s
	r = strings.ReplaceAll(r, "‘", "'")
	r = strings.ReplaceAll(r, "’", "'")
	r = strings.ReplaceAll(r, "′", "'")
	r = strings.ReplaceAll(r, "“", `"`)
	r = strings.ReplaceAll(r, "″", `"`)
	r = strings.ReplaceAll(r, "“", `"`)
	r = strings.ReplaceAll(r, "”", `"`)

	r = strings.ReplaceAll(r, "…", "...")
	
	r = strings.ReplaceAll(r, "-", "-")
	r = strings.ReplaceAll(r, "‐", "-")
	r = strings.ReplaceAll(r, "−", "-")
	r = strings.ReplaceAll(r, "‒", "-")
	r = strings.ReplaceAll(r, "–", "-")
	r = strings.ReplaceAll(r, "―", "-")
	r = strings.ReplaceAll(r, "—", "-")
	return r
}

// permuteString takes a string and returns a map contaioning different
// permutations of the string. The entries are calculated by manipulating the 
// string by adding, removing and replacing parts of the string. For example, 
// the string "Happy & Good" would return a set returning "Happy and Good".
func permuteString(s string) map[string]struct{} {
	var sep string = " & "
	s = strings.ToLower(s)

	permutes := make(map[string]struct{})
	permutes[s] = struct{}{}

	s = replaceSpecials(s)
	permutes[s] = struct{}{}
	permutes[strings.ReplaceAll(s, " & ", " and ")] = struct{}{}
	permutes[strings.ReplaceAll(s, " and ", " & ")] = struct{}{}
	permutes[strings.ReplaceAll(s, " vol. ", " volume ")] = struct{}{}
	permutes[strings.ReplaceAll(s, " volume ", " vol. ")] = struct{}{}
	permutes[strings.ReplaceAll(s, " pt. ", " part ")] = struct{}{}
	permutes[strings.ReplaceAll(s, " part ", " pt. ")] = struct{}{}

	n := strings.ReplaceAll(s, "(", "[")
	n = strings.ReplaceAll(n, ")", "]")
	permutes[n] = struct{}{}

	n = strings.ReplaceAll(s, "[", "(")
	n = strings.ReplaceAll(n, "]", ")")
	permutes[n] = struct{}{}

	matches, _ := regexp.MatchString(`\s+/\s+`, s)
	if matches {
		slashes := regexp.MustCompile(`\s+/\s+`)
		permutes[slashes.ReplaceAllString(s, "/")] = struct{}{}
	}	
	n = strings.ReplaceAll(s, "/", " / ")
	permutes[n] = struct{}{}

	if strings.Contains(s, " and ") {
		sep = " and "
	}

	addRemove := func(x string) []string {
		l := []string{}
		sizeThe := len("the ")
		sizeA := len("a ")
		the := "the "
		a := "a "

		if len(x) > sizeThe && x[:sizeThe] == the { 
			l = append(l, x[sizeThe:])
		} else if len(x) > sizeA && x[:sizeA] != a {
			l = append(l, the + x) 
		} 

		if len(x) > sizeA && x[:sizeA] == a {
			l = append(l, x[sizeA:])
		} else if len(x) > sizeThe && x[:sizeThe] != the {
			l = append(l, a + x)
		}
		return l
	}

	for _, x := range strings.Split(s, sep) {
		permutes[x] = struct{}{}

		for _, y := range addRemove(x) {
			permutes[y] = struct{}{}
		}
	}
	
	permutes[stripAllParentheses(s)] = struct{}{}
	permutes[convertNumbersToWords(s)] = struct{}{}
	return permutes
}

func stripLastParentheses(s string) string {
	r := s
	open := strings.LastIndex(s, "(")
	close := strings.LastIndex(s, ")")

	if open != -1 && close != -1 {
		close++
		if close == len(s) {
			return strings.Trim(s[:open], " ")
		} else if close > open {
			l := strings.Trim(s[:open], " ")
			r := strings.Trim(s[close:], " ")
			return strings.Join([]string{l, r}, " ")
		}
	}
	return r
}

func stripAllParentheses(s string) string {
	regex := regexp.MustCompile("(\\(.*\\)|\\[.*\\]|\\{.*\\})")
	matches := regex.FindAllString(s, -1)

	for _, match := range matches {
		s = strings.Replace(s, match, "", 1)
	}
	return strings.Trim(s, " ")
}

func convertNumbersToWords(s string) string {
	regex := regexp.MustCompile("[0-9]+")
	matches := regex.FindAllString(s, -1)
	
	for _, match := range matches {
		n, err := strconv.Atoi(match)
		if err == nil {
			n2w := num2words.Convert(n)
			s = strings.Replace(s, match, n2w, 1)
		}
	}
	return s
}

// isSimilarString takes two strings and compares them to see if they are the 
// same string
func isSimilarString(s1, s2 string) bool {
	s1Lower := strings.Trim(strings.ToLower(s1), " ")
	s1Stripped := stripAllParentheses(s1Lower)
	s1NumToWord := convertNumbersToWords(s1Lower)
	
	permutedStringMap := permuteString(s2)	
	_, regTitle := permutedStringMap[s1Lower]
	_, ignParen := permutedStringMap[s1Stripped]
	_, convNums := permutedStringMap[s1NumToWord]
	_, combined := permutedStringMap[
		stripAllParentheses(convertNumbersToWords(s1))]
	return regTitle || ignParen || convNums || combined
}

// repeatRequest repeats a request "maxRequestRetries" number of times. If it
// succeeds it instantly returns the result. If it errors more than 
// "maxRequestRetries" number of times nil as well as the error is returned.
func repeatRequest(req func() (interface{}, error)) (interface{}, error) {
	var e error
	for retry := 0; retry <= maxRequestRetries; retry++ {
		d := requestRetryTimer * retry
		time.Sleep(time.Duration(d) * time.Second)

		resp, err := req()
		if err == nil {
			return resp, nil
		}
		e = err
	}
	return nil, fmt.Errorf("Request failed: %v", e)
}

func formatOutput(s string, length int) string {
	if len(s) > length {
		l := length - 3
		s = s[:l] + "..."
	} else if len(s) < length {
		for {
			if len(s) == 15 {
				break
			}
			s += " "
		}
	}

	return s
}