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

// Track folder IDs to db.Art, db.Artist and db.Album
var folderToObjs = map[int][]interface{}{}

// Track folder IDs containing new art and hold the art IDs.
var artFiles map[int]int = make(map[int]int)

// Track folder IDs containing metadata and hold the metadata IDs.
var mdFiles map[int]int = make(map[int]int)

// Map folder IDs to db.Artist or db.Album
var attachArt map[int]AttachesArt = make(map[int]AttachesArt)

// Map folder IDs to db.Artist or db.Album
var hasMdItems map[int]HasMetadata = make(map[int]HasMetadata)

// Track Metrics
var artCount      int = 0
var artistCount   int = 0
var albumCount    int = 0
var songCount     int = 0
var songUpdate    int = 0
var folderCount   int = 0
var metadataCount int = 0

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
	artCount      = 0
	artistCount   = 0
	albumCount    = 0
	songCount     = 0
	songUpdate    = 0
	folderCount   = 0
	metadataCount = 0
	startTime := time.Now()

	folderCache = map[string]*db.Folder{}
	artistCache = map[string]*db.Artist{}
	albumCache  = map[string]*db.Album{}
	folderToObjs = map[int][]interface{}{}

	if fs.verbose	{
		log.Printf("FS: Scanning: %s", baseFolder)
	}

	godirwalk.Walk(baseFolder, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error	{
			info := de.ModeType()
			log.Printf("FS: Media Scan: Got new file: %s", osPathname)
			
			folder, err := handleFolder(osPathname, info)
			if err != nil {
				return fmt.Errorf("FS: Media Scan: Error handling folder: %s", err)
			}
			if _, ok := folderToObjs[folder.ID] ; !ok {
				folderToObjs[folder.ID] = make([]interface{}, 0)
			}
			
			ext := path.Ext(osPathname)
			var isMetadata bool = false
			if osPathname[len(osPathname)-len(metadataFile):] == metadataFile {
				isMetadata = true
			}

			if img, audio := imgType[ext], audioType[ext];
			!img && !audio && !isMetadata {
				return nil
			}
			
			if _, ok := imgType[ext]; ok {
				art, err := handleImg(osPathname, info)
				if err != nil {
					return fmt.Errorf("FS: Media Scan: Error handling image: %s", err)
				}
				if art != nil {
					artFiles[folder.ID] = art.ID
				}
				return nil
			}

			if isMetadata {
				md, err := handleMetadata(osPathname, folder)
				if err != nil {
					return fmt.Errorf("FS: Media Scan: Error handling Metadata: %s", err)
				}
				mdFiles[folder.ID] = md.ID
				return nil
			}

			if _, ok := audioType[ext]; ok {
				err := handleAudio(osPathname, info, folder)
				if err != nil {
					return fmt.Errorf("FS: Media Scan: Error handling audio file: %s", err)
				}
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
			if art, _ := v.GetArt(); art != nil {
				continue
			}
			if aID != 0 {		 
				if err := v.SetArtID(aID); err != nil {
					log.Printf("FS: Media Scan: Attach Art Error: ", err)
				}
			}
		}
	}

	for fID, mID := range mdFiles {
		if v, member := hasMdItems[fID]; member {
			if md, _ := v.GetMetadata(); md != nil {
				continue
			}
			if mID != 0 {
				if err := v.SetMetadataID(mID); err != nil {
					log.Printf("FS: Media Scan: Metadata attachment Error: %v", err)
				}
			}
		}
	}

	if fs.verbose {
		log.Printf("FS: Media Scan Complete [time: %s]", time.Since(startTime).String())
		log.Printf(
			"FS: Media Scan: Added: " +
			"[art: %d] [artists: %d] [albums: %d] [songs: %d] " + 
			"[folders: %d] [metadata: %d]",
			artCount, artistCount, albumCount, songCount, folderCount, metadataCount)
		
		log.Printf("FS: Updated: [songs: %d]", songUpdate)
	}

	sum := artCount + artistCount + albumCount + songCount + folderCount
	return sum, nil
}

