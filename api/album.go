package api

import (
	"strconv"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

// AlbumResponse is what is returned to the frontend.
type AlbumResponse struct {
	Albums []db.Album `json:"albums"`
	Songs  []db.Song  `json:"songs"`
}

// GetAlbums is the function called when a user accesses /albums on the frontend.
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
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	album := db.Album{ID: id}
	if err := album.Load(); err != nil {
		util.Logger.Print(err)
		c.JSON(500, ErrGeneric)
		return
	}

	resp := new(AlbumResponse)
	resp.Albums = []db.Album{album}

	songs, err := db.DB.SongsForAlbum(album.ID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	resp.Songs = songs
	c.IndentedJSON(200, resp)
	return
}

func handleAlbumNoID(c *gin.Context)	{
	albums, err := db.DB.AllAlbumsByTitle()
	if err != nil {
		util.Logger.Print(err)
		c.JSON(500, ErrGeneric)
		return
	}

	resp := new(AlbumResponse)
	resp.Albums = albums
	resp.Songs = []db.Song{}
	c.IndentedJSON(200, resp)
	return
}