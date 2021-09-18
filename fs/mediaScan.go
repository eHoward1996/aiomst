package fs

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/karrick/godirwalk"
)

// MediaScan is a filesystem task that scans the given path for new media
type MediaScan struct	{
	baseFolder 	string
	subFolder 	string
	verbose 		bool
}
type attachables struct {
	art *db.Art
	md  *db.Metadata
}

// Cache repetitive entries
var folderCache = map[string]*db.Folder{}
var artistCache = map[string]*db.Artist{}
var albumCache  = map[string]*db.Album{}

// Track folder IDs to HasAttachables (db.Artist & db.Album)
var fIDHasAttachables = map[int]HasAttachables{}

// Track folder IDs to Attachables (db.Art & db.Metadata)
var folderAttachables = map[int]attachables{}

// Track Metrics
var artCount      int = 0
var artistCount   int = 0
var albumCount    int = 0
var songCount     int = 0
var songUpdate    int = 0
var folderCount   int = 0
var metadataCount int = 0


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
func (fs *MediaScan) Scan(
	baseFolder, subFolder string, walkCancelChan chan struct{}) (int, error) {
		// Halt file system walk if needed
		var mutex sync.RWMutex
		haltWalk := false
		go func(haltWalk bool)	{
			// Wait until signal recieved
			<- walkCancelChan

			// Halt fs walk
			mutex.Lock()
			haltWalk = true
			mutex.Unlock()
		}(haltWalk)

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
		
		if fs.verbose	{
			util.Logger.Printf("FS: Scanning: %s", baseFolder)
		}

		godirwalk.Walk(baseFolder, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error	{
				info := de.ModeType()
				util.Logger.Printf("FS: Media Scan: Got new file: %s", osPathname)
				
				folder, err := handleFolder(osPathname, info)
				if err != nil {
					return fmt.Errorf("FS: Media Scan: Error handling folder: %s", err)
				}
				if _, ok := folderAttachables[folder.ID]; !ok {
					folderAttachables[folder.ID] = attachables{}
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
					art, err := handleImg(osPathname, folder)
					if err != nil {
						return fmt.Errorf("FS: Media Scan: Error handling image: %s", err)
					}
					currentVal := folderAttachables[folder.ID]
					currentVal.art = art
					folderAttachables[folder.ID] = currentVal
					return nil
				}

				if isMetadata {
					md, err := handleMetadata(osPathname, folder)
					if err != nil {
						return fmt.Errorf(
							"FS: Media Scan: Error handling Metadata: %s", err)
					}
					currentVal := folderAttachables[folder.ID]
					currentVal.md = md
					folderAttachables[folder.ID] = currentVal
					return nil
				}

				if _, ok := audioType[ext]; ok {
					err := handleAudio(osPathname, info, folder)
					if err != nil {
						return fmt.Errorf(
							"FS: Media Scan: Error handling audio file: %s", err)
					}
				}
				return nil
			},
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				util.Logger.Printf("%s", err)
				return godirwalk.SkipNode
			},
			Unsorted: true,
		})
		
		joinAttachables()
		if err := db.DB.TruncateLog(); err != nil {
			util.Logger.Printf("FS: Media Scan: Could not truncate WAL File: %v", err)
		}

		if fs.verbose {
			util.Logger.Printf(
				"FS: Media Scan Complete [time: %s]", time.Since(startTime).String())
			util.Logger.Printf(
				"FS: Media Scan: Added: " +
				"[art: %d] [artists: %d] [albums: %d] [songs: %d] " + 
				"[folders: %d] [metadata: %d]",
				artCount, artistCount, albumCount,
				songCount, folderCount, metadataCount,
			)
			
			util.Logger.Printf("FS: Updated: [songs: %d]", songUpdate)
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
	
	err := folder.Load()
	if err == nil {
		folderCache[cPath] = folder
		return folder, nil
	}

	if err == sql.ErrNoRows  {
		if _, err := os.Stat(path.Join(cPath, metadataFile)); os.IsNotExist(err) {
			if _, err := os.Create(path.Join(cPath, metadataFile)); err != nil {
				util.Logger.Printf("FS: Media Scan: Error creating Metadata file: %v", err)
			} else {
				util.Logger.Printf(
					"FS: Media Scan: Handle Folder: Created new Metadata File at: %#v",
					path.Join(cPath, metadataFile))
			}
		}

		files, err := godirwalk.ReadDirents(folder.Path, nil)
		if err != nil {
			return nil, err
		} else if len(files) == 0 {
			return nil, fmt.Errorf(
				"FS: Media Scan: Found no files in folder: %v", folder.Path)
		}

		util.Logger.Printf("FS: Media Scan: Found %v files in %v", len(files), cPath)
		folder.Title = path.Base(folder.Path)
		parent := new(db.Folder)
		if info.IsDir() {
			parent.Path = path.Dir(cPath)
		} else {
			parent.Path = path.Dir(path.Dir(cPath))
		}

		if err := parent.Load(); err != nil && err != sql.ErrNoRows {
			return nil, err
		} else if err == nil {
			folder.ParentID = parent.ID
		}

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

func handleImg(cPath string, folder *db.Folder) (*db.Art, error) {
	art := new(db.Art)
	art.Path = cPath

	err := art.Load()
	if err == nil {
		return art, nil
	}

	if err == sql.ErrNoRows {
		data, err := os.Stat(cPath)
		if err != nil {
			return nil, err
		}

		art.FolderID = folder.ID
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
	
	err := md.Load()
	if err == nil {
		return md, nil
	}

	if err == sql.ErrNoRows {
		data, err := os.Stat(cPath)
		if err != nil {
			return nil, err
		}

		md.FolderID = folder.ID
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
	fIDHasAttachables[artist.FolderID] = artist

	album, err := handleAlbum(song)
	if err != nil {
		return err
	}
	song.AlbumID = album.ID
	fIDHasAttachables[album.FolderID] = album

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

	err := artist.Load()
	if err == nil {
		artistCache[artist.Title] = artist
		return artist, nil
	}

	if err == sql.ErrNoRows 	{
		f := new(db.Folder)
		f.ID = song.FolderID
		
		artist.FolderID = 0
		if err := db.DB.LoadFolder(f); err == nil {
			artist.FolderID = f.ParentID
		}
		artist.MBID = errStartValueString
		artist.DiscogsID = errStartValueString
		artist.MetadataID = 0
		if err := artist.Save(); err != nil {
			return nil, fmt.Errorf(
				"FS: Media Scan: Handle Artist: Error Saving: %v", err)
		} 

		artistCount++
		artistCache[artist.Title] = artist
		util.Logger.Printf("FS: Media Scan: Artist: [#%05d] %s", artist.ID, artist.Title)		
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

	err := album.Load()
	if err == nil {
		albumCache[albumCacheKey] = album
		return album, nil
	}

	if err == sql.ErrNoRows {
		album.ArtistID = song.ArtistID
		album.FolderID = song.FolderID
		album.MBID = errStartValueString
		album.DiscogsID = errStartValueInt
		album.MetadataID = 0
		if err := album.Save(); err != nil {
			return nil, fmt.Errorf(
				"FS: Media Scan: Handle Album: Error Saving: %v", err)
		}
		
		albumCount++
		albumCache[albumCacheKey] = album
		util.Logger.Printf(
			"FS: Media Scan: Album: [#%05d] %s - %d - %s",
			album.ID, album.Artist, album.Year, album.Title)
		return album, nil
	}

	return nil, fmt.Errorf("FS: Media Scan: Handle Album: Other Error: %s", err)
}

func checkForModification(origin *db.Song) error	{
	song2 := new(db.Song)
	song2.Path = origin.Path

	// Check if the song exists
	err := song2.Load()
	if err == nil  {
		if origin.LastModified > song2.LastModified {
			// Update Existing
			origin.ID = song2.ID
			if err2 := origin.Update(); err2 != nil {
				return fmt.Errorf("FS: Media Scan: Check Modifications: %v", err2)
			}

			songUpdate++
			return nil
		}

		// The song exists, but it hasn't been changed
		origin = song2
		return nil
	}

	if err == sql.ErrNoRows {
		// The song didn't save. So do that...
		origin.MBID = errStartValueString
		if err2 := origin.Save(); err2 != nil && err2 != sql.ErrNoRows {
			return fmt.Errorf("FS: Media Scan: Check Modifications: %v", err2)
		}	else if err2 == nil {
			songCount++
			return nil
		}
	}	
	return fmt.Errorf("FS: Media Scan: Check Modification Other Error: %v", err)
}

func joinAttachables() {
	for k, v := range fIDHasAttachables {
		if _, ok := folderAttachables[k]; !ok {
			util.Logger.Printf("FS: Media Scan: No attachables for Folder ID: %v", k)
			continue
		}

		switch x := v.(type) {
		case *db.Artist:
			artist := *x
			if folderAttachables[k].md != nil {
				md := *folderAttachables[k].md
				artist.MetadataID = md.ID
			}
			if folderAttachables[k].art != nil {
				art := *folderAttachables[k].art
				artist.ArtID = art.ID
			}
			if err := artist.Update(); err != nil {
				util.Logger.Printf(
					"FS: Media Scan: Error updating Artist %v: %v", artist.Title, err)
			}
		case *db.Album:	
			album := *x	
			if folderAttachables[k].md != nil {
				md := *folderAttachables[k].md
				album.MetadataID = md.ID
			}
			if folderAttachables[k].art != nil {
				art := *folderAttachables[k].art
				album.ArtID = art.ID
			}
			if err := album.Update(); err != nil {
				util.Logger.Printf(
					"FS: Media Scan: Error updating Album %v - %v: %v", 
					album.Artist, album.Title, err)
			}
		default:
			util.Logger.Printf(
				"FS: Media Scan: Unable to attach objects to unknown Type: %T", v)
		}
	}
}