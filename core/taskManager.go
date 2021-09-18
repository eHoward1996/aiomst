package core

import (
	"github.com/eHoward1996/aiomst/util"
)

var dbKillChan, fsKillChan chan struct{}

// TaskManager begins AIOMST. It controls all sub-tasks.
func TaskManager(killChan chan struct{}, exitChan chan int)	{
	util.Logger.Print("TASK MANAGER STARTED")
	
	util.C = util.LoadConfig()
	config := util.C 
	
	dbLaunchFinishChan := make(chan struct{})
	dbKillChan         = make(chan struct{})
	go dbManager(config, dbLaunchFinishChan, dbKillChan)
	<- dbLaunchFinishChan

	fsKillChan := make(chan struct{})
	go fsManager(config.MediaFolderPath(), config.SqlFilePath(), fsKillChan)

	apiKillChan := make(chan struct{})
	go apiManager(apiKillChan)

	watchKillChans(killChan, exitChan)
}

func watchKillChans(killChan chan struct{}, exitChan chan int)	{
	for {
		select {
		case <- killChan:
			util.Logger.Print("TASK MANAGER: Triggering shutdown")

			fsKillChan <- struct{}{}
			<- fsKillChan
			close(fsKillChan)

			dbKillChan <- struct{}{}
			<- dbKillChan
			close(dbKillChan)

			util.Logger.Print("TASK MANAGER STOPPED")
			exitChan <- 0
		}
	}
}