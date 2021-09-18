package core

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/eHoward1996/aiomst/api"
	"github.com/eHoward1996/aiomst/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func apiManager(apikillChan chan struct{})	{
	util.Logger.Print("API MANAGER STARTED")
	gracefulChan := make(chan struct{}, 0)

	gin.SetMode(gin.ReleaseMode)
	// Initialize Gin (api router)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:     	true,
		AllowHeaders:     	 	[]string{
			"Content-Type",
			"Origin",
			"Accept-Ranges",
			"Ranges",
			"Access-Control-Allow-Origin",
		},
		ExposeHeaders:    []string{
			"Content-Type", 
			"Content-Length", 
			"Origin",
		},
		MaxAge: 12 * time.Hour,
	}))
	
	r.Use(func(c *gin.Context)  {
		srv := fmt.Sprintf("AIOMST ==> (%s_%s)", runtime.GOOS, runtime.GOARCH)
		c.Set("Server", srv)
		c.Next()
		return
	})
	
	// r.GET("/", func(c *gin.Context) {c.String(http.StatusOK, "pong")})
	r.StaticFile("/", "./core/public.html")
	r.GET("/albums",  api.GetAlbums)
	r.GET("/artists", api.GetArtist)
	r.GET("/songs",   api.GetSongs)
	r.GET("/search",  api.GetSearch)
	r.GET("/art",     api.GetArt)
	r.GET("/stream",  api.GetStream)

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
			util.Logger.Print("API MANAGER STOPPED")
			apiKillChan <- struct{}{}
			return
		}
	}
}