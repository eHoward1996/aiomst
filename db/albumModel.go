package db

// Album represents an album and contains information extracted from song tags
type Album struct {
	ID       	int    	`json:"id"`
	ArtID 		int 		`db:"art_id" json:"artId"`
	Artist   	string 	`json:"artist"`
	ArtistID 	int    	`db:"artist_id" json:"artistId"`
	FolderID  int     `db:"folder_id" json:"folderId"`
	Title    	string 	`db:"title" json:"title"`
	Year     	int    	`db:"year" json:"year"`
}

// GetAlbumFromSong creates a new Album from a Song model, extracting its
// fields as needed to build the struct
func GetAlbumFromSong(song *Song) *Album {
	return &Album{
		Artist: song.Artist,
		Title:  song.Album,
		Year:   song.Year,
	}
}

// GetArtID returns the ArtID from the struct
func (a *Album) GetArtID() int	{
	return a.ArtID
}

// SetArtID set the ArtID for the struct
func (a *Album) SetArtID(aID int) error {
	a.ArtID = aID
	return DB.UpdateAlbumArt(a)
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