func handleFolder(cPath string, info os.FileMode)	(*db.Folder, error) {
	// Check for cached folder
	if seenFolder, ok := folderCache[cPath]; ok	{
		return seenFolder, nil
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
		return &existing, nil
	}

	if err == sql.ErrNoRows  {
		if _, err := os.Stat(path.Join(cPath, metadataFile)); os.IsNotExist(err) {
			if _, err := os.Create(path.Join(cPath, metadataFile)); err != nil {
				log.Printf("FS: Media Scan: Error creating Metadata file: %v", err)
			} else {
				log.Printf(
					"FS: Media Scan: Handle Folder: Created new Metadata File at: %#v",
					path.Join(cPath, metadataFile))
			}
		}

		files, err := godirwalk.ReadDirents(folder.Path, nil)
		if err != nil {
			return nil, err
		} else if len(files) == 0 {
			return nil, nil
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
			return nil, err
		}

		folder.ParentID = parent.ID
		if err := folder.Save(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	folderCache[folder.Path] = folder
	folderCount++
	return folder, nil
}

func handleImg(cPath string, info os.FileMode) (*db.Art, error) {
	art := new(db.Art)
	art.Path = cPath

	existing, err := art.Load()
	if existing.Path == art.Path {
		return &existing, nil
	}

	if err == sql.ErrNoRows {
		data, err := os.Stat(cPath)
		if err != nil {
			return nil, err
		}

		art.FileSize = data.Size()
		art.LastModified = data.ModTime().Unix()
		if art.FileSize == 0 {
			return nil, errors.New("Art File Size is 0")
		}
		if err := art.Save(); err != nil {
			return nil, err
		} 

		artCount++
		return art, nil
	}
	return nil, err
}

func handleMetadata(cPath string, folder *db.Folder) (*db.Metadata, error) {
	md := new(db.Metadata)
	md.Path = cPath
	
	existing, err := md.Load()
	if existing.Path == md.Path {
		return &existing, nil
	}

	if err == sql.ErrNoRows {
		data, err := os.Stat(cPath)
		if err != nil {
			return nil, err
		}

		mdDir := filepath.Dir(filepath.Dir(md.Path))
		if folder, ok := folderCache[mdDir]; ok {
			md.FolderID = folder.ID
		} else {
			md.FolderID = 0
		}
		
		md.FileSize = data.Size()
		md.LastModified = data.ModTime().Unix()
		if err := md.Save(); err != nil {
			return nil, err
		} 

		metadataCount++
		return md, nil
	}
	return nil, err
}

func handleAudio(cPath string, info os.FileMode, folder *db.Folder) error	{
	song, err := db.SongFromFile(cPath)
	if err != nil {
		return err
	}

	data, err := os.Stat(cPath)
	if err != nil {
		return err
	}
	if data.Size() == 0 {
		return errors.New("Audio File Size is 0")
	}

	song.Path = cPath
	song.FileSize = data.Size()
	song.LastModified = data.ModTime().Unix()
	song.FolderID = folder.ID
	song.FileTypeID = db.FileTypeMap[path.Ext(cPath)]

	artist, err := handleArtist(song)
	if err != nil {
		return err
	}
	song.ArtistID = artist.ID

	album, err := handleAlbum(song)
	if err != nil {
		return err
	}
	song.AlbumID = album.ID
	
	err = checkForModification(song)
	if err != nil {
		return err
	}
	return nil
}

func handleArtist(song *db.Song) (*db.Artist, error) {
	artist := db.GetArtistFromSong(song)
	if seenArtist, ok := artistCache[artist.Title]; ok {
		return seenArtist, nil
	}

	existing, err := artist.Load()
	if existing != (db.Artist{}) {
		artistCache[artist.Title] = &existing
		return &existing, nil
	}

	if err == sql.ErrNoRows 	{
		artistDir := filepath.Dir(filepath.Dir(song.Path))
		artist.FolderID = folderCache[artistDir].ID
		artist.MBID = errMBIDStartValue
		artist.MetadataID = 0
		if err := artist.Save(); err != nil {
			return nil, fmt.Errorf("FS: Media Scan: Handle Artist: Error Saving: %v", err)
		} 

		attachArt[artist.FolderID] = artist
		hasMdItems[artist.FolderID] = artist
		artistCache[artist.Title] = artist
		log.Printf("FS: Media Scan: Artist: [#%05d] %s", artist.ID, artist.Title)
		
		artistCount++
		return artist, nil
	}	
	
	return nil, fmt.Errorf("FS: Media Scan: Handle Artist: Other Error: %v", err)
}

func handleAlbum(song *db.Song) (*db.Album, error) {
	album := db.GetAlbumFromSong(song)
	album.ArtistID = song.ArtistID
	albumCacheKey := strconv.Itoa(album.ArtistID) + "_" + album.Title 
	if seenAlbum, ok := albumCache[albumCacheKey]; ok {
		return seenAlbum, nil
	} 

	existing, err := album.Load();
	if existing != (db.Album{}) {
		albumCache[albumCacheKey] = &existing
		return &existing, nil
	}

	if err == sql.ErrNoRows {
		album.ArtistID = song.ArtistID
		album.FolderID = song.FolderID
		album.MBID = errMBIDStartValue
		album.MetadataID = 0
		if err := album.Save(); err != nil {
			return nil, fmt.Errorf(
				"FS: Media Scan: Handle Album: Error Saving: %v", err)
		}
		
		attachArt[album.FolderID] = album
		hasMdItems[album.FolderID] = album
		albumCache[albumCacheKey] = album
		log.Printf(
			"FS: Media Scan: Album: [#%05d] %s - %d - %s",
			album.ID, album.Artist, album.Year, album.Title)

		albumCount++
		return album, nil
	}

	return nil, fmt.Errorf("FS: Media Scan: Handle Album: Other Error: %s", err)
}

func checkForModification(origin *db.Song) error	{
	song2 := new(db.Song)
	song2.Path = origin.Path

	// Check if the song exists
	existing, err := song2.Load()
	if existing != (db.Song{})  {
		if origin.LastModified > existing.LastModified {
			// Update Existing
			origin.ID = song2.ID
			if err2 := origin.Update(); err2 != nil {
				return fmt.Errorf("FS: Media Scan: Check Modifications: %v", err2)
			}

			songUpdate++
			return nil
		}

		// The song exists, but it hasn't been changed
		origin = &existing
		return nil
	}

	if err == sql.ErrNoRows {
		// The song didn't save. So do that...
		origin.MBID = errMBIDStartValue
		if err2 := origin.Save(); err2 != nil && err2 != sql.ErrNoRows {
			return fmt.Errorf("FS: Media Scan: Check Modifications: %v", err2)
		}	else if err2 == nil {
			songCount++
			return nil
		}
	}	
	return fmt.Errorf("FS: Media Scan: Check Modification Other Error: %v", err)
}