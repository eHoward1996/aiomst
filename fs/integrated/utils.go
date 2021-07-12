package integrated

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/divan/num2words"
	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"
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

type sendArtist struct {
	ID 				 int
	Title      string
	MetadataID int
	AlbumMBID  string
	DiscogsID  int
}

// Object that holds a musicbrainz response. Either a Album or Artist type.
type musicbrainzResponse struct {
	AlbumResponse  *gomusicbrainz.Release
	ArtistResponse *gomusicbrainz.Artist
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
	ArtistID      int
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


func printReceivedAlbum(count, numAlbums int, result *recvAlbum) {	
	mbId := "errored"
	dsId := -1
	if result.MBRelease != nil {
		mbId = string(result.MBRelease.ID)
	}
	if result.DiscogsMaster != nil {
		dsId = result.DiscogsMaster.ID
	}

	artistString := formatOutput(result.Artist, 18)
	albumString := formatOutput(result.Album, 18)
	processedStr := fmt.Sprintf(
		"(%v of %v) %v - %v MBID: %v\tDiscogsID: %v", count, numAlbums,
		artistString, albumString, formatOutput(mbId, 13), dsId)
	util.Logger.Print(processedStr)
}

func updateAlbum(recv *recvAlbum) error {
	album := &db.Album{ID: recv.AlbumID}
	if err := db.DB.LoadAlbum(album); err != nil {
		return err
	}

	if recv.MBRelease != nil {
		album.MBID = string(recv.MBRelease.ID)
	}
	if recv.DiscogsMaster != nil {
		album.DiscogsID = recv.DiscogsMaster.ID
	}
	
	return album.Update()
}

func updateArtist(recv *recvAlbum) error {
	artist := &db.Artist{ID: recv.ArtistID}
	if err := db.DB.LoadArtist(artist); err != nil {
		return err
	}

	if recv.MBRelease != nil && recv.MBRelease.ArtistCredit != nil {
		mbidList := make([]string, 0)
		for _, aCredit := range recv.MBRelease.ArtistCredit.NameCredits {
			mbArtist := aCredit.Artist
			if isSimilarString(mbArtist.Name, artist.Title) {
				mbidList = append(mbidList, string(mbArtist.ID))
			}
		}
		mbidString := strings.Join(mbidList, ",")
		if mbidString != "" {
			artist.MBID = mbidString
		}
	}

	if recv.DiscogsMaster != nil {
		discList := make([]string, 0)
		for _, aSource := range recv.DiscogsMaster.Artists {
			if isSimilarString(aSource.Name, artist.Title) {
				discList = append(discList, strconv.Itoa(aSource.ID))
			}
		}
		if len(discList) != 0 {
			artist.DiscogsID = strings.Join(discList, ",")
		}
	}

	return artist.Update()
}

func (i TPIntegrator) createAlbumMetadata(recv *recvAlbum) {
	albumObj := &db.Album{ID: recv.AlbumID}
	if err := albumObj.Load(); err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error loading album with ID: %v",
			albumObj.ID)
		return
	}

	md, err := albumObj.GetMetadataObj()
	if err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error loading metadata with ID: %v",
			albumObj.MetadataID)
		return
	}

	mdOnFileAsObj := new(db.AlbumMetadata)
	md.ToStruct(mdOnFileAsObj)
	mMD := mdOnFileAsObj.MusicBrainz
	dMD := mdOnFileAsObj.Discogs

	if mMD.Release.Title == "" && recv.MBRelease != nil {
		mMD = buildMBAlbumMetadata(recv.MBRelease)	
	}
	if dMD.Title == "" && recv.DiscogsMaster != nil {
		dMD = buildDiscogsAlbumMetadata(recv.DiscogsMaster)
	}

	albumMD := db.AlbumMetadata{
		AlbumName: recv.Album,
		MusicBrainz: mMD,
		Discogs: dMD,
	}
	writeMetadata(recv.MetadataID, albumMD)
}

func (i TPIntegrator) createArtistMetadata(recv *recvAlbum) {
	artistObj := &db.Artist{ID: recv.ArtistID}
	if err := artistObj.Load(); err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error loading artist with ID: %v",
			artistObj.ID)
			return
	}

	md, err := artistObj.GetMetadataObj()
	if err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error loading metadata with ID: %v",
			artistObj.MetadataID)
			return
	}

	mdOnFileAsObj := new(db.ArtistMetadata)
	md.ToStruct(mdOnFileAsObj)
	mMD := mdOnFileAsObj.MusicBrainz
	dMD := mdOnFileAsObj.Discogs

	mbArtistLen := len(mMD.Artists)
	diArtistLen := len(dMD.Artists)

	if mbArtistLen == 0 && recv.MBRelease != nil {
		mMD = buildMBArtistMetadata(recv.Artist, recv.MBRelease)
	}
	if diArtistLen == 0 && recv.DiscogsMaster != nil {
		discMaster := *recv.DiscogsMaster
		dMD = buildDiscogsArtistMetadata(discMaster.Artists)
	}

	artistMD := db.ArtistMetadata{
		ArtistName: recv.Artist,	
		MusicBrainz: mMD,
		Discogs: dMD,
	}
	writeMetadata(artistObj.MetadataID, artistMD)
}

func writeMetadata(metadataID int, md interface{}) {
	bytes, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error writing Metadata file: %v", err)
		return
	}

	tmp := new(db.Metadata)
	tmp.ID = metadataID
	if err := tmp.Load(); err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error loading Metadata object: %v", err)
		return
	}
	if err := ioutil.WriteFile(tmp.Path, bytes, 0644); err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error writing Metadata File: %v", err)
		return
	}
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