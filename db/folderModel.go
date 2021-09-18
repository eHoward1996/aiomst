package db

// Folder represents a filesystem folder
type Folder struct	{
	ID       int    `json:"id"`
	ParentID int    `db:"parent_id" json:"parentId"`
	Title    string `json:"title"`
	Path		 string `json:"path"`
}

// SubFolders retrieves all folder with this folder as the parent
func (f *Folder) SubFolders()	([]Folder, error)	{
	return DB.SubFolders(f.ID)
}

// Delete removes an existing folder from the DB
func (f *Folder) Delete()	error {
	return DB.DeleteFolder(f)
}

// Load pulls an existing folder from the db
func (f *Folder) Load() error {
	return DB.LoadFolder(f)
}

// Save creates a new folder in the db
func (f *Folder) Save() error {
	return DB.SaveFolder(f)
}