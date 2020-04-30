package util

import (
	"os"
	"os/user"
	"runtime"
	"sync/atomic"
	"time"
)

// startTime represents the application's starting UNIX timestamp
var startTime = time.Now().Unix()

// scanTime is the last time the application did a media or orphan scan which
// created, modified, or deleted one or more items.  It defaults to the startup
// time, and is then updated by the filesystem manager.
var scanTime = startTime

// System is a grouped global which stores static information about the
// host operating system.
var System struct	{
	Hostname string
	User *user.User
}

// App is a global that stores static information about the App.
var App struct	{
	Name string
}

// osInfo represents basic, static information about the host operating system for this process
type osInfo struct {
	Architecture string
	Hostname     string
	NumCPU       int
	PID          int
	Platform     string
}

// Status represents information about the current process, including the basic, static
// information provided by osInfo
type Status struct {
	Architecture string  `json:"architecture"`
	Hostname     string  `json:"hostname"`
	MemoryMB     float64 `json:"memoryMb"`
	NumCPU       int     `json:"numCpu"`
	NumGoroutine int     `json:"numGoroutine"`
	PID          int     `json:"pid"`
	Platform     string  `json:"platform"`
	Uptime       int64   `json:"uptime"`
}

// init fetches information from the operating system and stores it in System
// for common access from different components of the service.
func init()	{
	setupSystem()
	setupApp()	
}

func setupSystem()	{
	hostname, err := os.Hostname()
	if err != nil	{
		panic(err)
	}
	System.Hostname = hostname

	user, err := user.Current()
	if err != nil 	{
		panic(err)
	}
	System.User = user
}

func setupApp()	{
	App.Name = "AIOMST"
}

// ScanTime returns the UNIX timestamp of the last time a media scan made changes
// to the database
func ScanTime() int64 {
	return atomic.LoadInt64(&scanTime)
}

// ScanTimePretty returns the scan time in a pretty date format
func ScanTimePretty()	time.Time {
	return time.Unix(ScanTime(), 0)
}

// UpdateScanTime updates the scanTime to the current UNIX timestamp
func UpdateScanTime() {
	atomic.StoreInt64(&scanTime, time.Now().Unix())
}

// OSInfo returns information about the host operating system for this process
func OSInfo() *osInfo {
	return &osInfo{
		Architecture: runtime.GOARCH,
		Hostname:     System.Hostname,
		NumCPU:       runtime.NumCPU(),
		PID:          os.Getpid(),
		Platform:     runtime.GOOS,
	}
}

// ServerStatus returns information about the current process status
func ServerStatus() *Status {
	// Retrieve basic OS information
	osStat := OSInfo()

	// Get current memory profile
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	// Report memory usage in MB
	memMB := float64((float64(mem.Alloc) / 1000) / 1000)

	// Get current uptime
	uptime := time.Now().Unix() - startTime

	// Return status
	return &Status{
		Architecture: osStat.Architecture,
		Hostname:     osStat.Hostname,
		MemoryMB:     memMB,
		NumCPU:       osStat.NumCPU,
		NumGoroutine: runtime.NumGoroutine(),
		PID:          osStat.PID,
		Platform:     osStat.Platform,
		Uptime:       uptime,
	}
}