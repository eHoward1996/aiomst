package fs

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eHoward1996/aiomst/db"

	"github.com/pascoej/gomusicbrainz"
)

// Tables to scan for specifed MBIDs
var tableNamesForScan = [3]string{"artists", "albums", "songs"}

// MusicBrainzScan scans the database and attempts to find MusicBrainz IDs
// for different objects
type MusicBrainzScan struct {
	sqlFile string
}

// SetSqlFile sets the path to the sql file for the python script to interact
// with the database.
func (fs *MusicBrainzScan) SetSqlFile(f string) {
	fs.sqlFile = f
}

// Scan calls a Python Script that looks for objects in the database that have 
// errored or otherwise non-unique MBIDs and attempts to grab the correct MBID
// from the MusicBrainz API
func (fs *MusicBrainzScan) Scan() {
	log.Print("FS: MusicBrainz Scan: Starting Scan...")
	startTime := time.Now()

	// cmdStr := fmt.Sprintf("python scripts/__main__.py %s", fs.sqlFile)
	// cmd := exec.Command("bash", "-c", cmdStr)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stdout

	// cmd.Run()
	client, err := gomusicbrainz.NewWS2Client(
    "https://musicbrainz.org/ws/2",
    "All In One Media Server Tool",
    "0.0.1-beta",
		"http://github.com/eHoward1996/aiomst",
		new(http.Client),
	)
	
	if err != nil {
		log.Printf("FS: MusicBrainz Scan: Error creating MB client: %v", err)
		return
	}

	handleAlbums(client)
	log.Printf(
		"FS: MusicBrainz Scan: Completed [time: %v]", time.Since(startTime))
}

func handleAlbums(c *gomusicbrainz.WS2Client) {
	albums, err := db.DB.GetBadAlbumMBIDs()
	if err != nil {
		log.Printf("FS: MusicBrainz Scan: Errored finding bad Album MBIDs: %v", err)
		return
	}

	for _, album := range albums {
		title := album.Title
		artist := album.Artist
		q := fmt.Sprintf(`artist:"%s" AND releaseaccent:"%s"`, artist, title)

		resp, err := c.SearchRelease(q, 10, 0)
		if err != nil {
			log.Printf(
				"FS: MusicBrainz Scan: Error searching for album: %v with %v - %v",
				err, artist, title,
			)
			continue
		}

		maxScore := -1
		var closestAlbum *gomusicbrainz.Release
		// var priority map[int]gomusicbrainz.Release = make(map[int]gomusicbrainz.Release)
		for _, rel := range resp.Releases {
			if _, member := permuteString(rel.Title)[strings.ToLower(title)]; !member {
				continue
			}

			if s := scoreAlbum(album, rel, c); s > maxScore {
				maxScore = s
				closestAlbum = rel
			}
		}

		album.MBID = string(closestAlbum.Id())
		if err := db.DB.UpdateAlbum
	}
}

func permuteString(s string) map[string]struct{} {
	var sep string = ""
	s = strings.ToLower(s)
	permutes := make(map[string]struct{})
	
	if strings.Contains(s, " & ") {
		sep = " & "
		permutes[strings.ReplaceAll(s, " & ", " and ")] = struct{}{}
	}

	if strings.Contains(s, " and ") {
		sep = " and "
		permutes[strings.ReplaceAll(s, " and ", " & ")] = struct{}{}
	}

	addRemove := func(x string) []string {
		l := []string{}
		if x[:4] == "the " { 
			l = append(l, x[:4])
		} else if x[:2] != "a " {
			l = append(l, "the " + x) 
		} 

		if x[:2] == "a " {
			l = append(l, x[2:])
		} else if x[:4] != "the " {
			l = append(l, "a " + x)
		}

		if strings.Contains(x, " vol. ") {
			l = append(l, strings.ReplaceAll(x, " vol. ", " volume "))
		} else if strings.Contains(x, " volume ") {
			l = append(l, strings.ReplaceAll(x, " volume ", " vol. "))
		}
		return l
	}

	if sep != "" {
		for _, x := range strings.Split(s, sep) {
			permutes[x] = struct{}{}
			for _, y := range addRemove(x) {
				permutes[y] = struct{}{}
			}
		}
	} else {
		for _, x := range addRemove(s) {
			permutes[x] = struct{}{}
		}
	}
	return permutes
}

