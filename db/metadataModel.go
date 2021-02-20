package db

import (
	"encoding/json"
	"io/ioutil"
)

// Metadata represents a metadata file for an Artist or an Album
type Metadata struct {
	ID 					 int    `json:"id"`
	FolderID 		 int    `db:"folder_id" json:"folderId"`
	FileSize     int64  `db:"file_size"`
	LastModified int64  `db:"last_modified"`
	Path 				 string `db:"path" json:"path"`
}

// Save creates a new metadata object in the database
func (m *Metadata) Save() error {
	return DB.SaveMetadata(m)
}

// Load pulls existing metadata from the db
func (m *Metadata) Load() (Metadata, error) {
	return DB.LoadMetadata(m)
}

// Delete removes existing metadata from the db
func (m *Metadata) Delete() error {
	return DB.DeleteMetadata(m)
}

// ReadMetadata returns a map with the data found in the metadata file
func (m *Metadata) ReadMetadata() (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return nil, err
	}

	var j map[string]interface{}
	if err := json.Unmarshal(b, &j); err != nil {
		return nil, err
	}

	return j, nil
} 