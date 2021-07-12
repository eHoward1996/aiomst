package db

import "time"

type ArtistMetadata struct {
	ArtistName  string              `json:"artist_name"`
	MusicBrainz MusicBrainzMetadata `json:"musicbrainz,omitempty"`
	Discogs 		DiscogsMetadata     `json:"discogs,omitempty"`
}

type AlbumMetadata struct {
	AlbumName    string              `json:"album_name"`
	MusicBrainz  MusicBrainzMetadata `json:"musicbrainz,omitempty"`
	Discogs      DiscogsMetadata     `json:"discogs,omitempty"`
}

type MusicBrainzMetadata struct {
	Artists        []MusicBrainzArtist       `json:"artists,omitempty"`
	AssociatedActs []AssociatedAct           `json:"associated_acts,omitempty"`
	RelatedUrls    []RelatedUrl              `json:"related_urls,omitempty"`
	RelatedTags    []string    						   `json:"related_tags,omitempty"`
	Release        MusicBrainzRelease        `json:"release,omitempty"`
	ReleaseGroups  []MusicBrainzReleaseGroup `json:"release_groups,omitempty"`
}

type MusicBrainzArtist struct {
	ID             string             `json:"id,omitempty"`
	Name           string             `json:"name,omitempty"`
	Disambiguation string             `json:"disambiguation,omitempty"`
	SortName       string             `json:"sort_name,omitempty"`
	Type           string             `json:"type,omitempty"`
	Aliases        []MusicBrainzAlias `json:"aliases,omitempty"`
	Area           string             `json:"area,omitempty"`
	Country        string             `json:"country,omitempty"`
}

type RelatedUrl struct {
	Type string `json:"type,omitempty"`
	Url  string `json:"url,omitempty"`
}

type MusicBrainzRelease struct {
	ID             string             `json:"id,omitempty"`	
	Title          string             `json:"title,omitempty"`
	Disambiguation string             `json:"disambiguation,omitempty"`
	Status         string             `json:"status,omitempty"`
	DiscCount      int                `json:"disc_count,omitempty"` 
	TrackCount     int                `json:"track_count,omitempty"`
	Tracks         []MusicBrainzTrack `json:"tracks,omitempty"`
}

type AssociatedAct struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Disambiguation string `json:"disambiguation,omitempty"`
	Type           string `json:"type,omitempty"`
	Relation       string `json:"relation,omitempty"`
}

type MusicBrainzReleaseGroup struct {
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Type         string    `json:"type,omitempty"`
	ReleaseCount int       `json:"release_count,omitempty"`
	ReleaseDate  time.Time `json:"release_date,omitempty"`
}

type MusicBrainzAlias struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type MusicBrainzTrack struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Number int    `json:"number,omitempty"`
	Length int 		`json:"length,omitempty"`
}

type DiscogsMetadata struct {
	Styles    []string        `json:"styles,omitempty"`
	Genres    []string        `json:"genres,omitempty"`
	Title     string          `json:"title,omitempty"`
	Year      int             `json:"year,omitempty"`
	Tracklist []DiscogsTrack  `json:"tracklist,omitempty"`
	Artists   []DiscogsArtist `json:"artists"`
	Images    []DiscogsImage  `json:"images,omitempty"`
	Videos    []DiscogsVideo  `json:"videos,omitempty"`
	URI       string          `json:"uri,omitempty"`
}

type DiscogsTrack struct {
	Duration     string          `json:"duration"`
	Position     string          `json:"position"`
	Title        string          `json:"title"`
	Type         string          `json:"type"`
	ExtraArtists []DiscogsArtist `json:"extra_artists,omitempty"`
	Artists      []DiscogsArtist `json:"artists,omitempty"`
}

type DiscogsArtist struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Realname    string   `json:"real_name"`
	ResourceURL string   `json:"resource_url"`
	Role        string   `json:"role,omitempty"`
	Tracks      string   `json:"tracks,omitempty"`
	Members     []string `json:"members,omitempty"`
}

type DiscogsImage struct {
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	ResourceURL string `json:"resource_url"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
}

type DiscogsVideo struct {
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Title       string `json:"title"`
	URI         string `json:"uri"`
}