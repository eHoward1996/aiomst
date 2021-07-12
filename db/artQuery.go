package db

import (
	"database/sql"
)

// artQuery loads a slice of Art structs matching the input query
func (s *SqlBackend) artQuery(query string, args ...interface{}) ([]Art, error) {
	// Perform input query with arguments
	rows, err := s.db.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// Iterate all rows
	art := make([]Art, 0)
	a := Art{}
	for rows.Next() {
		// Scan artist into struct
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}

		// Append to list
		art = append(art, a)
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return art, nil
}

// AllArt returns a slice of All Art structs
func (s *SqlBackend) AllArt() ([]Art, error)	{
	return s.artQuery("SELECT * FROM art;")
}

// ArtInPath loads a slice of all Art structs contained within the specified file path
func (s *SqlBackend) ArtInPath(path string) ([]Art, error) {
	return s.artQuery("SELECT * FROM art WHERE path LIKE ?;", path+"%")
}

// ArtNotInPath loads a slice of all Art structs NOT contained within the specified file path
func (s *SqlBackend) ArtNotInPath(path string) ([]Art, error) {
	return s.artQuery("SELECT * FROM art WHERE path NOT LIKE ?;", path+"%")
}

// CountArt fetches the total number of Art structs from the database
func (s *SqlBackend) CountArt() (int64, error) {
	return s.integerQuery("SELECT COUNT(*) AS int FROM art;")
}

// DeleteArt removes Art from the database
func (s *SqlBackend) DeleteArt(a *Art) error {
	// Attempt to delete this art by its ID
	tx := s.db.MustBegin()
	tx.Exec("DELETE FROM art WHERE id = ?;", a.ID)

	// Update any songs using this art ID to have a zero ID
	tx.Exec("UPDATE songs SET art_id = 0 WHERE art_id = ?;", a.ID)
	return tx.Commit()
}

// LoadArt loads Art from the database, populating the parameter struct
func (s *SqlBackend) LoadArt(a *Art) error {
	// Load the artist via ID if available
	if a.ID != 0 {
		if err := s.db.Get(a, "SELECT * FROM art WHERE id = ?;", a.ID);
		err != nil {
			return err
		}
		return nil
	}

	// Load via file name
	if err := s.db.Get(a, "SELECT * FROM art WHERE path = ?;", a.Path);
	err != nil {
		return err
	}
	return nil
}

// SaveArt attempts to save Art to the database
func (s *SqlBackend) SaveArt(a *Art) error {
	// Insert new artist
	query := "INSERT INTO art " +
	"(path, file_size, last_modified) VALUES (?, ?, ?);"
	tx := s.db.MustBegin()
	tx.Exec(query, a.Path, a.FileSize, a.LastModified)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// If no ID, reload to grab it
	if a.ID == 0 {
		if err := s.LoadArt(a); err != nil {
			return err
		}
	}
	return nil
}