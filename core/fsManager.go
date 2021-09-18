package core

import (
	// "errors"
	"log"
	"time"

	"github.com/eHoward1996/aiomst/fs"
	"github.com/eHoward1996/aiomst/fs/integrated"
	"github.com/eHoward1996/aiomst/util"

	"github.com/radovskyb/watcher"
)

// fsTaskQueue is a queue of tasks to be performed by the filesystem
var fsTaskQueue = make(chan fs.Task, 10)

// Track the number of file system events
var fsTaskCount = 0

// Initialize a queue to cancel filesystem tasks
var cancelQueue = make(chan chan struct{}, 10)

func fsManager(mediaPath string, sqlFile string, fsKillChan chan struct{})	{
	util.Logger.Print("FS MANAGER STARTED")

	// Initialize filesystem watcher
	watcherChan := make(chan struct{})

	// Queue an orphan scan
	o := new(fs.OrphanScan)
	o.SetFolders(mediaPath, "")
	o.Verbose(true)
	fsTaskQueue <- o

	// Queue a media scan
	m := new(fs.MediaScan)
	m.SetFolders(mediaPath, "")
	m.Verbose(true)
	fsTaskQueue <- m

	// Queue a job to integrate third party data
	go func(watcherChan chan struct{}) {
		<- watcherChan

		for {
			tpi := new(integrated.TPIntegrator)
			tpi.Integrate()
			time.Sleep(time.Hour * 24)
		}
	}(watcherChan)

	go handleFSTasks(watcherChan)
	go handleFSEvents(watcherChan)
	fsWatchKillSig(fsKillChan)
}

// Handle fs tasks in goroutine so they can be halted by the Task Manager
func handleFSTasks(watcherChan chan struct{}) {
	for {
		select {
		case task := <- fsTaskQueue:
			// Create a channel to halt the scan
			util.Logger.Printf("FS: Got new task (WhoAmI ==> %s)", task.WhoAmI())
			cancelChan := make(chan struct{})
			cancelQueue <- cancelChan

			// Retrieve the folder for the scan
			baseFolder, subFolder := task.Folders()
			
			changes, err := task.Scan(baseFolder, subFolder, cancelChan)
			if err != nil	{
				util.Logger.Printf("FS: Task Errored: %v", err)
			}

			if changes > 0 {
				util.UpdateScanTime()
				util.Logger.Printf("FS: New Scan Time: %v", util.ScanTimePretty())
			}

			cancelChan = <- cancelQueue
			close(cancelChan)
			fsTaskCount++

			if fsTaskCount == 2 {
				util.Logger.Print("FS: Finished initial media and orphan scans")
				close(watcherChan)
			}
		}
	}
}

// Handle fs events such as modify/rename/delete/create files.
func handleFSEvents(watcherChan chan struct{}) {
	<- watcherChan
	w := watcher.New()

	go func() 	{
		for {
			select {
			case ev := <- w.Event: 
				switch ev.Op.String() {
				case "MOVE":
					o := new(fs.OrphanScan)
					o.SetFolders("", ev.OldPath)
					o.Verbose(false)
					fsTaskQueue <- o
					fallthrough
				case "CREATE":
					m := new(fs.MediaScan)
					m.SetFolders(ev.Path, "")
					m.Verbose(false)
					fsTaskQueue <- m
				case "RENAME":
					fallthrough
				case "REMOVE":
					o := new(fs.OrphanScan)
					o.SetFolders("", ev.Path)
					o.Verbose(false)
					fsTaskQueue <- o
				}
			case err := <- w.Error:
				util.Logger.Print(err)
				return 
			}
		}
	}()

	// Watch media folder
	if err := w.AddRecursive(util.C.MediaFolderPath()); err != nil 	{
		log.Fatal(err)
	}
	if err := w.Start(1 * time.Minute); err != nil {
		log.Fatal(err)
	}
	util.Logger.Print("FS: Watching folder:", util.C.MediaFolderPath())
}

func fsWatchKillSig(fsKillChan chan struct{})	{
	for {		
		select {
		// Stop filesystem manager
		case <- fsKillChan:
			// Halt any in-progress tasks
			util.Logger.Print("FS: halting tasks")
			for i := 0; i < len(cancelQueue); i++ {
				// Receive a channel
				f := <-cancelQueue
				if f == nil {
					continue
				}

				// Send termination
				f <- struct{}{}
				util.Logger.Print("FS: task halted")
			}

			// Inform manager that shutdown is complete
			util.Logger.Print("FS MANAGER STOPPED!")
			fsKillChan <- struct{}{}
			return
		}
	}
}