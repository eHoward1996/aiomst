package util

import (
	"flag"
	"path"
)

var (
	hostPort = flag.String("host", ":8090", "The port to bind to.")
	driveLoc = flag.String("media", "~/MediaDrive/Media/Music", "The path to the media folder.")
	sqlDBLoc = flag.String("sqlite", "~/MediaDrive/Media/mediadb.db", "The sql db location.")
)

// C is the current configuration
var C Config

// SqliteFile is the configuration for the SQLite backend.
type SqliteFile struct {
	File string  `json:"file"`
}

// Config is the programs configuration options.
type Config struct {
	Host 				string 		  `json:"host"`
	MediaFolder string		  `json:"mediaFolder"`
	Sqlite			*SqliteFile	`json:"sqlite"`
}

// MediaFolderPath returns the fully expanded path to the media folder.
func (c Config) MediaFolderPath()	string {
	return path.Clean(ExpandHomeDir(c.MediaFolder))
}

// SqlFilePath return the fully expanded path to the SQL file.
func (c Config) SqlFilePath() string {
	return path.Clean(ExpandHomeDir(c.Sqlite.File))
}

// LoadConfig returns the configuration from the command line.
func LoadConfig() Config	{
	flag.Parse()

	return Config	{
		Host:        *hostPort,
		MediaFolder: *driveLoc,
		Sqlite: &SqliteFile {
			File: *sqlDBLoc,
		},
	}
}