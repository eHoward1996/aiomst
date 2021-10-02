package api

import (
	"database/sql"
	"io"
	"strconv"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

func GetStream(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Accept-Ranges", "bytes")

	// Attempt to load the song with matching ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.Logger.Print(err)
		c.JSON(400, ErrGeneric)
		return
	}

	song := &db.Song{ID: id}
	if err := song.Load(); err != nil {
		// Check for invalid ID
		if err == sql.ErrNoRows {
			c.JSON(404, "song ID not found")
			return
		}

		// All other errors
		util.Logger.Print(err)
		c.JSON(500, serverErr)
		return
	}

	stream, err := song.Stream()
	if err != nil {
		util.Logger.Print(err)
		c.JSON(500, serverErr)
		return
	}

	rs, ok := toReadSeeker(stream)
	if !ok {
		util.Logger.Print("Unable to create readSeeker")
		c.JSON(500, serverErr)
		return
	}

	c.Stream(func(w io.Writer) bool {
		for {
			_, err := io.CopyN(w, rs, 2048)
			if err != nil && err != io.EOF {
				return false
			} else if err == io.EOF {
				return false
			}
			return true
		}
	})
}

func toReadSeeker(in io.Reader) (io.ReadSeeker, bool) {
	s, ok := in.(io.ReadSeeker)
	return s, ok
}
