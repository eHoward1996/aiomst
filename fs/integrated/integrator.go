package integrated

import (
	"net/http"
	"sync"
	"time"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"
	"github.com/gammazero/workerpool"
	"github.com/irlndts/go-discogs"
	"github.com/pascoej/gomusicbrainz"
)

var mbClient *gomusicbrainz.WS2Client
var discogsClient *discogs.Discogs

type TPIntegrator struct{}

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
	util.Logger.Printf("FS: Third Party Integrator: Starting to retrieve...")
	startTime := time.Now()

	mbClient = getMusicBrainzClient()
	discogsClient = getDiscogsClient()
	i.makeRequests()

	util.Logger.Printf("FS: Third Party Integrator Complete [time: %s]",
		time.Since(startTime).String())
	if err := db.DB.TruncateLog(); err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Error truncating WAL file: %v", err)
	}
}

func (i TPIntegrator) makeRequests() {
	albums, err := db.DB.AlbumsWithErroredThirdPartyId()
	if err != nil {
		util.Logger.Printf(
			"FS: Third Party Integrator: Errored finding bad Album IDs: %v", err)
		return
	}
	util.Logger.Print(
		"FS: Third Party Integrator: Starting Albums Worker Pool...")
	
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
			recv.ArtistID = album.ArtistID
			recv.Artist = send.Artist
			recv.MetadataID = send.MetadataID

			util.Logger.Printf(
				"FS: Third Party Integrator: Requesting data for: %v - %v",
				send.Artist, send.Title)

			var wg sync.WaitGroup
			wg.Add(2)

			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				mbResp, err := i.getMusicBrainzData(send)
				if err != nil {
					util.Logger.Printf(
						"FS: Third Party Integrator: Error getting MusicBrainz content for " +
						"%v - %v: %v", send.Artist, send.Title, err)
				} else if mbResp.AlbumResponse != nil {
					recv.MBRelease = mbResp.AlbumResponse
				}	
			}(&wg)
			
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				dsResp, err := i.getDiscogsData(send)
				if err != nil {
					util.Logger.Printf(
						"FS: Third Party Integrator: Error getting Discogs content for " +
						"%v - %v: %v",
						send.Artist, send.Title, err)
				} else if dsResp.AlbumResponse != nil {
					recv.DiscogsMaster = dsResp.AlbumResponse
				}
			}(&wg)	
			
			wg.Wait()
			albumResults <- recv
		})
	}

	go func() {
		wp.StopWait()
		close(albumResults)
	}()

	count := 0
	for result := range albumResults {
		count++
		go func() {
			printReceivedAlbum(count, len(albums), result)
			i.createAlbumMetadata(result)
			i.createArtistMetadata(result)
		}()
		if err := updateAlbum(result); err != nil {
			util.Logger.Printf(
				"FS: Third Party Integrator: Error updating %v - %v: %v", 
				result.Artist, result.Album, err)
		}	
		if err := updateArtist(result); err != nil {
			util.Logger.Printf(
				"FS: Third Party Integrator: Error updating %v: %v", result.Artist, err)
		}
	}
}

