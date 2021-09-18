package db

import (
	"database/sql"
)

// metadataQuery loads a slice of Metadata structs matching the input query
func (s *SqlBackend) metadataQuery(query string, args ...interface{}) ([]Metadata, error) {
	// Perform input query with arguments
	rows, err := s.db.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// Iterate all rows
	meta := make([]Metadata, 0)
	m := Metadata{}
	for rows.Next() {
		// Scan album into struct
		if err := rows.StructScan(&m); err != nil {
			return nil, err
		}

		// Append to list
		meta = append(meta, m)
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return meta, nil
}

// LoadMetadata loads Metadata from the database, populating the parameter struct
func (s *SqlBackend) LoadMetadata(m *Metadata) error {
	// Load metadata via ID if available
	if m.ID != 0 {
		if err := s.db.Get(m, "SELECT * FROM metadata WHERE id = ?;", m.ID); 
		err != nil {
			return err	
		}
		return nil
	}

	// Load via Path
	if err := s.db.Get(m, "SELECT * FROM metadata WHERE path = ?;", m.Path);
	err != nil {
		return err
	}
	return nil
}


// SaveMetadata attempts to save an Metadata to the database
func (s *SqlBackend) SaveMetadata(m *Metadata) error {
	// Insert new album
	query := "INSERT INTO metadata " +
		"(folder_id, file_size, last_modified, path) VALUES (?, ?, ?, ?);"
	tx := s.db.MustBegin()
	tx.MustExec(query, m.FolderID, m.FileSize, m.LastModified, m.Path);
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// If no ID, reload to grab it
	if m.ID == 0 {
		if err := s.LoadMetadata(m); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMetadata updates the Metadata object
func (s *SqlBackend) UpdateMetadata(m *Metadata) error {
	query := "UPDATE metadata SET folder_id = ?, path = ? WHERE id = ?;"
	tx := s.db.MustBegin()
	tx.Exec(query, m.FolderID, m.Path, m.ID)
	return tx.Commit()
}

// DeleteMetadata removes Metadata from the database
func (s *SqlBackend) DeleteMetadata(m *Metadata) error {
	// Attempt to delete metadata by its ID, if available
	tx := s.db.MustBegin()
	if m.ID != 0 {
		tx.Exec("DELETE FROM metadata WHERE id = ?;", m.ID)
		return tx.Commit()
	}

	// Else, attempt to remove metadata by its path
	tx.Exec("DELETE FROM metadata WHERE path = ?;", m.Path)
	return tx.Commit()
}