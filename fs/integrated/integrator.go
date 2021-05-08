package integrated

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/eHoward1996/aiomst/db"
	"github.com/gammazero/workerpool"
	"github.com/irlndts/go-discogs"
	"github.com/pascoej/gomusicbrainz"
)

type TPIntegrator struct{
	mbClient      *gomusicbrainz.WS2Client
	discogsClient *discogs.Discogs
}

func getMusicBrainzClient() *gomusicbrainz.WS2Client {
	mb, _ := gomusicbrainz.NewWS2Client(
    "https://musicbrainz.org/ws/2",
    "All In One Media Server Tool",
    "0.0.1-beta",
		"http://github.com/eHoward1996/aiomst",
		new(http.Client),
	)
	return mb
}

func getDiscogsClient() *discogs.Discogs {
	d, _ := discogs.New(&discogs.Options{
		UserAgent: "All In One Media Server Tool",
		Token: "myPnYxaJxaTHrBHKNxzEEXuLDIviSIayFgmbmMiD",
	})
	return &d
}

func (i TPIntegrator) Integrate() {
	log.Printf("FS: Third Party Integrator: Starting to retrieve...")
	startTime := time.Now()

	i.mbClient = getMusicBrainzClient()
	i.discogsClient = getDiscogsClient()
	i.makeAlbumRequests()

	log.Printf("FS: Third Party Integrator Complete [time: %s]",
		time.Since(startTime).String())
}

func (i TPIntegrator) makeAlbumRequests() {
	albums, err := db.DB.AlbumsWithErroredThirdPartyId()
	if err != nil {
		log.Printf(
			"FS: Third Party Integrator: Errored finding bad Album IDs: %v", err)
		return
	}
	log.Println("FS: Third Party Integrator: Starting Albums Worker Pool...")
	
	albumResults := make(chan *recvAlbum)
	wp := workerpool.New(maxRoutines)	
	for _, album := range albums {
		album := album
		send := &sendAlbum{
			ID:         album.ID,
			MetadataID: album.MetadataID,
			Artist:     album.Artist,
			Title:      album.Title,
		}

		wp.Submit(func() {
			// For the time being, I'm just going to assume that if an album landed 
			// here that its better to just make all requests regardless of whether or
			// not all ids errored.
			// Right now, I can't think of a way to test whether or not a singular id
			// was "errored" (from not being able to be found or from initialization)
			// or was found to be non-unique. So, as stated above, if the album was 
			// returned from the SQL query, I'm going to run through this entire 
			// process.
			// TODO: find a more efficient way to do this.

			recv := new(recvAlbum)
			recv.AlbumID = send.ID
			recv.Album = send.Title
			recv.Artist = send.Artist
			recv.MetadataID = send.MetadataID

			log.Printf(
				"FS: Third Party Integrator: Requesting data for: %v - %v",
				send.Artist, send.Title)

			mbResp, err := i.getMusicBrainzData(send)
			if err != nil {
				log.Printf(
					"FS: Third Party Integrator: Error getting MusicBrainz content for " +
					"%v - %v: %v", send.Artist, send.Title, err)
			} else if mbResp.AlbumResponse != nil {
				recv.MBRelease = mbResp.AlbumResponse
			}	
				
			dsResp, err := i.getDiscogsData(send)
			if err != nil {
				log.Printf(
					"FS: Third Party Integrator: Error getting Discogs content for " +
					"%v - %v: %v",
					send.Artist, send.Title, err)
			} else if dsResp.AlbumResponse != nil {
				recv.DiscogsMaster = dsResp.AlbumResponse
			}
			
			albumResults <- recv
			// time.Sleep(2 * time.Second)
		})
	}

	go func() {
		wp.StopWait()
		close(albumResults)
	}()

	count := 0
	for result := range albumResults {
		count++
		go createAlbumMetadata(result)
		if err := updateAlbum(result); err != nil {
			log.Printf("FS: Third Party Integrator: Error updating %v - %v: %v", 
				result.Artist, result.Album, err)
		}
		printReceivedAlbum(count, len(albums), result)
	}
}

func printReceivedAlbum(count, numAlbums int, result *recvAlbum) {	
	mbId := "errored"
	dsId := -1
	if result.MBRelease != nil {
		mbId = string(result.MBRelease.ID)
	}
	if result.DiscogsMaster != nil {
		dsId = result.DiscogsMaster.ID
	}

	artistString := formatOutput(result.Artist, 18)
	albumString := formatOutput(result.Album, 18)
	processedStr := fmt.Sprintf(
		"(%v of %v) %v - %v MBID: %v\tDiscogsID: %v", count, numAlbums,
		artistString, albumString, formatOutput(mbId, 13), dsId)
	log.Printf("%v", processedStr)
}

func updateAlbum(recv *recvAlbum) error {
	album := &db.Album{ID: recv.AlbumID}
	if err := db.DB.LoadAlbum(album); err != nil {
		return err
	}

	if recv.MBRelease != nil {
		album.MBID = string(recv.MBRelease.ID)
	}
	if recv.DiscogsMaster != nil {
		album.DiscogsID = recv.DiscogsMaster.ID
	}
	
	return album.Update()
}

func createAlbumMetadata(recv *recvAlbum) {
	albumMD := db.AlbumMetadata{
		AlbumName: recv.Album,
		MusicBrainz: buildMBAlbumMetadata(recv.MBRelease),
		Discogs: buildDiscogsAlbumMetadata(recv.DiscogsMaster),
	}

	writeMetadata(recv.MetadataID, albumMD)
}

func writeMetadata(metadataID int, md interface{}) {
	bytes, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		log.Printf(
			"FS: Third Party Integrator: Error writing Metadata file: %v", err)
		return
	}

	tmp := new(db.Metadata)
	tmp.ID = metadataID
	if err := tmp.Load(); err != nil {
		log.Printf(
			"FS: Third Party Integrator: Error loading Metadata object: %v", err)
		return
	}
	if err := ioutil.WriteFile(tmp.Path, bytes, 0644); err != nil {
		log.Printf(
			"FS: Third Party Integrator: Error writing Metadata File: %v", err)
		return
	}
}
