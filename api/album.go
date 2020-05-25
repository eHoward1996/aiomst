package api

import (
	"aiomst/db"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlbumResponse struct {
	Album  db.Album   `json:"album"`
	Songs  []db.Song  `json:"songs"`
}

func GetAlbums(c *gin.Context)	{
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
	sID := c.Query("id")
	if sID == ""	{
		handleAlbumNoID(c)
		return
	}
	handleAlbumID(sID, c)
	return
}

func handleAlbumID(sID string, c *gin.Context)	{
	id, err := strconv.Atoi(sID)
	if err != nil {
		log.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	album := db.Album{ID: id}
	a, err := album.Load()
	if err != nil {
		log.Print(err)
		c.JSON(500, ErrGeneric)
		return
	}

	resp := new(AlbumResponse)
	resp.Album = a

	songs, err := db.DB.SongsForAlbum(album.ID)
	if err != nil {
		log.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	resp.Songs = songs
	c.IndentedJSON(200, resp)
	return
}

func handleAlbumNoID(c *gin.Context)	{
	albums, err := db.DB.AllAlbums()
	if err != nil {
		log.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}
	c.IndentedJSON(200, albums)
	return
}