package fs

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"aiomst/db"
)

// Cache repetitive entries
var folderCache = map[string]*db.Folder{}
var artistCache = map[string]*db.Artist{}
var albumCache  = map[string]*db.Album{}

// Track folder IDs containing new art and hold the art IDs.
var artFiles []folderArtPair

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
	artFiles := make([]folderArtPair, 0)
	startTime := time.Now()

	folderCache = map[string]*db.Folder{}
	artistCache = map[string]*db.Artist{}
	albumCache  = map[string]*db.Album{}

	if fs.verbose	{
		log.Printf("FS: Scanning: %s", baseFolder)
	}

	err := filepath.Walk(baseFolder, func(cPath string, info os.FileInfo, e error) error	{
		mutex.RLock()
		if haltWalk {
			return errors.New("FS: Media Scan: Halted by channel")
		}
		mutex.RUnlock()

		// This should never happen but just to be sure
		if info == nil 	{
			return errors.New("FS: Media Scan: invalid path: " + cPath)
		}
		
		log.Printf("FS: Media Scan: Got new file: %s", cPath)
		folder, err := handleFolder(cPath, info)
		if err != nil {
			return err
		}
		if folder != nil {
			folderCount++
		}

		ext := path.Ext(cPath)
		if img, audio := imgType[ext], audioType[ext]; !img && !audio {
			return nil
		}
		
		if _, ok := imgType[ext]; ok {
			art, err := handleImg(cPath, info)
			if err != nil {
				return err
			}
			if art != nil {
				artFiles = append(artFiles, folderArtPair {
					folderID: folder.ID,
					artID: art.ID,
				})
				artCount++
			}
			return nil
		}

		if _, ok := audioType[ext]; ok {
			changes, err := handleAudio(cPath, info, folder)
			if err != nil {
				return err
			}
			if err == nil {
				artistCount += changes[0]
				albumCount  += changes[1]
				songCount   += changes[2]
				songUpdate  += changes[3]
			}
		}
		return nil
	})
	
	if err != nil {
		return 0, err
	}

	for _, a := range artFiles {
		songs, err := db.DB.SongsForFolder(a.folderID)
		if err != nil {
			return 0, err
		}

		for _, s := range songs {
			s.ArtID = a.artID
			if err := s.Update(); err != nil {
				return 0, err
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

func handleFolder(cPath string, info os.FileInfo)	(*db.Folder, error) {
	folder := new(db.Folder)
	if info.IsDir() {
		folder.Path = cPath
	}	else	{
		folder.Path = path.Dir(cPath)
	}

	// Check for cached folder
	if seenFolder, ok := folderCache[folder.Path]; ok	{
		folder = seenFolder
	}	else if err := folder.Load(); err != nil && err == sql.ErrNoRows  {
		files, err := ioutil.ReadDir(folder.Path)
		if err != nil {
			return nil, err
		}
		
		log.Printf("FS: Media Scan: Found %v files in %v", len(files), cPath)
		if len(files) == 0 {
			return nil, nil
		}

		folder.Title = path.Base(folder.Path)
		parent := new(db.Folder)
		if info.IsDir() {
			parent.Path = path.Dir(cPath)
		} else {
			parent.Path = path.Dir(path.Dir(cPath))
		}

		if err := parent.Load(); err != nil && err != sql.ErrNoRows {
			log.Print(err)
			return nil, err
		}

		folder.ParentID = parent.ID
		if err := folder.Save(); err != nil {
			log.Print(err)
			return nil, err
		}
	}

	folderCache[folder.Path] = folder
	return folder, nil
}

func handleImg(cPath string, info os.FileInfo) (*db.Art, error) {
	art := new(db.Art)
	art.FileName = cPath
	if err := art.Load(); err == sql.ErrNoRows {
		art.FileSize = info.Size()
		art.LastModified = info.ModTime().Unix()

		if art.FileSize == 0 {
			return nil, nil
		}

		if err := art.Save(); err != nil {
			return nil, err
		} 
		return art, nil
	}
	return nil, nil
}

func handleAudio(cPath string, info os.FileInfo, folder *db.Folder) ([]int, error)	{
	song, err := db.SongFromFile(cPath)
	if err != nil {
		return []int{}, fmt.Errorf("FS: Media Scan: Handle Audio Error: %v", err)
	}

	if info.Size() == 0 {
		return []int{}, nil
	}

	song.FileName = cPath
	song.FileSize = info.Size()
	song.LastModified = info.ModTime().Unix()
	song.FolderID = folder.ID
	song.FileTypeID = db.FileTypeMap[path.Ext(cPath)]

	artist, countArtist := handleArtist(song)
	song.ArtistID = artist.ID

	album, countAlbum := handleAlbum(song)
	song.AlbumID = album.ID
	
	countSong, countUpdate := checkForModification(song)
	return []int{countArtist, countAlbum, countSong, countUpdate}, nil
}

func handleArtist(song *db.Song) (*db.Artist, int) {
	count := 0
	artist := db.ArtistFromSong(song)
	if tempArtist, ok := artistCache[artist.Title]; ok {
		artist = tempArtist
	} else if err := artist.Load(); err == sql.ErrNoRows {
		if err := artist.Save(); err != nil {
			log.Printf("FS: Media Scan: Handle Artist: %v", err)
		} else if err == nil {
			log.Printf("FS: Media Scan: Artist: [#%05d] %s", artist.ID, artist.Title)
		  count++
		}
	}
	
	artistCache[artist.Title] = artist
	return artist, count
}

func handleAlbum(song *db.Song) (*db.Album, int) {
	count := 0
	album := db.AlbumFromSong(song)
	album.ArtistID = song.ArtistID
	albumCacheKey := strconv.Itoa(album.ArtistID) + "_" + album.Title
	if temp, ok := albumCache[albumCacheKey]; ok {
		album = temp
	} else if err := album.Load(); err == sql.ErrNoRows {
		if err := album.Save(); err != nil {
			log.Printf("FS: Media Scan: Handle Album: %v", err)
		} else if err == nil {
			log.Printf("FS: Media Scan: Album: [#%05d] %s - %d - %s", album.ID, album.Artist, album.Year, album.Title)
			count++
		}
	}
	
	albumCache[albumCacheKey] = album
	return album, count
}

func checkForModification(origin *db.Song) (int, int)	{
	songCount := 0
	songUpdate := 0
	song2 := new(db.Song)
	song2.FileName = origin.FileName

	// Check if the song exists
	if err := song2.Load(); err == sql.ErrNoRows {
		// The song didn't save. So do that...
		if err2 := origin.Save(); err2 != nil && err2 != sql.ErrNoRows {
			log.Println(err2)
		}	else if err2 == nil {
			songCount++
		}
	} else {
		// Song already exists. Check for updates.
		if origin.LastModified > song2.LastModified {
			// Update Existing
			origin.ID = song2.ID
			if err2 := origin.Update(); err2 != nil {
				log.Println(err2)
			}
			songUpdate++
		}
	}

	return songCount, songUpdate
}