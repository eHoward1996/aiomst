package api

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

type ArtistResponse struct {
	Artists  	[]db.Artist		`json:"artists"`
	Albums		[]db.Album 		`json:"albums"`
	Songs 		[]db.Song  		`json:"songs"`
}

func GetArtist(c *gin.Context)	{
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
	sID := c.Query("id")
	if sID == ""	{
		handleArtistNoID(c)
		return
	}
	handleArtistID(sID, c)
}

func handleArtistID(sID string, c *gin.Context)	{
	id, err := strconv.Atoi(sID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(400, "Invalid integer artist ID")
		return
	}

	artist := db.Artist{ID: id}
	if err := artist.Load(); err != nil {
		if err == sql.ErrNoRows 	{
			c.IndentedJSON(400, "Artist ID not found.")
			return
		}

		util.Logger.Print(err)
		c.JSON(500, serverErr)
		return
	}

	resp := new(ArtistResponse)
	resp.Artists = []db.Artist{artist}

	albums, err := db.DB.AlbumsForArtist(artist.ID)
	if err != nil {
		util.Logger.Print(err)
		c.IndentedJSON(500, serverErr)
		return
	}
	resp.Albums = albums
	
	songs, err := db.DB.SongsForArtist(artist.ID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}
	resp.Songs = songs
	c.IndentedJSON(200, resp)
	return
}

func handleArtistNoID(c *gin.Context)	{
	if limit := c.Query("limit"); limit != "" {
		var offset, count int
		if n, err := fmt.Sscanf(limit, "%d,%d", &offset, &count); n < 2 || err != nil {
			c.IndentedJSON(400, "Invalid comma-seperated integer pair.")
			return
		}

		artists, err := db.DB.LimitArtists(offset, count)
		if err != nil {
			util.Logger.Print(err)
			c.IndentedJSON(500, serverErr)
			return
		}

		c.IndentedJSON(200, artists)
		return
	}

	artists, err := db.DB.AllArtistsByTitle()
	if err != nil {
		util.Logger.Print(err)
		c.JSON(500, ErrGeneric)
		return
	}

	resp := new(ArtistResponse)
	resp.Artists = artists
	resp.Albums = []db.Album{}
	resp.Songs = []db.Song{}
	c.IndentedJSON(200, resp)
	return
}