package api

import (
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"

	"github.com/eHoward1996/aiomst/db"

	"github.com/gin-gonic/gin"
)

// type ArtResponse struct {

// }

func GetArt(c *gin.Context) {
	sID := c.Query("id")
	if sID == "" {
		return
	}

	id, err := strconv.Atoi(sID)
	if err != nil {
		c.JSON(200, "Invalid integer art ID")
		return
	}

	art := &db.Art{ID: id}
	if err := art.Load(); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(404, "Art ID not found")
			return
		}

		c.IndentedJSON(500, serverErr)
		return
	}

	stream, err := art.Stream()
	if err != nil {
		c.IndentedJSON(404, err)
		return
	}
	
	img, imgFormat, err := image.Decode(stream)
	if err != nil {
		c.IndentedJSON(404, err)
		return
	}

	buffer := bytes.NewBuffer(nil)
	if imgFormat == "jpeg" {
		// JPEG, lossy encoding, default quality
		if err := jpeg.Encode(buffer, img, nil); err != nil {
			c.IndentedJSON(404, err)
			return
		}
		c.Header("Content-Type", "image/jpeg")
	} else {
		// Always send PNG as a backup
		// PNG, lossless encoding
		if err := png.Encode(buffer, img); err != nil {
			c.IndentedJSON(404, err)
			return
		}
		c.Header("Content-Type", "image/png")
	}

	// Serve content directly, account for range headers, and enabling caching.
	// http.ServeContent(w, r, art.FileName, time.Unix(art.LastModified, 0), bytes.NewReader(buffer.Bytes()))
	c.Data(200, c.ContentType(), buffer.Bytes())
	return
}