func scoreAlbum(
		album db.Album,
		mbResp *gomusicbrainz.Release,
		c *gomusicbrainz.WS2Client) int {
	
			var wg sync.WaitGroup
			ch := make(chan int)
			score := 0
		
			calculateDisambigScore := func(albumTitle, mbDisambig string) {
				defer wg.Done()
				
				disambigScore := 0
				var disambigPriority map[string]int = make(map[string]int)
				disambigPriority["deluxe"] = 4
				disambigPriority["explicit"] = 3
				disambigPriority["mastered for itunes"] = 2
				disambigPriority["bonus tracks"] = 1
				
				albumDisambig := ""
				openParen := strings.LastIndex(albumTitle, "(")
				closeParen := strings.LastIndex(albumTitle, ")")
				if openParen != -1 && closeParen != -1 {
					albumDisambig = strings.ToLower(albumTitle[openParen + 1:closeParen])
				}
				
				mbDisLower := strings.ToLower(mbDisambig)
				regexMatch, _ := regexp.MatchString(
					"(.* edition|.* version)", mbDisLower)
				switch {
				case albumDisambig == mbDisLower:
					disambigScore += 20
				case regexMatch:
					disambigScore += 10
				default:
					if val, ok := disambigPriority[mbDisLower]; ok {
						disambigScore += val
					}
				}

				ch <- disambigScore
			}
			
			calculateArtistsScore := func(
				albumArtist string, mbArtists []gomusicbrainz.NameCredit) {
					
					defer wg.Done()
					artistsScore := 0
					permutedAlbumArtist := permuteString(albumArtist)
					for _, mbArtist := range mbArtists {
						if _, ok := permutedAlbumArtist[mbArtist.Artist.Name]; ok {
							artistsScore += 20
						}
					}

					ch <- artistsScore
			}

			calculateTracksScore := func(album db.Album, relMBID gomusicbrainz.MBID) {
				defer wg.Done()
				tracksScore := 0
				mbResp, err := c.LookupRelease(relMBID, "media", "recordings")
				if err != nil {
					log.Printf(
						"FS: MusicBrainz Scan: Error while looking up Album %v - %v. %v",
						album.Title, album.Artist, err,
					)
					ch <- tracksScore
					return
				}

				songsOnAlbum, err := db.DB.SongsForAlbum(album.ID)
				if err != nil {
					log.Printf("FS: MusicBrainz Scan: Error while queryign DB: %v", err)
				}

				mediumDiscNum := 1
				for _, mediums := range mbResp.Mediums {	
					for i, track := range mediums.TrackList.Tracks {
						songTitle := songsOnAlbum[i].Title
						if track.Recording.Title == songTitle {
							tracksScore += 5
						} else if strings.Contains(songTitle, track.Recording.Title) {
							tracksScore += 2
						}

						songDiscStr := strings.Split(songsOnAlbum[i].Disc, "/")[0]
						songDiscNum, err := strconv.Atoi(songDiscStr)
						if err != nil {
							songDiscNum = 0
						}
						if mediumDiscNum == songDiscNum {
							tracksScore++
						}
					}
					mediumDiscNum++
				}
				ch <- tracksScore
			}

			wg.Add(3)
			go calculateDisambigScore(album.Title, mbResp.Disambiguation)
			go calculateArtistsScore(album.Artist, mbResp.ArtistCredit.NameCredits)
			go calculateTracksScore(album, mbResp.Id())
			score += <- ch
			
			wg.Wait()
			close(ch)
			
			return score
}