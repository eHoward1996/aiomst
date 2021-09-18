package db

import (
	"github.com/eHoward1996/aiomst/util"
)

// Artist represents an artist in the db and contains a unique ID and name
type Artist struct {
	ID    		      int    	`json:"id"`
	MBID						string  `db:"mb_id" json:"mbId"`
	DiscogsID       string     `db:"discogs_id" json:"discogsId"`
	MetadataID      int     `db:"metadata_id" json:"metadataId"`
	ArtID  		      int 		`db:"art_id" json:"artId"` 
	FolderID 	      int 		`db:"folder_id" json:"folderId"`
	Title 		      string 	`db:"title" json:"title"`
	NormalizedTitle string  `db:"normalized_title" json:"normalizedTitle"`
}

// GetArtistFromSong creates a new Artist from a Song model, extracting its
// fields as needed to build the struct
func GetArtistFromSong(song *Song) *Artist {
	// Copy the artist name to 
	return &Artist{
		Title:           song.Artist,
		NormalizedTitle: util.NormalizeString(song.Artist),
	}
}

// GetArt returns an Art object based on the ArtID from this struct
func (a *Artist) GetArt() (*Art, error) {
	art := &Art{ID: a.ArtID}
	if err := art.Load(); err != nil {
		return nil, err
	}
	return art, nil
}

// GetMetadataObj returns a Metadata object based on the MetadataID from this
// struct
func (a *Artist) GetMetadataObj() (*Metadata, error) {
	md := &Metadata{ID: a.MetadataID}
	if err := md.Load(); err != nil {
		return nil, err
	}
	return md, nil
}


// Delete removes an existing Artist from the database
func (a *Artist) Delete() error {
	return DB.DeleteArtist(a)
}

// Load pulls an existing Artist from the database
func (a *Artist) Load() error {
	return DB.LoadArtist(a)
}

// Save creates a new Artist in the database
func (a *Artist) Save() error {
	return DB.SaveArtist(a)
}

// Update saves an existing Artist in the database
func (a *Artist) Update() error {
	return DB.UpdateArtist(a)
}

// HasAttachables denotes this object has "attachable" or associated objects.
// Objects that can be attached have an "IsAttachable" method.
// TODO: Find a better way to make loose object associations
func (a *Artist) HasAttachables() {}

// String is a method that returns a string with simple information about this
// object.
func (a Artist) String() string {
	return a.Title
}