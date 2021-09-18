package core

import (
	"log"
	"os"

	"github.com/eHoward1996/aiomst/db"
	"github.com/eHoward1996/aiomst/util"
)

func dbManager(conf util.Config, dbLaunchChan, dbKillChan chan struct{})	{
	util.Logger.Print("DB MANAGER STARTED")

	if conf.Sqlite == nil {
		log.Fatalf("DB: Invalid database file")
	}

	path := conf.SqlFilePath()
	db.DB.DSN(path)
	util.Logger.Print("DB: SQLite:", db.DB.Path)

	// Setup the db
	if err := db.DB.Setup(); err != nil 	{
		log.Fatalf("DB: Could not set up database: %s", err)
	}

	// Verify DB file exists 
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("DB: Database file does not exist: %s", conf.Sqlite.File)
	}

	// Open the database connection
	if err := db.DB.Open(); err != nil	{
		log.Fatalf("DB: Could not open database: %s", err)
	}

	close(dbLaunchChan)
	dbWatchKillSig(dbKillChan)
}

func dbWatchKillSig(dbKillChan chan struct{})	{
	// Trigger events via channel
	for {
		select {
		// Stop database manager
		case <- dbKillChan:
			// Close the database connection
			if err := db.DB.Close(); err != nil {
				log.Fatalf("DB: Could not close connection")
			}

			// Inform manager that shutdown is complete
			util.Logger.Print("DB: Stopped!")
			dbKillChan <- struct{}{}
			return
		}
	}	
}