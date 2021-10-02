package api

import (
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

func GetArtists(c *gin.Context)	{
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
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

	c.IndentedJSON(200, artists)
}

func GetArtist(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")

	sID := c.Param("id")
	util.Logger.Print(sID)
	id, err := strconv.Atoi(sID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(400, "Invalid integer artist ID")
		return
	}

	artist := db.Artist{ID: id}
	if err := artist.Load(); err != nil {
		util.Logger.Print(err)
		c.JSON(500, serverErr)
		return
	}

	albums, err := db.DB.AlbumsForArtist(artist.ID)
	if err != nil {
		util.Logger.Print(err)
		c.IndentedJSON(500, serverErr)
		return
	}
	
	songs, err := db.DB.SongsForArtist(artist.ID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}
	
	type ArtistResponse struct {
		Artist  	db.Artist  `json:"artist"`
		Albums		[]db.Album `json:"albums"`
		Songs 		[]db.Song  `json:"songs"`
	}

	resp := ArtistResponse{
		Artist: artist,
		Albums: albums,
		Songs: songs,
	}
	c.IndentedJSON(200, resp)
	return
}