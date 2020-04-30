package core

import (
	"log"
)

func apiManager(apikillChan chan struct{})	{
	log.Print("API MANAGER STARTED")

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