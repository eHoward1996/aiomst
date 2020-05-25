package core

import (
	"aiomst/api"
	"aiomst/util"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	// "github.com/mdlayher/wavepipe/api"
)

func apiManager(apikillChan chan struct{})	{
	log.Print("API MANAGER STARTED")
	gracefulChan := make(chan struct{}, 0)

	// Initialize Gin (api router)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	
	r.Use(func(c *gin.Context)  {
		srv := fmt.Sprintf("AIOMST ==> (%s_%s)", runtime.GOOS, runtime.GOARCH)
		c.Set("Server", srv)
		c.Next()
		return
	})

	r.GET("/", func(c *gin.Context) {c.String(http.StatusOK, "pong")})
	r.GET("/albums", api.GetAlbums)
	r.GET("/artists", api.GetArtist)
	r.GET("/songs", api.GetSongs)
	r.GET("/search", api.GetSearch)
	sConf := util.LoadConfig()
	server := &http.Server{
		Addr:    sConf.Host,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API: Encountered Error: %s", err)
		}
		close(gracefulChan)
	}()

	watchKillSig(apikillChan)
}

func watchKillSig(apiKillChan chan struct{})	{
	for {
		select {
		// Stop API
		case <-apiKillChan:
			// Inform manager that shutdown is complete
			log.Println("API MANAGER STOPPED")
			apiKillChan <- struct{}{}
			return
		}
	}
}