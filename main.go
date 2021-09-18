package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eHoward1996/aiomst/core"
	"github.com/eHoward1996/aiomst/util"
)

func main()	{
	util.InitializeLogger()
	util.Logger.Print("AIOMST: Starting")

	stat := util.ServerStatus()
	util.Logger.Printf(`AIOMST: Initial Server Status:
				  -- Hostname: %s
				  -- Platform: %s
				  -- Architecture: %s
				  -- # CPUs: %v
				  -- Memory Being Used (MB): %v
				  -- NumGoroutine: %v
				  -- Process PID: %v
				  -- Uptime: %d`,
		stat.Hostname, stat.Platform, stat.Architecture, stat.NumCPU, stat.MemoryMB,
		stat.NumGoroutine, stat.PID, stat.Uptime)

	killChan := make(chan struct {})
	exitChan := make(chan int)
	go core.TaskManager(killChan, exitChan)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func()	{
		for sig := range sigChan {
			util.Logger.Printf("AIOMST caught signal: %v... force halting", sig)
			killChan <- struct{}{}
			os.Exit(1)
		}
	}()
	
	code := <-exitChan
	util.Logger.Print("AIOMST GRACEFUL SHUTDOWN")
	os.Exit(code)
}
