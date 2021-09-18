package db

import (
	"database/sql"
)

// albumQuery loads a slice of Album structs matching the input query
func (s *SqlBackend) albumQuery(query string, args ...interface{}) ([]Album, error) {
	// Perform input query with arguments
	rows, err := s.db.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// Iterate all rows
	albums := make([]Album, 0)
	a := Album{}
	for rows.Next() {
		// Scan album into struct
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}

		// Append to list
		albums = append(albums, a)
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

// AllAlbums loads a slice of all Album structs from the database
func (s *SqlBackend) AllAlbums() ([]Album, error) {
	return s.albumQuery(
		`SELECT 
			albums.*,
			artists.title AS artist 
		FROM albums
		JOIN artists ON albums.artist_id = artists.id;`,
	)
}

// AllAlbumsByTitle loads a slice of all Album structs from the database ordered
// by their title case insensitive
func (s *SqlBackend) AllAlbumsByTitle() ([]Album, error) {
	return s.albumQuery(
		`SELECT 
			albums.*,
			artists.title AS artist 
		FROM albums
		JOIN artists ON albums.artist_id = artists.id
		ORDER BY albums.title
		COLLATE NOCASE ASC;`,
	)
}

// LimitAlbums loads a slice of Album structs from the database using SQL limit, where the first parameter
// specifies an offset and the second specifies an item count
func (s *SqlBackend) LimitAlbums(offset int, count int) ([]Album, error) {
	return s.albumQuery(
		`SELECT 
			albums.*,
			artists.title AS artist 
		FROM albums
		JOIN artists ON albums.artist_id = artists.id 
		LIMIT ?, ?;`, 
		offset, 
		count,
	)
}

// AlbumsForArtist loads a slice of all Album structs with matching artist ID
func (s *SqlBackend) AlbumsForArtist(ID int) ([]Album, error) {
	return s.albumQuery(
		`SELECT 
			albums.*,
			artists.title AS artist 
		FROM albums
		JOIN artists ON albums.artist_id = artists.id 
		WHERE albums.artist_id = ?;`,
		ID,
	)
}

// SearchAlbums loads a slice of all Album structs from the database which contain
// titles that match the specified search query
func (s *SqlBackend) SearchAlbums(query string) ([]Album, error) {
	return s.albumQuery(
		`SELECT 
			albums.*,
			artists.title AS artist 
		FROM albums
		JOIN artists ON albums.artist_id = artists.id 
		WHERE albums.normalized_title LIKE ?;`, 
		"%"+query+"%",
	)
}

// CountAlbums fetches the total number of Album structs from the database
func (s *SqlBackend) CountAlbums() (int64, error) {
	return s.integerQuery("SELECT COUNT(*) AS int FROM albums;")
}

// DeleteAlbum removes a Album from the database
func (s *SqlBackend) DeleteAlbum(a *Album) error {
	// Attempt to delete this album by its ID, if available
	tx := s.db.MustBegin()
	if a.ID != 0 {
		tx.Exec("DELETE FROM albums WHERE id = ?;", a.ID)
		return tx.Commit()
	}

	// Else, attempt to remove the album by its artist ID and title
	tx.Exec("DELETE FROM albums WHERE artist_id = ? AND title = ?;",
		a.ArtistID,
		a.Title,
	)
	return tx.Commit()
}

// LoadAlbum loads an Album from the database, populating the parameter struct
func (s *SqlBackend) LoadAlbum(a *Album) error {
	// Load the album via ID if available
	if a.ID != 0 {
		if err := s.db.Get(
			a,
			`SELECT 
				albums.*, 
				artists.title AS artist 
			FROM albums 
			JOIN artists ON albums.artist_id = artists.id 
			WHERE albums.id = ?;`,
			a.ID,
		);
		err != nil {
			return err
		}
		return nil
	}

	// Load via artist ID and album title
	if err := s.db.Get(
		a,
		`SELECT 
			albums.*, 
			artists.title AS artist 
		FROM albums
		JOIN artists ON albums.artist_id = artists.id 
		WHERE 
		albums.artist_id = ? AND albums.title = ?;`,
		a.ArtistID, a.Title,
	);
	err != nil {
		return err
	}
	return nil
}

// SaveAlbum attempts to save an Album to the database
func (s *SqlBackend) SaveAlbum(a *Album) error {
	// Insert new album
	query := `INSERT INTO albums
		(
			art_id, mb_id, discogs_id,
			metadata_id, artist_id, folder_id,
			title, normalized_title, year
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`
	tx := s.db.MustBegin()
	tx.MustExec(query, 
		a.ArtID, a.MBID, a.DiscogsID, 
		a.MetadataID, a.ArtistID,	a.FolderID, 
		a.Title, a.NormalizedTitle, a.Year);
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// If no ID, reload to grab it
	if a.ID == 0 {
		if err := s.LoadAlbum(a); err != nil {
			return err
		}
	}
	return nil
}

// UpdateAlbum updates the row in the database based where the id is the same 
// as the passed album objects ID
func (s *SqlBackend) UpdateAlbum(a *Album) error {
	query := `UPDATE albums 
		SET
			mb_id = ?,
			discogs_id = ?,
			metadata_id = ?,
			art_id = ?,
			artist_id = ?,
			folder_id = ?,
			title = ?,
			normalized_title = ?,
			year = ?
		WHERE id = ?;`
	tx := s.db.MustBegin()
	tx.Exec(
		query,
		a.MBID, 
		a.DiscogsID,
		a.MetadataID, 
		a.ArtID, 
		a.ArtistID, 
		a.FolderID, 
		a.Title, 
		a.NormalizedTitle, 
		a.Year,
		a.ID,
	)
	return tx.Commit()
}

// PurgeOrphanAlbums deletes all albums who are "orphaned", meaning that they no
// longer have any songs which reference their ID
func (s *SqlBackend) PurgeOrphanAlbums() (int, error) {
	// Select all albums without a song referencing their album ID
	rows, err := s.db.Queryx(
		`SELECT 
			albums.id 
		FROM albums 
		LEFT JOIN songs ON albums.id = songs.album_id 
		WHERE songs.album_id IS NULL;`,
	)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}
	defer rows.Close()

	// Open a transaction to remove all orphaned albums
	tx := s.db.MustBegin()

	// Iterate all rows
	album := new(Album)
	total := 0
	for rows.Next() {
		// Scan ID into struct
		if err := rows.StructScan(album); err != nil {
			return -1, err
		}

		// Remove album
		tx.Exec("DELETE FROM albums WHERE id = ?;", album.ID)
		total++
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return -1, err
	}

	return total, tx.Commit()
}

// AlbumsWithErroredThirdPartyId returns a list of Albums where any id used by 
// a third party is "errored" or non-unique.
// Currently, the only Third Party APIs being used are
//     1. MusicBrainz (MusicBrainzID or MBID)
//     2. Discogs (DiscogsID)
func (s *SqlBackend) AlbumsWithErroredThirdPartyId() ([]Album, error) {
	return s.albumQuery(
		`SELECT
			albums.*, 
			artists.title AS artist
		FROM albums 
		JOIN artists 
		ON artists.id = albums.artist_id 
		WHERE albums.mb_id = "errored" 
		OR albums.mb_id IN (
			SELECT albums.mb_id 
			FROM albums 
			GROUP BY albums.mb_id 
			HAVING COUNT (*) > 1
		)
		OR albums.discogs_id = "errored"
		OR albums.discogs_id IN (
			SELECT albums.discogs_id
			FROM albums
			GROUP BY albums.discogs_id
			HAVING COUNT (*) > 1
		)
		ORDER BY artist;`)
}