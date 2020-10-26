package fs

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/eHoward1996/aiomst/db"

	"github.com/karrick/godirwalk"
)

// Cache repetitive entries
var folderCache = map[string]*db.Folder{}
var artistCache = map[string]*db.Artist{}
var albumCache  = map[string]*db.Album{}

// Track folder IDs containing new art and hold the art IDs.
var artFiles map[int]int

// Map folder IDs to db.Artist or db.Album
var attachArt map[int]AttachesArt = make(map[int]AttachesArt)

// MediaScan is a filesystem task that scans the given path for new media
type MediaScan struct	{
	baseFolder 	string
	subFolder 	string
	verbose 		bool
}

// Folders returns the base and sub folders for scanning
func (fs *MediaScan) Folders()	(string, string)	{
	return fs.baseFolder, fs.subFolder
}

// SetFolders sets the base and sub folders for scanning
func (fs *MediaScan) SetFolders(baseFolder, subFolder string)	{
	fs.baseFolder = baseFolder
	fs.subFolder  = subFolder
}

// Verbose is whether scanning has verbose output or not
func (fs *MediaScan) Verbose(v bool)	{
	fs.verbose = v
}

// WhoAmI returns Media Scan
func (fs *MediaScan) WhoAmI() string {
	return "Media Scan"
}

