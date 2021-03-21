package db

import (
	"encoding/json"
	"fmt"
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
func (m *Metadata) Load() error {
	return DB.LoadMetadata(m)
}

// Delete removes existing metadata from the db
func (m *Metadata) Delete() error {
	return DB.DeleteMetadata(m)
}

// ReadMetadata returns a map with the data found in the metadata file
func (m *Metadata) ReadMetadata() ([]byte, error) {
	return ioutil.ReadFile(m.Path)	
} 


func (m *Metadata) ToStruct(t interface{}) error {
	mdBytes, e := m.ReadMetadata()
	if e != nil {
		return e
	}

	switch t.(type) {
	case *ArtistMetadata:
		return json.Unmarshal(mdBytes, t)
	case *AlbumMetadata:
		return json.Unmarshal(mdBytes, t)
	default:
		return fmt.Errorf("Unknown Object: %T", t)
	}	
}