package core

import (
	"aiomst/db"
	"aiomst/util"
	"log"
	"os"
)

func dbManager(conf util.Config, dbLaunchChan, dbKillChan chan struct{})	{
	log.Print("DB MANAGER STARTED")

	if conf.Sqlite == nil {
		log.Fatalf("DB: Invalid database file")
	}

	path := conf.SqlFilePath()
	db.DB.DSN(path)
	log.Println("DB: SQLite:", db.DB.Path)

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
			log.Println("DB: Stopped!")
			dbKillChan <- struct{}{}
			return
		}
	}	
}