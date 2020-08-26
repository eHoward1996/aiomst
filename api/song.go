package api

import (
	"aiomst/db"
	"database/sql"
	"log"
	"strconv"

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
		log.Print(err)
		c.JSON(500, serverErr)
		return
	}

	resp := new(SongsResponse)
	resp.Songs = songs
	c.IndentedJSON(200, resp)
	return
}

func handleSongID(sID string, c *gin.Context)	{
	id, err := strconv.Atoi(sID)
	if err != nil {
		log.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	song := db.Song{ID: id}
	s, err := song.Load()
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(404, "Song ID not found.")
			return
		}

		log.Print(err)
		c.IndentedJSON(500, ErrGeneric)
		return
	}

	resp := new(SongsResponse)
	resp.Songs = []db.Song{s}
	c.IndentedJSON(200, resp)
	return
}