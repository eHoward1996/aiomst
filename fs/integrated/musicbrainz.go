package integrated

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gammazero/workerpool"
	"github.com/pascoej/gomusicbrainz"
)

// getMusicBrainzID makes a request to the MusicBrainz API and then parses
// responses to find the best matching MBID
func (i TPIntegrator) getMusicBrainzData(
	obj interface{}) (musicbrainzResponse, error) {	
	
		var resp musicbrainzResponse = musicbrainzResponse{}
		var err error = nil

		if mbClient == nil {
			return resp, fmt.Errorf("No MusicBrainz Client available")
		}

		switch t := obj.(type) {
		case *sendAlbum:
			resp.AlbumResponse, err = findClosestAlbum(t)
		// case gomusicbrainz.MBID:
		// 	resp.ArtistResponse, err = mbLookupArtist(fmt.Sprintf("%v", obj))
		}	

		return resp, err
}

func findClosestAlbum(album *sendAlbum) (*gomusicbrainz.Release, error) {
	releaseList, err := getReleaseList(album)
	if err != nil {
		return nil, err
	}
	
	type result struct {
		release *gomusicbrainz.Release
		score int
	}

	wp := workerpool.New(maxWorkers)
	max := &result{score: -1}
	resultCh := make(chan *result, len(releaseList))		
	for _, release := range releaseList {
		release := release
		wp.Submit(func() {
			score := scoreAlbum(album, release)
			r := new(result)
			r.score = score
			r.release = release
			resultCh <- r
		})
	}
	
	go func() {
		wp.StopWait()
		close(resultCh)
	}()

	for r := range resultCh {
		if r.score > max.score {
			max = r
		}
	}

	if max.score <= 0 {
		return nil, fmt.Errorf("MusicBrainz: Scores were at or below 0")
	}

	return max.release, nil
}

