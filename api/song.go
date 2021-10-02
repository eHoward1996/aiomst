package api

import (
	"fmt"
	"strconv"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

func GetSongs(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
	if limit := c.Query("limit"); limit != "" {
		var offset, count int
		if n, err := fmt.Sscanf(limit, "%d,%d", &offset, &count); n < 2 || err != nil {
			c.IndentedJSON(400, "Invalid comma-seperated integer pair.")
			return
		}

		songs, err := db.DB.LimitSongs(offset, count)
		if err != nil {
			util.Logger.Print(err)
			c.IndentedJSON(500, serverErr)
			return
		}

		c.IndentedJSON(200, songs)
		return
	}

	songs, err := db.DB.AllSongsByTitle()
	if err != nil {
		util.Logger.Print(err)
		c.JSON(500, serverErr)
		return
	}

	c.IndentedJSON(200, songs)
}

func GetSong(c *gin.Context) {
	sID := c.Param("id")
	id, err := strconv.Atoi(sID)
	if err != nil {
		util.Logger.Print(err)
		c.JSON(200, ErrGeneric)
		return
	}

	song := db.Song{ID: id}
	if err := song.Load(); err != nil {
		util.Logger.Print(err)
		c.IndentedJSON(500, ErrGeneric)
		return
	}

	c.IndentedJSON(200, song)
}