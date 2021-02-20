package db

import (
	"fmt"
)

// Album represents an album and contains information extracted from song tags
type Album struct {
	ID       	      int    	`json:"id"`
	MBID						string  `db:"mb_id" json:"mbId"`
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
		NormalizedTitle: normalizeString(song.Album),
		Year:            song.Year,
	}
}

// GetArt returns an Art object based on the ArtId from this struct
func (a *Album) GetArt() (*Art, error)	{
	art := new(Art)
	art.ID = a.ArtID
	artObj, err := art.Load()
	if err != nil {
		return nil, err
	}
	return &artObj, nil
}

// GetMetadata returns a Metadata object based on the MetadataID from this struct
func (a *Album) GetMetadata() (*Metadata, error) {
	md := new(Metadata)
	md.ID = a.MetadataID
	mdObj, err := md.Load()
	if err != nil {
		return nil, err
	}
	return &mdObj, nil
}

// Delete removes an existing Album from the database
func (a *Album) Delete() error {
	return DB.DeleteAlbum(a)
}

// Load pulls an existing Album from the database
func (a *Album) Load() (Album, error) {
	return DB.LoadAlbum(a)
}

// Save creates a new Album in the database
func (a *Album) Save() error {
	return DB.SaveAlbum(a)
}

// ToString is a method that returns a string with simple information about this
// object. 
func (a Album) ToString() string {
	return fmt.Sprintf("%s - %s", a.Title, a.Artist)	
}