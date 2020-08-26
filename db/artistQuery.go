package db

import (
	"database/sql"
)

// artistQuery loads a slice of Artist structs matching the input query
func (s *SqlBackend) artistQuery(query string, args ...interface{}) ([]Artist, error) {
	// Perform input query with arguments
	rows, err := s.db.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// Iterate all rows
	artists := make([]Artist, 0)
	a := Artist{}
	for rows.Next() {
		// Scan artist into struct
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}

		// Append to list
		artists = append(artists, a)
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

// AllArtists loads a slice of all Artist structs from the database
func (s *SqlBackend) AllArtists() ([]Artist, error) {
	return s.artistQuery("SELECT * FROM artists;")
}

// AllArtistsByTitle loads a slice of all Artist structs from the database,
// sorted alphabetically by title case insensitive
func (s *SqlBackend) AllArtistsByTitle() ([]Artist, error) {
	return s.artistQuery("SELECT * FROM artists ORDER BY title COLLATE NOCASE ASC;")
}

// LimitArtists loads a slice of Artist structs from the database using SQL limit, where the first parameter
// specifies an offset and the second specifies an item count
func (s *SqlBackend) LimitArtists(offset int, count int) ([]Artist, error) {
	return s.artistQuery("SELECT * FROM artists LIMIT ?, ?;", offset, count)
}

// SearchArtists loads a slice of all Artist structs from the database which contain
// titles that match the specified search query
func (s *SqlBackend) SearchArtists(query string) ([]Artist, error) {
	return s.artistQuery("SELECT * FROM artists WHERE title LIKE ?;", "%"+query+"%")
}

// CountArtists fetches the total number of Artist structs from the database
func (s *SqlBackend) CountArtists() (int64, error) {
	return s.integerQuery("SELECT COUNT(*) AS int FROM artists;")
}

// DeleteArtist removes an Artist from the database
func (s *SqlBackend) DeleteArtist(a *Artist) error {
	// Attempt to delete this artist by its ID, if available
	tx := s.db.MustBegin()
	if a.ID != 0 {
		tx.Exec("DELETE FROM artists WHERE id = ?;", a.ID)
		return tx.Commit()
	}

	// Else, attempt to remove the artist by its title
	tx.Exec("DELETE FROM artists WHERE title = ?;", a.Title)
	return tx.Commit()
}

// LoadArtist loads an Artist from the database, populating the parameter struct
func (s *SqlBackend) LoadArtist(a *Artist) (Artist, error) {
	// Load the artist via ID if available
	r := *a
	if a.ID != 0 {
		if err := s.db.Get(&r, "SELECT * FROM artists WHERE id = ?", a.ID);
		err != nil {
			return Artist{}, err
		}
		return r, nil
	}

	// Load via title
	if err := s.db.Get(&r, "SELECT * FROM artists WHERE title = ?;", a.Title);
	err != nil {
		return Artist{}, err
	}
	return r, nil
}

// SaveArtist attempts to save an Artist to the database
func (s *SqlBackend) SaveArtist(a *Artist) error {
	// Insert new artist
	query := "INSERT INTO artists (art_id, folder_id, title) VALUES (?, ?, ?);"
	tx := s.db.MustBegin()
	tx.Exec(query, a.ArtID, a.FolderID, a.Title)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// If no ID, reload to grab it
	if a.ID == 0 {
		artist, err := s.LoadArtist(a)
		if err != nil {
			return err
		}
		*a = artist
	}
	return nil
}

// UpdateArtistArt updates the Artists artId
func (s *SqlBackend) UpdateArtistArt(a *Artist) error {
	query := "UPDATE artists SET art_id = ? WHERE id = ?;"
	tx := s.db.MustBegin()
	tx.Exec(query, a.ArtID, a.ID)
	return tx.Commit()
}

// PurgeOrphanArtists deletes all artists who are "orphaned", meaning that they no
// longer have any songs which reference their ID
func (s *SqlBackend) PurgeOrphanArtists() (int, error) {
	// Select all artists without a song referencing their artist ID
	rows, err := s.db.Queryx("SELECT artists.id FROM artists LEFT JOIN songs ON " +
		"artists.id = songs.artist_id WHERE songs.artist_id IS NULL;")
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}
	defer rows.Close()

	// Open a transaction to remove all orphaned artists
	tx := s.db.MustBegin()

	// Iterate all rows
	artist := new(Artist)
	total := 0
	for rows.Next() {
		// Scan ID into struct
		if err := rows.StructScan(artist); err != nil {
			return -1, err
		}

		// Remove artist
		tx.Exec("DELETE FROM artists WHERE id = ?;", artist.ID)
		total++
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return -1, err
	}

	return total, tx.Commit()
}