func getReleaseList(a *sendAlbum) ([]*gomusicbrainz.Release, error) {
	title := stripAllParentheses(a.Title)
	q := fmt.Sprintf(`releasegroupaccent:"%s" OR release:"%s" OR artist:"%s"`,
		title, a.Title, a.Artist)

	resp, err := repeatRequest(func() (interface{}, error) {
		r, err := mbClient.SearchReleaseGroup(q, mbLimit, mbOffset)
		if r != nil {
			return r.ReleaseGroups, err
		}
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	var releases []*gomusicbrainz.Release = make([]*gomusicbrainz.Release, 0)
	rg := resp.([]*gomusicbrainz.ReleaseGroup)
	for _, group := range rg {
		for _, release := range group.Releases.Releases {
			if !isSimilarString(a.Title, release.Title) {
				continue
			}

			rel, err := repeatRequest(func() (interface{}, error) {
				return mbClient.LookupRelease(
					release.ID,
					[]string{"recordings", "artists", "release-groups", "url-rels"}...)
			})
			if err != nil {
				return nil, err
			}
			releases = append(releases, rel.(*gomusicbrainz.Release))
		}
	}
	return releases, nil
}

func scoreAlbum(a *sendAlbum, mbRelease *gomusicbrainz.Release) int {
	var wg sync.WaitGroup
	scoreCh := make(chan int, 3)
	score := 0

	wg.Add(3)
	mbArtists := mbRelease.ArtistCredit.NameCredits

	go a.calculateDisambiguationScore(mbRelease.Disambiguation, scoreCh, &wg)
	go a.calculateArtistsScore(mbArtists, scoreCh, &wg)
	go a.calculateTracksScore(mbRelease.Mediums, scoreCh, &wg)

	go func() {
		wg.Wait()
		close(scoreCh)
	}()

	for s := range scoreCh {
		score += s
	}	
	return score
}

func (a *sendAlbum) calculateDisambiguationScore(
	mbDisambig string, ch chan int, wg *sync.WaitGroup) {

		defer wg.Done()

		disambigScore := 0
		disambigPriority := map[string]int{
			"deluxe": 4,
			"explicit": 3,
			"mastered for itunes": 2,
			"bonus tracks": 1,
		}

		albumDisambig := ""
		openParen := strings.LastIndex(a.Title, "(")
		closeParen := strings.LastIndex(a.Title, ")")
		if openParen != -1 && closeParen != -1 {
			albumDisambig = strings.ToLower(a.Title[openParen + 1:closeParen])
		}

		mbDisambigLower := strings.ToLower(mbDisambig)
		mbRegexMatch, _ := regexp.MatchString(
			"(.* edition|.* version|.* exclusive)", mbDisambigLower)
		abRegexMatch, _ := regexp.MatchString(
			"(.* edition|.* version|.* exclusive)", albumDisambig)

		switch {
		case albumDisambig == mbDisambigLower:
			disambigScore = 20
		case mbRegexMatch && abRegexMatch:
			disambigScore = 5
		default:
			if val, ok := disambigPriority[mbDisambigLower]; ok {
				disambigScore = val
			}
		}
		ch <- disambigScore
}

func (a *sendAlbum) calculateArtistsScore(
	mbArtists []gomusicbrainz.NameCredit, ch chan int, wg *sync.WaitGroup) {

		defer wg.Done()
		
		artistsScore := 0
		for _, mbArtist := range mbArtists {
			aLower := strings.ToLower(mbArtist.Artist.Name)
			if isSimilarString(aLower, a.Artist) {
				artistsScore += 50
			} else {
				artistsScore -= 5
			}
		}
		ch <- artistsScore
}

func (a *sendAlbum) calculateTracksScore(
	mediums []*gomusicbrainz.Medium, ch chan int, wg *sync.WaitGroup) {

		defer wg.Done()

		tracksScore := 0	
		songsOnAlbum, err := db.DB.SongsForAlbum(a.ID)
		if err != nil {
			ch <- tracksScore
			return
		}

		type mbTrack struct {
			title string
			number string
			position int
			disc int
		}
		
		var trackList map[string]mbTrack = make(map[string]mbTrack)
		for mediumDiscNum, m := range mediums {
			for _, track := range m.TrackList.Tracks {
				key := fmt.Sprintf("%v:%v", mediumDiscNum + 1, track.Position)
				mbTitle := track.Title
				if mbTitle == "" {
					mbTitle = track.Recording.Title
				}

				value := mbTrack{
					title: mbTitle,
					number: track.Number,
					position: track.Position,
					disc: mediumDiscNum+1,
				}
				trackList[key] = value
			}
		}

		for _, song := range songsOnAlbum {
			songDiscStr := strings.Split(song.Disc, "/")[0]
			songDisc, err := strconv.Atoi(songDiscStr)
			if err != nil {
				songDisc = 1
			}

			key := fmt.Sprintf("%v:%v", songDisc, song.Track)
			if track, ok := trackList[key];
				ok && isSimilarString(song.Title, track.title) {
					tracksScore += 5
			} else {
				for _, t := range trackList {	
					if isSimilarString(song.Title, t.title) {
						tracksScore += 2
						break
					}
				}
			}
		}
		if tracksScore == 0 {
			tracksScore = -1000
		}
		ch <- tracksScore
}

func buildMBAlbumMetadata(rel *gomusicbrainz.Release) db.MusicBrainzMetadata {
	var mbmd db.MusicBrainzMetadata = db.MusicBrainzMetadata{}
	if rel == nil {
		return mbmd
	}

	var artists []db.MusicBrainzArtist = make([]db.MusicBrainzArtist, 0)
	var artistCredits gomusicbrainz.ArtistCredit
	if rel.ArtistCredit != nil {
		artistCredits = *rel.ArtistCredit
		for _, artist := range artistCredits.NameCredits {
			x := new(db.MusicBrainzArtist)
			x.ID = string(artist.Artist.ID)
			x.Name = artist.Artist.Name
			x.Disambiguation = artist.Artist.Disambiguation
			x.SortName = artist.Artist.SortName
			x.Type = artist.Artist.Type
			
			aliases := make([]db.MusicBrainzAlias, 0)
			for _, a := range artist.Artist.Aliases {
				m := db.MusicBrainzAlias{
					Type: a.Type,
					Name: a.Name,
				}
				aliases = append(aliases, m)
			}
			x.Aliases = aliases

			x.Country = artist.Artist.CountryCode
			x.Area = "Unknown"
			if artist.Artist.Area != nil {
				x.Area = artist.Artist.Area.Name
			}
			artists = append(artists, *x)
		}
	}
	mbmd.Artists = artists
	
	release := db.MusicBrainzRelease{
		ID: string(rel.ID),
		Title: rel.Title,
		Disambiguation: rel.Disambiguation,
		Status: "Unknown",
	}
	if rel.Status != nil {
		release.Status = rel.Status.Status
	}

	totalTracks := 0
	discs := 0
	tracks := make([]db.MusicBrainzTrack, 0)
	for _, disc := range rel.Mediums {
		for _, track := range disc.TrackList.Tracks {
			m := db.MusicBrainzTrack{
				ID: string(track.ID),
				Title: track.Recording.Title,
				Number: track.Position,
				Length: track.Length,
			}
			tracks = append(tracks, m)
		}
		totalTracks += disc.TrackList.Count
		discs++
	}
	release.Tracks = tracks
	release.TrackCount = totalTracks
	release.DiscCount = discs
	mbmd.Release = release

	releaseGroups := make([]db.MusicBrainzReleaseGroup, 1)
	if rel.ReleaseGroup != nil {
		r := *rel.ReleaseGroup
		rg := db.MusicBrainzReleaseGroup{
			ID: string(r.ID),
			Title: r.Title,
			Type: r.Type,
		}
		releaseGroups[0] = rg
	}
	mbmd.ReleaseGroups = releaseGroups

	urls := make([]db.RelatedUrl, 0)
	if _, ok := rel.Relations["url"]; rel.Relations != nil && ok {
		for _, x := range rel.Relations["url"] {
			url := x.(*gomusicbrainz.URLRelation).RelationAbstract
			m := db.RelatedUrl{
				Type: url.Type,
				Url: url.Target,
			}
			urls = append(urls, m)
		}
	}
	mbmd.RelatedUrls = urls
	return mbmd
}

func buildMBArtistMetadata(
	rArtist string, rel *gomusicbrainz.Release) db.MusicBrainzMetadata {

		var mbmd db.MusicBrainzMetadata = db.MusicBrainzMetadata{}
		mbmd.Artists = make([]db.MusicBrainzArtist, 0)
		mbmd.AssociatedActs = make([]db.AssociatedAct, 0)
		mbmd.RelatedUrls = make([]db.RelatedUrl, 0)
		mbmd.ReleaseGroups = make([]db.MusicBrainzReleaseGroup, 0)
		mbmd.RelatedTags = make([]string, 0)
		tagMap := make(map[string]struct{})

		for _, artist := range rel.ArtistCredit.NameCredits {
			if strings.Contains(rArtist, artist.Artist.Name) {
				if artistData := mbLookupArtist(artist.Artist.ID); artistData != nil {
					aData := db.MusicBrainzArtist{
						ID: string(artistData.ID),
						Name: artistData.Name,
						Disambiguation: artistData.Disambiguation,
						SortName: artistData.SortName,
						Type: artistData.Type,
						Aliases: []db.MusicBrainzAlias{},
						Area: artistData.Area.Name,
						Country: artistData.CountryCode,
					}
					for _, val := range artistData.Aliases {
						aData.Aliases = append(aData.Aliases, db.MusicBrainzAlias{
							Name: val.Name,
							Type: val.Type,
						})
					}
					mbmd.Artists = append(mbmd.Artists, aData)

					if artistData.Relations["artist"] != nil {
						for _, val := range artistData.Relations["artist"] {
							aRel := val.(*gomusicbrainz.ArtistRelation).Artist
							associated := db.AssociatedAct{
								ID: string(aRel.ID),
								Name: aRel.Name,
								Disambiguation: aRel.Disambiguation,
								Type: aRel.Type,
								Relation: val.(*gomusicbrainz.ArtistRelation).Type,
							}
							mbmd.AssociatedActs = append(mbmd.AssociatedActs, associated)
						}
					}
					
					if artistData.Relations["url"] != nil {
						for _, val := range artistData.Relations["url"] {
							urlRel := val.(*gomusicbrainz.URLRelation).RelationAbstract
							mbmd.RelatedUrls = append(mbmd.RelatedUrls, db.RelatedUrl{
								Type: urlRel.Type,
								Url: urlRel.Target,
							})
						}
					}
					 
					if artistData.ReleaseGroups != nil {
						relGroups := *artistData.ReleaseGroups
						for _, rg := range relGroups.ReleaseGroups {
							rCount := 0
							if rg.Releases != nil {
								rCount = rg.Releases.Count
							}
							
							mbmd.ReleaseGroups = append(
								mbmd.ReleaseGroups,
								db.MusicBrainzReleaseGroup{
									ID: string(rg.ID),
									Title: rg.Title,
									Type: rg.Type,
									ReleaseCount: rCount,
									ReleaseDate: rg.FirstReleaseDate.Time,
								},
							)
						}
					}

					for _, tag := range artistData.Tags {
						tagMap[tag.Name] = struct{}{}
					}
				}
			}
		}

		for k := range tagMap {
			mbmd.RelatedTags = append(mbmd.RelatedTags, k)
		}
		return mbmd
}

func mbLookupArtist(mbid gomusicbrainz.MBID) *gomusicbrainz.Artist {
	resp, err := repeatRequest(func() (interface{}, error) {
		return mbClient.LookupArtist(mbid,
			[]string{
				"genres", "tags", "aliases", 
				"release-groups", "artist-rels", "url-rels",
			}...)
	})

	if err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Errored finding MusicBrainz Artist given " +
			"MBID: %v", mbid)
		return nil
	}
	return resp.(*gomusicbrainz.Artist)
}