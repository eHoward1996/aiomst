package core

import (
	"aiomst/util"
	"log"
)

var dbKillChan, fsKillChan chan struct{}

// TaskManager begins AIOMST. It controls all sub-tasks.
func TaskManager(killChan chan struct{}, exitChan chan int)	{
	log.Println("TASK MANAGER STARTED")
	
	util.C = util.LoadConfig()
	config := util.C 
	
	dbLaunchFinishChan := make(chan struct{})
	dbKillChan         = make(chan struct{})
	go dbManager(config, dbLaunchFinishChan, dbKillChan)
	<- dbLaunchFinishChan

	fsKillChan := make(chan struct{})
	fsLaunchFinishChan := make(chan struct{})
	go fsManager(config.MediaFolderPath(), fsLaunchFinishChan, fsKillChan)
	<- fsLaunchFinishChan 

	apiKillChan := make(chan struct{})
	go apiManager(apiKillChan)

	watchKillChans(killChan, exitChan)
}

func watchKillChans(killChan chan struct{}, exitChan chan int)	{
	for {
		select {
		case <- killChan:
			log.Print("TASK MANAGER: Triggering shutdown")

			fsKillChan <- struct{}{}
			<- fsKillChan
			close(fsKillChan)

			dbKillChan <- struct{}{}
			<- dbKillChan
			close(dbKillChan)

			log.Print("TASK MANAGER STOPPED")
			exitChan <- 0
		}
	}
}