package db

import (
	"fmt"
)

// Artist represents an artist in the db and contains a unique ID and name
type Artist struct {
	ID    		      int    	`json:"id"`
	MBID						string  `db:"mb_id" json:"mbId"`
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
		NormalizedTitle: normalizeString(song.Artist),
	}
}

// GetArt returns an Art object based on the ArtID from this struct
func (a *Artist) GetArt() (*Art, error) {
	art := new(Art)
	art.ID = a.ArtID
	artObj, err := art.Load()
	if err != nil {
		return nil, err
	}
	return &artObj, nil
}

// SetArtID sets the ArtID for the struct
func (a *Artist) SetArtID(aID int) error {
	a.ArtID = aID
	return DB.UpdateArtistArt(a)
}

// GetMetadata returns a Metadata object based on the MetadataID from this struct
func (a *Artist) GetMetadata() (*Metadata, error) {
	md := new(Metadata)
	md.ID = a.MetadataID
	mdObj, err := md.Load()
	if err != nil {
		return nil, err
	}
	return &mdObj, nil
}

// SetMetadataID sets the MetadataID for the struct
func (a *Artist) SetMetadataID(mID int) error {
	a.MetadataID = mID
	return DB.UpdateArtistMetadata(a)
}

// Delete removes an existing Artist from the database
func (a *Artist) Delete() error {
	return DB.DeleteArtist(a)
}

// Load pulls an existing Artist from the database
func (a *Artist) Load() (Artist, error) {
	return DB.LoadArtist(a)
}

// Save creates a new Artist in the database
func (a *Artist) Save() error {
	return DB.SaveArtist(a)
}

// ToString is a method that returns a string with simple information about this
// object.
func (a Artist) ToString() string {
	return fmt.Sprintf("%s", a.Title)	
}