// Scan scans for media files in the filesystem
func (fs *MediaScan) Scan(baseFolder, subFolder string, walkCancelChan chan struct{}) (int, error) {
	// Halt file system walk if needed
	var mutex sync.RWMutex
	haltWalk := false
	go func()	{
		// Wait until signal recieved
		<- walkCancelChan

		// Halt fs walk
		mutex.Lock()
		haltWalk = true
		mutex.Unlock()
	}()

	// Track metrics
	artCount    := 0
	artistCount := 0
	albumCount  := 0
	songCount   := 0
	songUpdate  := 0
	folderCount := 0
	artFiles := make(map[int]int)
	startTime := time.Now()

	folderCache = map[string]*db.Folder{}
	artistCache = map[string]*db.Artist{}
	albumCache  = map[string]*db.Album{}

	if fs.verbose	{
		log.Printf("FS: Scanning: %s", baseFolder)
	}

	godirwalk.Walk(baseFolder, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error	{
			info := de.ModeType()
			log.Printf("FS: Media Scan: Got new file: %s", osPathname)
			
			folder, inc, err := handleFolder(osPathname, info)
			if err != nil {
				return fmt.Errorf("FS: Media Scan: Error handling folder: %s", err)
			}
			if inc {
				folderCount++
			}
			
			ext := path.Ext(osPathname)
			if img, audio := imgType[ext], audioType[ext]; !img && !audio {
				return nil
			}
			
			if _, ok := imgType[ext]; ok {
				art, inc, err := handleImg(osPathname, info)
				if err != nil {
					return fmt.Errorf("FS: Media Scan: Error handling image: %s", err)
				}
				if art != nil {
					artFiles[folder.ID] = art.ID
					if inc {
						artCount++
					}
				}
				return nil
			}

			if _, ok := audioType[ext]; ok {
				changes, err := handleAudio(osPathname, info, folder)
				if err != nil {
					return fmt.Errorf("FS: Media Scan: Error handling audio file: %s", err)
				}

				artistCount += changes[0]
				albumCount  += changes[1]
				songCount   += changes[2]
				songUpdate  += changes[3]
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			log.Printf("%s", err)
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	for fID, aID := range artFiles {
		if v, member := attachArt[fID]; member {
			if aID != 0 && aID != v.GetArtID() {		 
				if err := v.SetArtID(aID); err != nil {
					log.Printf("FS: Media Scan: Attach Art Error: ", err)
				}
			}
		}
	}

	if fs.verbose {
		log.Printf("FS: Media Scan Complete [time: %s]", time.Since(startTime).String())
		log.Printf("FS: Media Scan: Added: [art: %d] [artists: %d] [albums: %d] [songs: %d] [folders: %d]",
			artCount, artistCount, albumCount, songCount, folderCount)
		log.Printf("FS: Updated: [songs: %d]", songUpdate)
	}

	sum := artCount + artistCount + albumCount + songCount + folderCount
	return sum, nil
}

func handleFolder(cPath string, info os.FileMode)	(*db.Folder, bool, error) {
	// Check for cached folder
	if seenFolder, ok := folderCache[cPath]; ok	{
		return seenFolder, false, nil
	}

	folder := new(db.Folder)
	if info.IsDir() {
		folder.Path = cPath
	}	else	{
		folder.Path = path.Dir(cPath)
	}
	
	existing, err := folder.Load()
	if existing != (db.Folder{}) {
		folderCache[cPath] = &existing
		return &existing, false, nil
	}

	if err == sql.ErrNoRows  {
		files, err := godirwalk.ReadDirents(folder.Path, nil)
		if err != nil {
			return nil, false, err
		} else if len(files) == 0 {
			return nil, false, nil
		}

		log.Printf("FS: Media Scan: Found %v files in %v", len(files), cPath)
		folder.Title = path.Base(folder.Path)
		parent := new(db.Folder)
		if info.IsDir() {
			parent.Path = path.Dir(cPath)
		} else {
			parent.Path = path.Dir(path.Dir(cPath))
		}

		if _, err := parent.Load(); err != nil && err != sql.ErrNoRows {
			return nil, false, err
		}

		folder.ParentID = parent.ID
		if err := folder.Save(); err != nil {
			return nil, false, err
		}
	} else {
		return nil, false, err
	}

	folderCache[folder.Path] = folder
	return folder, true, nil
}

func handleImg(cPath string, info os.FileMode) (*db.Art, bool, error) {
	art := new(db.Art)
	art.FileName = cPath

	existing, err := art.Load()
	if existing.FileName == art.FileName {
		return &existing, false, nil
	}

	if err == sql.ErrNoRows {
		data, err := os.Stat(cPath)
		if err != nil {
			return nil, false, err
		}

		art.FileSize = data.Size()
		art.LastModified = data.ModTime().Unix()

		if art.FileSize == 0 {
			return nil, false, errors.New("Art File Size is 0")
		}
		if err := art.Save(); err != nil {
			return nil, false, err
		} 
		return art, true, nil
	}
	return nil, false, err
}

func handleAudio(cPath string, info os.FileMode, folder *db.Folder) ([]int, error)	{
	song, err := db.SongFromFile(cPath)
	if err != nil {
		return []int{}, err
	}

	data, err := os.Stat(cPath)
	if err != nil {
		return []int{}, err
	}
	if data.Size() == 0 {
		return []int{}, errors.New("Audio File Size is 0")
	}

	song.FileName = cPath
	song.FileSize = data.Size()
	song.LastModified = data.ModTime().Unix()
	song.FolderID = folder.ID
	song.FileTypeID = db.FileTypeMap[path.Ext(cPath)]

	artist, countArtist, err := handleArtist(song)
	if err != nil {
		return []int{}, err
	}
	song.ArtistID = artist.ID

	album, countAlbum, err := handleAlbum(song)
	if err != nil {
		return []int{}, err
	}
	song.AlbumID = album.ID
	
	countSong, countUpdate, err := checkForModification(song)
	if err != nil {
		return []int{}, err
	}
	return []int{countArtist, countAlbum, countSong, countUpdate}, nil
}

func handleArtist(song *db.Song) (*db.Artist, int, error) {
	artist := db.GetArtistFromSong(song)
	if seenArtist, ok := artistCache[artist.Title]; ok {
		return seenArtist, 0, nil
	}

	existing, err := artist.Load()
	if existing != (db.Artist{}) {
		artistCache[artist.Title] = &existing
		return &existing, 0, nil
	}

	if err == sql.ErrNoRows 	{
		artistDir := filepath.Dir(filepath.Dir(song.FileName))
		artist.FolderID = folderCache[artistDir].ID
		if err := artist.Save(); err != nil {
			return nil, 0, fmt.Errorf("FS: Media Scan: Handle Artist: Error Saving: %v", err)
		} 

		attachArt[artist.FolderID] = artist
		artistCache[artist.Title] = artist
		log.Printf("FS: Media Scan: Artist: [#%05d] %s", artist.ID, artist.Title)
		return artist, 1, nil
	}	
	
	return nil, 0, fmt.Errorf("FS: Media Scan: Handle Artist: Other Error: %v", err)
}

func handleAlbum(song *db.Song) (*db.Album, int, error) {
	album := db.GetAlbumFromSong(song)
	album.ArtistID = song.ArtistID
	albumCacheKey := strconv.Itoa(album.ArtistID) + "_" + album.Title 
	if seenAlbum, ok := albumCache[albumCacheKey]; ok {
		return seenAlbum, 0, nil
	} 

	existing, err := album.Load();
	if existing != (db.Album{}) {
		albumCache[albumCacheKey] = &existing
		return &existing, 0, nil
	}

	if err == sql.ErrNoRows {
		album.ArtistID = song.ArtistID
		album.FolderID = folderCache[filepath.Dir(song.FileName)].ID
		if err := album.Save(); err != nil {
			return nil, 0, fmt.Errorf("FS: Media Scan: Handle Album: Error Saving: %v", err)
		}
		
		attachArt[album.FolderID] = album
		albumCache[albumCacheKey] = album
		log.Printf("FS: Media Scan: Album: [#%05d] %s - %d - %s", album.ID, album.Artist, album.Year, album.Title)
		return album, 1, nil
	}

	return nil, 0, fmt.Errorf("FS: Media Scan: Handle Album: Other Error: %s", err)
}

func checkForModification(origin *db.Song) (int, int, error)	{
	song2 := new(db.Song)
	song2.FileName = origin.FileName

	// Check if the song exists
	existing, err := song2.Load()
	if existing != (db.Song{})  {
		if origin.LastModified > existing.LastModified {
			// Update Existing
			origin.ID = song2.ID
			if err2 := origin.Update(); err2 != nil {
				return 0, 0, fmt.Errorf("FS: Media Scan: Check Modifications: %v", err2)
			}
			return 0, 1, nil
		}

		// The song exists, but it hasn't been changed
		origin = &existing
		return 0, 0, nil
	}

	if err == sql.ErrNoRows {
		// The song didn't save. So do that...
		if err2 := origin.Save(); err2 != nil && err2 != sql.ErrNoRows {
			return 0, 0, fmt.Errorf("FS: Media Scan: Check Modifications: %v", err2)
		}	else if err2 == nil {
			return 1, 0, nil
		}
	}
	
	return 0, 0, fmt.Errorf("FS: Media Scan: Check Modification Other Error: %v", err)
}