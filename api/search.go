package api

import (
	"strings"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-gonic/gin"
)

const allTypes = "artists,albums,songs,folders"

// SearchResponse represents the JSON output for /api/search
type SearchResponse struct {
	Artists []db.Artist 	`json:"artists"`
	Albums  []db.Album  	`json:"albums"`
	Songs   []db.Song   	`json:"songs"`
	// Folders []db.Folder 	`json:"folders"`
}


func GetSearch(c *gin.Context)	{
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	
	query := c.Query("q")
	if query == ""	{
		c.IndentedJSON(400, "No search query specified.")
		return
	}

	resp := new(SearchResponse)
	types := c.Query("types")
	if types == ""	{
		types = allTypes
	}

	for _, t := range strings.Split(types, ",") 	{
		switch t {
		case "artists":
			artists, err := db.DB.SearchArtists(query)
			if err != nil {
				util.Logger.Print(err)
				c.IndentedJSON(500, serverErr)
				return
			}
			resp.Artists = artists
		case "albums":
			albums, err := db.DB.SearchAlbums(query)
			if err != nil {
				util.Logger.Print(err)
				c.IndentedJSON(500, serverErr)
				return
			}
			resp.Albums = albums
		case "songs":
			songs, err := db.DB.SearchSongs(query)
			if err != nil {
				util.Logger.Print(err)
				c.IndentedJSON(500, serverErr)
				return
			}
			resp.Songs = songs
		// case "folders":
		// 	folders, err := db.DB.SearchFolders(query)
		// 	if err != nil {
		// 		util.Logger.Print(err)
		// 		c.IndentedJSON(500, serverErr)
		// 		return
		// 	}
		// 	resp.Folders = folders
		}
	}

	c.IndentedJSON(200, resp)
	return
}