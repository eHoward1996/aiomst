package api

import (
	"strconv"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

// GetAlbums is the function called when a user accesses /albums on the frontend.
func GetAlbums(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
	albums, err := db.DB.AllAlbumsByTitle()
	if err != nil {
		util.Logger.Print(err)
		c.JSON(500, ErrGeneric)
		return
	}

	c.IndentedJSON(200, albums)
	return
}

func GetAlbum(c *gin.Context)	{
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")

	sID := c.Param("id")
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

	songs, err := db.DB.SongsForAlbum(album.ID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	type AlbumResponse struct{
		album db.Album   `json:"album"`
		songs []db.Song  `json:"songs"`
	}

	resp := AlbumResponse{
		album: album,
		songs: songs,
	}
	c.IndentedJSON(200, resp)
	return
}
