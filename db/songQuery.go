package db

import (
	"database/sql"
)

// songQuery loads a slice of Song structs matching the input query
func (s *SqlBackend) songQuery(query string, args ...interface{}) ([]Song, error) {
	// Perform input query with arguments
	rows, err := s.db.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// Iterate all rows
	songs := make([]Song, 0)
	a := Song{}
	for rows.Next() {
		// Scan song into struct
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}

		// Append to list
		songs = append(songs, a)
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

// AllSongs loads a slice of all Song structs from the database
func (s *SqlBackend) AllSongs() ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*, 
			artists.title AS artist, 
			albums.title AS album
		FROM songs 
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id;`,
	)
}

// AllSongsByTitle loads a slice of all Song structs from the database by Song
// title case insensitive
func (s *SqlBackend) AllSongsByTitle() ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album
		FROM songs
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id
		ORDER BY title COLLATE NOCASE ASC;`,
	)
}

// LimitSongs loads a slice of Song structs from the database using SQL limit, where the first parameter
// specifies an offset and the second specifies an item count
func (s *SqlBackend) LimitSongs(offset int, count int) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album 
		FROM songs
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id 
		LIMIT ?, ?;`,
		offset,
		count,
	)
}

// RandomSongs loads a slice of 'n' random song structs from the database
func (s *SqlBackend) RandomSongs(n int) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album 
		FROM songs 
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id
		ORDER BY RANDOM() LIMIT ?;`,
		n,
	)
}

// SearchSongs loads a slice of all Song structs from the database which contain
// titles that match the specified search query
func (s *SqlBackend) SearchSongs(query string) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album 
		FROM songs
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id
		WHERE songs.normalized_title LIKE ?;`,
		"%"+query+"%",
	)
}

// SongsForAlbum loads a slice of all Song structs which have the matching album ID
func (s *SqlBackend) SongsForAlbum(ID int) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album 
		FROM songs
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id
		WHERE songs.album_id = ?;`,
		ID,
	)
}

// SongsForArtist loads a slice of all Song structs which have the matching artist ID
func (s *SqlBackend) SongsForArtist(ID int) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album 
		FROM songs
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id
		WHERE songs.artist_id = ?;`, 
		ID,
	)
}

// SongsForFolder loads a slice of all Song structs which have the matching folder ID
func (s *SqlBackend) SongsForFolder(ID int) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album 
		FROM songs
		JOIN artists ON songs.artist_id = artists.id
		JOIN albums ON songs.album_id = albums.id
		WHERE songs.folder_id = ?;`,
		ID,
	)
}

// SongsInPath loads a slice of all Song structs residing under the specified
// filesystem path from the database
func (s *SqlBackend) SongsInPath(path string) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*,
			artists.title AS artist,
			albums.title AS album
		FROM songs 
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id 
		WHERE songs.path LIKE ?;`, 
		path+"%",
	)
}

// SongsNotInPath loads a slice of all Song structs that do not reside under the specified
// filesystem path from the database
func (s *SqlBackend) SongsNotInPath(path string) ([]Song, error) {
	return s.songQuery(
		`SELECT 
			songs.*, 
			artists.title AS artist,
			albums.title AS album
		FROM songs 
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id 
		WHERE songs.path NOT LIKE ?;`,
		path+"%",
	)
}

// CountSongs fetches the total number of Artist structs from the database
func (s *SqlBackend) CountSongs() (int64, error) {
	return s.integerQuery("SELECT COUNT(*) AS int FROM songs;")
}

// DeleteSong removes a Song from the database
func (s *SqlBackend) DeleteSong(a *Song) error {
	// Attempt to delete this song by its ID, if available
	tx := s.db.MustBegin()
	if a.ID != 0 {
		tx.Exec("DELETE FROM songs WHERE id = ?;", a.ID)
		return tx.Commit()
	}

	// Else, attempt to remove the song by its file name
	tx.Exec("DELETE FROM songs WHERE path = ?;", a.Path)
	return tx.Commit()
}

// LoadSong loads a Song from the database, populating the parameter struct
func (s *SqlBackend) LoadSong(a *Song) error {
	// Load the song via ID if available
	if a.ID != 0 {
		if err := s.db.Get(
			a, 
			`SELECT 
				songs.*, 
				artists.title AS artist, 
				albums.title AS album 
			FROM songs 
			JOIN artists ON songs.artist_id = artists.id  
			JOIN albums ON songs.album_id = albums.id 
			WHERE songs.id = ?;`,
			a.ID,
		);
		err != nil {
			return err
		}
		return nil
	}

	// Load via file name
	if err := s.db.Get(
		a,
		`SELECT 
			songs.*, 
			artists.title AS artist, 
			albums.title AS album 
		FROM songs 
		JOIN artists ON songs.artist_id = artists.id 
		JOIN albums ON songs.album_id = albums.id 
		WHERE songs.path = ?;`, 
		a.Path,
	); err != nil {
		return err
	}
	return nil
}

// SaveSong attempts to save a Song to the database
func (s *SqlBackend) SaveSong(a *Song) error {
	// Insert new song
	query := `INSERT INTO songs (
			mb_id, album_id, artist_id,
			bitrate, channels, comment, 
			path, file_size, file_type_id, 
			folder_id, genre, last_modified, 
			length, sample_rate, title,	
			normalized_title, track, year
		)  
		VALUES (
			?, ?, ?, 
			?, ?, ?,
			?, ?, ?, 
			?, ?, ?,
			?, ?, ?, 
			?, ?, ?
		);`
	tx := s.db.MustBegin()
	tx.MustExec(
		query, 
		a.MBID, a.AlbumID, a.ArtistID,
		a.Bitrate, a.Channels, a.Comment,
		a.Path, a.FileSize, a.FileTypeID,
		a.FolderID, a.Genre, a.LastModified,
		a.Length, a.SampleRate, a.Title,
		a.NormalizedTitle, a.Track, a.Year,
	)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// If no ID, reload to grab it
	if a.ID == 0 {
		if err := s.LoadSong(a); err != nil {
			return err
		}
	}
	return nil
}

// UpdateSong attempts to update a Song in the database
func (s *SqlBackend) UpdateSong(a *Song) error {
	// Update existing song
	query := `UPDATE songs 
		SET 
			mb_id = ?, album_id = ?, artist_id = ?, 
			bitrate = ?, channels = ?, comment = ?,
			file_size = ?, folder_id = ?, genre = ?,
			last_modified = ?, length = ?, sample_rate = ?, 
			title = ?, track = ?, year = ? 
		WHERE id = ?;`
	tx := s.db.MustBegin()
	tx.Exec(
		query, 
		a.MBID, a.AlbumID, a.ArtistID,
		a.Bitrate, a.Channels, a.Comment, 
		a.FileSize, a.FolderID, a.Genre, 
		a.LastModified, a.Length, a.SampleRate, 
		a.Title, a.Track, a.Year,
		a.ID,
	)

	return tx.Commit()
}