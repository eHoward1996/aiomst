package core

import (
	"aiomst/api"
	"aiomst/db"
	"aiomst/util"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func apiManager(apikillChan chan struct{})	{
	log.Print("API MANAGER STARTED")
	gracefulChan := make(chan struct{}, 0)

	gin.SetMode(gin.ReleaseMode)
	// Initialize Gin (api router)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Type", "Content-Length", "origin"},
		// AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	
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

	imgFiles, err := db.DB.AllArt()
	if err != nil || len(imgFiles) == 0 {
		log.Printf("API: Couldn't get Art files: %s", err)
	}

	imgRoute := r.Group("/art")
	{
		for _, art := range imgFiles {
			id := strconv.Itoa(art.ID)
			imgRoute.StaticFile(id, art.FileName)
		}
	}

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