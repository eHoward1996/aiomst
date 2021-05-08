package integrated

import (
	"fmt"

	"github.com/eHoward1996/aiomst/db"
	"github.com/irlndts/go-discogs"
)


func (i TPIntegrator) getDiscogsData(obj interface{}) (discogsResponse, error) {
	var resp discogsResponse = discogsResponse{}
	var err error = nil

	if i.discogsClient == nil {
		return resp, fmt.Errorf("No Discogs Client available")
	}

	client := *i.discogsClient
	switch t := obj.(type) {
	case *sendAlbum:
		resp.AlbumResponse, err = i.requestDiscogsMaster(client, t)
	}

	return resp, err
}

func (i TPIntegrator) requestDiscogsMaster(
	client discogs.Discogs, a *sendAlbum) (*discogs.Master, error) {
		// Remove any disambiguation/special edition info from the album title
		// Discogs really doesn't like it.
		albumTitle := stripAllParentheses(a.Title)
		request := discogs.SearchRequest{
			Artist: a.Artist, 
			ReleaseTitle: albumTitle, 
			Type: "master",
			PerPage: 1,
		}

		search, err := client.Search(request)
		if err != nil {
			e := fmt.Errorf("Error making Discogs Search Request: %v", err)
			return nil, e
		}
		if search != nil && len(search.Results) == 1 {
			result := search.Results[0]
			return client.Master(result.MasterID)
		}
		
		return nil, nil
}

func buildDiscogsAlbumMetadata(master *discogs.Master) db.DiscogsMetadata {
	if master == nil {
		return db.DiscogsMetadata{}
	}
	
	var dmd db.DiscogsMetadata = db.DiscogsMetadata{
		Styles: master.Styles,
		Genres: master.Genres,
		Title:  master.Title,
		Year:   master.Year,
		URI:    master.URI,
	}

	var tracklist []db.DiscogsTrack = make([]db.DiscogsTrack, 0)
	for _, dTrack := range master.Tracklist {
		track := db.DiscogsTrack {
			Duration: dTrack.Duration,
			Position: dTrack.Position,
			Title:    dTrack.Title,
			Type:     dTrack.Type,
		}

		extraArtists := make([]db.DiscogsArtist, 0)
		for _, ea := range dTrack.Extraartists {
			a := db.DiscogsArtist{
				ID: ea.ID,
				Name: ea.Name,
				ResourceURL: ea.ResourceURL,
				Role: ea.Role,
				Tracks: ea.Tracks,
			}
			extraArtists = append(extraArtists, a)
		}
		track.ExtraArtists = extraArtists

		artists := make([]db.DiscogsArtist, 0)
		for _, artist := range dTrack.Artists {
			a := db.DiscogsArtist {
				ID:          artist.ID,
				Name:        artist.Name,
				ResourceURL: artist.ResourceURL,
				Role:        artist.Role,
				Tracks:      artist.Tracks,
			}
			artists = append(artists, a)
		}
		track.Artists = artists
	}
	dmd.Tracklist = tracklist

	artists := make([]db.DiscogsArtist, 0)
	for _, disArtist := range master.Artists {
		a := db.DiscogsArtist {
			ID:          disArtist.ID,
			Name:        disArtist.Name,
			ResourceURL: disArtist.ResourceURL,
			Role:        disArtist.Role,
			Tracks:      disArtist.Tracks,
		}
		artists = append(artists, a)
	}
	dmd.Artists = artists

	images := make([]db.DiscogsImage, 0)
	for _, img := range master.Images {
		i := db.DiscogsImage {
			Height:      img.Height,
			Width:       img.Width,
			ResourceURL: img.ResourceURL,
			Type:        img.Type,
			URI:         img.URI,
		}
		images = append(images, i)
	}
	dmd.Images = images

	videos := make([]db.DiscogsVideo, 0)
	for _, vid := range master.Videos {
		v := db.DiscogsVideo {
			Description: vid.Description,
			Duration:    vid.Duration,
			Title:       vid.Title,
			URI:         vid.URI,
		}
		videos = append(videos, v)
	}
	dmd.Videos = videos
	return dmd
}