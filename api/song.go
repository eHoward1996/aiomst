package api

import (
	"database/sql"
	"strconv"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

type SongsResponse struct {
	Songs []db.Song `json:"songs"`
}

func GetSongs(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
	sID := c.Query("id")
	if sID == ""	{
		handleSongNoID(c)
		return
	}
	handleSongID(sID, c)
	return
}

func handleSongNoID(c *gin.Context)	{
	songs, err := db.DB.AllSongsByTitle()
	if err != nil {
		util.Logger.Print(err)
		c.JSON(500, serverErr)
		return
	}

	resp := new(SongsResponse)
	resp.Songs = songs
	c.IndentedJSON(200, songs)
	return
}

func handleSongID(sID string, c *gin.Context)	{
	id, err := strconv.Atoi(sID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	song := db.Song{ID: id}
	if err := song.Load(); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(404, "Song ID not found.")
			return
		}

		util.Logger.Print(err)
		c.IndentedJSON(500, ErrGeneric)
		return
	}

	resp := new(SongsResponse)
	resp.Songs = []db.Song{song}
	c.IndentedJSON(200, resp)
	return
}