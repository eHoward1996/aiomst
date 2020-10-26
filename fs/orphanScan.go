package fs

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/eHoward1996/aiomst/db"
)


// OrphanScan is a filesystem task that scans the given path for orphaned media
type OrphanScan struct	{
	baseFolder 	string
	subFolder 	string
	verbose 		bool
}

// Folders returns the base and sub folders for scanning
func (fs *OrphanScan) Folders()	(string, string)	{
	return fs.baseFolder, fs.subFolder
}

// SetFolders sets the base and sub folders for scanning
func (fs *OrphanScan) SetFolders(baseFolder, subFolder string)	{
	fs.baseFolder = baseFolder
	fs.subFolder  = subFolder
}

// Verbose is whether scanning has verbose output or not
func (fs *OrphanScan) Verbose(v bool)	{
	fs.verbose = v
}

// WhoAmI returns Orphan Scan
func (fs *OrphanScan) WhoAmI() string {
	return "Orphan Scan"
}

// Scan scans for missing "orphaned" media files in the local filesystem
func (fs *OrphanScan) Scan(baseFolder, subFolder string, orphanCancelChan chan struct{}) (int, error) {
	// Halt scan if needed
	var mutex sync.RWMutex
	haltWalk := false
	go func() {
		// Wait for signal
		<- orphanCancelChan

		// Halt fs walk
		mutex.Lock()
		haltWalk = true
		mutex.Unlock()
	}()

	// Track metrics about the scan
	artCount := 0
	folderCount := 0
	songCount := 0
	startTime := time.Now()

	// Check if a baseFolder is set, meaning remove ANYTHING not under this base
	if baseFolder != "" {
		if fs.verbose {
			log.Println("FS: Orphan Scan: Base Folder:", baseFolder)
		}

		// Scan for all art NOT under the base folder
		art, err := db.DB.ArtNotInPath(baseFolder)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		// Remove all art which is not in this path
		for _, a := range art {
			// Remove art from database
			filename := a.FileName
			if err := a.Delete(); err != nil {
				log.Println(err)
				return 0, err
			}
			log.Printf("FS: Orphan Scan: Removed File: %v", filename)
			artCount++
		}

		// Scan for all songs NOT under the base folder
		songs, err := db.DB.SongsNotInPath(baseFolder)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		// Remove all songs which are not in this path
		for _, s := range songs {
			// Remove song from database
			filename := s.FileName
			if err := s.Delete(); err != nil {
				log.Println(err)
				return 0, err
			}
			log.Printf("FS: Orphan Scan: Removed File: %v", filename)
			songCount++
		}

		// Scan for all folders NOT under the base folder
		folders, err := db.DB.FoldersNotInPath(baseFolder)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		// Remove all folders which are not in this path
		for _, f := range folders {
			// Remove folder from database
			path := f.Path
			if err := f.Delete(); err != nil {
				log.Println(err)
				return 0, err
			}
			log.Printf("FS: Orphan Scan: Removed Path: %v", path)
			folderCount++
		}
	}

	// If no subfolder set, use the base folder to check file existence
	if subFolder == "" {
		subFolder = baseFolder
	}

	if fs.verbose {
		log.Println("FS: Orphan Scanning: Scanning subfolder: ", subFolder)
	} else {
		log.Println("FS: Orphan Scan: Removing: ", subFolder)
	}

	// Scan for all art in subfolder
	art, err := db.DB.ArtInPath(subFolder)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// Iterate all art in this path
	for _, a := range art {
		// Check that the art still exists in this place
		if _, err := os.Stat(a.FileName); os.IsNotExist(err) {
			// Remove art from database
			filename := a.FileName
			if err := a.Delete(); err != nil {
				log.Println(err)
				return 0, err
			}
			log.Printf("FS: Orphan Scan: File Does Not Exist: %v", filename)
			artCount++
		}
	}

	// Scan for all songs in subfolder
	songs, err := db.DB.SongsInPath(subFolder)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// Iterate all songs in this path
	for _, s := range songs {
		// Check that the song still exists in this place
		if _, err := os.Stat(s.FileName); os.IsNotExist(err) {
			// Remove song from database
			filename := s.FileName
			if err := s.Delete(); err != nil {
				log.Println(err)
				return 0, err
			}
			log.Printf("FS: Orphan Scan: File Does Not Exist: %v", filename)
			songCount++
		}
	}

	// Scan for all folders in subfolder
	folders, err := db.DB.FoldersInPath(subFolder)
	if err != nil {
		return 0, err
	}

	// Iterate all folders in this path
	for _, f := range folders {
		// Check that the folder still has items within it
		files, err := ioutil.ReadDir(f.Path)
		if err != nil && !os.IsNotExist(err) {
			log.Println(err)
			return 0, err
		}

		// Delete any folders with 0 items
		if len(files) == 0 {
			path := f.Path
			if err := f.Delete(); err != nil {
				log.Println(err)
				return 0, err
			}
			log.Printf("FS: Orphan Scan: Folder Has No Items: %v", path)
			folderCount++
		}
	}

	// Now that songs have been purged, check for albums
	albumCount, err := db.DB.PurgeOrphanAlbums()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// Check for artists
	artistCount, err := db.DB.PurgeOrphanArtists()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// Print metrics
	if fs.verbose {
		log.Printf("FS: Orphan Scan: Complete [time: %s]", time.Since(startTime).String())
		log.Printf("FS: Orphan Scan: Removed: [art: %d] [artists: %d] [albums: %d] [songs: %d] [folders: %d]",
			artCount, artistCount, albumCount, songCount, folderCount)
	}

	// Sum up changes
	sum := artCount + artistCount + albumCount + songCount + folderCount

	return sum, nil
}
