package db

// Artist represents an artist in the db and contains a unique ID and name
type Artist struct {
	ID    		      int    	`json:"id"`
	ArtID  		      int 		`db:"art_id" json:"artId"` 
	FolderID 	      int 		`db:"folder_id" json:"folderId"`
	Title 		      string 	`db:"title" json:"title"`
	NormalizedTitle string  `db:"normalized_title" json:"normalizedTitle"`
}

// GetArtistFromSong creates a new Artist from a Song model, extracting its
// fields as needed to build the struct
func GetArtistFromSong(song *Song) *Artist {
	// Copy the artist name to title
	return &Artist{
		Title: song.Artist,
	}
}

// GetArtID returns the ArtID of this struct
func (a *Artist) GetArtID() int {
	return a.ArtID
}

// SetArtID sets the ArtID for the struct
func (a *Artist) SetArtID(aID int) error {
	a.ArtID = aID
	return DB.UpdateArtistArt(a)
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
