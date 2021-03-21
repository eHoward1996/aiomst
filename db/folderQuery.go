package db

import (
	"database/sql"
)

// folderQuery loads a slice of Folder structs matching the input query
func (s *SqlBackend) folderQuery(query string, args ...interface{}) ([]Folder, error) {
	// Perform input query with arguments
	rows, err := s.db.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	// Iterate all rows
	folders := make([]Folder, 0)
	a := Folder{}
	for rows.Next() {
		// Scan folder into struct
		if err := rows.StructScan(&a); err != nil {
			return nil, err
		}

		// Append to list
		folders = append(folders, a)
	}

	// Error check rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

// AllFolders loads a slice of all Folder structs from the database
func (s *SqlBackend) AllFolders() ([]Folder, error) {
	return s.folderQuery("SELECT * FROM folders;")
}

// LimitFolders loads a slice of Folder structs from the database using SQL limit, where the first parameter
// specifies an offset and the second specifies an item count
func (s *SqlBackend) LimitFolders(offset int, count int) ([]Folder, error) {
	return s.folderQuery("SELECT * FROM folders LIMIT ?, ?;", offset, count)
}

// SubFolders loads a slice of all Folder structs residing directly beneath this one from the database
func (s *SqlBackend) SubFolders(parentID int) ([]Folder, error) {
	return s.folderQuery("SELECT * FROM folders WHERE parent_id = ?;", parentID)
}

// FoldersInPath loads a slice of all Folder structs contained within the specified file path
func (s *SqlBackend) FoldersInPath(path string) ([]Folder, error) {
	return s.folderQuery("SELECT * FROM folders WHERE path LIKE ?;", path+"%")
}

// FoldersNotInPath loads a slice of all Folder structs NOT contained within the specified file path
func (s *SqlBackend) FoldersNotInPath(path string) ([]Folder, error) {
	return s.folderQuery("SELECT * FROM folders WHERE path NOT LIKE ?;", path+"%")
}

// SearchFolders loads a slice of all Folder structs from the database which contain
// titles that match the specified search query
func (s *SqlBackend) SearchFolders(query string) ([]Folder, error) {
	return s.folderQuery(
		"SELECT * FROM folders WHERE title LIKE ?;", "%"+query+"%")
}

// CountFolders fetches the total number of Folder structs from the database
func (s *SqlBackend) CountFolders() (int64, error) {
	return s.integerQuery("SELECT COUNT(*) AS int FROM folders;")
}

// DeleteFolder removes a Folder from the database
func (s *SqlBackend) DeleteFolder(f *Folder) error {
	// Attempt to delete this folder by its ID, if available
	tx := s.db.MustBegin()
	if f.ID != 0 {
		tx.Exec("DELETE FROM folders WHERE id = ?;", f.ID)
		return tx.Commit()
	}

	// Else, attempt to remove the folder by its path
	tx.Exec("DELETE FROM folders WHERE path = ?;", f.Path)
	return tx.Commit()
}

// LoadFolder loads a Folder from the database, populating the parameter struct
func (s *SqlBackend) LoadFolder(f *Folder) error {
	// Load the folder via ID if available
	if f.ID != 0 {
		if err := s.db.Get(f, "SELECT * FROM folders WHERE id = ?;", f.ID);
		err != nil {
			return err
		}
		return nil
	}

	// Load via path
	if err := s.db.Get(f, "SELECT * FROM folders WHERE path = ?;", f.Path);
	err != nil {
		return err
	}
	return nil
}

// SaveFolder attempts to save an Folder to the database
func (s *SqlBackend) SaveFolder(f *Folder) error {
	// Insert new folder
	query := "INSERT INTO folders (parent_id, title, path) VALUES (?, ?, ?);"
	tx := s.db.MustBegin()
	tx.Exec(query, f.ParentID, f.Title, f.Path)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// If no ID, reload to grab it
	if f.ID == 0 {
		if err := s.LoadFolder(f); err != nil {
			return err
		}
	}
	return nil
}