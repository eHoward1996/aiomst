package db

import (
	"fmt"

	"github.com/eHoward1996/aiomst/util"
)

// Album represents an album and contains information extracted from song tags
type Album struct {
	ID       	      int    	`json:"id"`
	MBID						string  `db:"mb_id" json:"mbId"`
	DiscogsID       int     `db:"discogs_id" json:"discogsId"`
	MetadataID      int     `db:"metadata_id" json:"metadataId"`
	ArtID 		      int 		`db:"art_id" json:"artId"`
	Artist   	      string 	`json:"artist"`
	ArtistID 	      int    	`db:"artist_id" json:"artistId"`
	FolderID        int     `db:"folder_id" json:"folderId"`
	Title    	      string 	`db:"title" json:"title"`
	NormalizedTitle string  `db:"normalized_title" json:"normalizedTitle"`
	Year     	      int    	`db:"year" json:"year"`
}

// GetAlbumFromSong creates a new Album from a Song model, extracting its
// fields as needed to build the struct
func GetAlbumFromSong(song *Song) *Album {
	return &Album{
		Artist:          song.Artist,
		Title:           song.Album,
		NormalizedTitle: util.NormalizeString(song.Album),
		Year:            song.Year,
	}
}

// GetArt returns an Art object based on the ArtId from this struct
func (a *Album) GetArt() (*Art, error)	{
	art := new(Art)
	art.ID = a.ArtID
	if err := art.Load(); err != nil {
		return nil, err
	}
	return art, nil
}

// GetMetadataObj returns a Metadata object based on the MetadataID from this
// struct
func (a *Album) GetMetadataObj() (*Metadata, error) {
	md := new(Metadata)
	md.ID = a.MetadataID
	if err := md.Load(); err != nil {
		return nil, err
	}
	return md, nil
}

// Delete removes an existing Album from the database
func (a *Album) Delete() error {
	return DB.DeleteAlbum(a)
}

// Load pulls an existing Album from the database
func (a *Album) Load() error {
	return DB.LoadAlbum(a)
}

// Save creates a new Album in the database
func (a *Album) Save() error {
	return DB.SaveAlbum(a)
}

// Update saves an existing Album in the database
func (a *Album) Update() error {
	return DB.UpdateAlbum(a)
}

// HasAttachables denotes this object has "attachable" or associated objects.
// Objects that can be attached have an "IsAttachable" method.
// TODO: Find a better way to make loose object associations
func (a *Album) HasAttachables() {}

// String is a method that returns a string with simple information about this
// object. 
func (a Album) String() string {
	return fmt.Sprintf("%s - %s", a.Title, a.Artist)	
}