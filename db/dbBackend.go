package db

import (
	"os"
	"path"

	"github.com/eHoward1996/aiomst/util"

	"github.com/jmoiron/sqlx"

	// Include the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// SqlBackend represents a sqlite3 db backend
type SqlBackend struct {
	Path string
	db   *sqlx.DB
}

// DB is the database backend.
var DB SqlBackend

// DSN sets the Path for Sqlite
func (s *SqlBackend) DSN(path string)	{
	s.Path = util.ExpandHomeDir(path)
}

// Setup copies the sqlite db into the config into AIOMST directory
func (s *SqlBackend) Setup()	error {
	// Check for an existing config file
	_, err := os.Stat(s.Path)
	if err == nil	{
		// DB file exists
		return nil
	}

	// If error is something other than file not exists, return
	if !os.IsNotExist(err)	{
		return err
	}

	// Create a new DB file
	util.Logger.Print("DB: Creating new DB file: ", s.Path)
	dir := path.Dir(s.Path) + "/"
	file := path.Base(s.Path)

	// Make directory
	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}

	// Attempt to open destination
	dest, err := os.Create(dir + file)
	if err != nil {
		return err
	}

	initSchema := getSchema()
	database, err := sqlx.Connect("sqlite3", s.Path)
	if err != nil {
		panic(err)
	}

	util.Logger.Print("DB: Setup: Executing Init Schema")
	database.MustExec(initSchema)

	// Close file
	if err := dest.Close(); err != nil {
		return err
	}
	return nil
}

// Open initializes a new sqlite db connection
func (s *SqlBackend) Open() error	{
	sqlDB, err := sqlx.Open("sqlite3", s.Path)
	if err != nil 	{
		return err
	}

	// Do not wait for OS to respond to data write to disk
	if _, err := sqlDB.Exec("PRAGMA synchronous = OFF;"); err != nil {
		return err
	}

	// Keep rollback journal in memory, instead of on disk
	if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		return err
	}
	s.db = sqlDB
	return nil
}

// Close closes the sqlite sqlx database connection
func (s *SqlBackend) Close() error	{
	return s.db.Close()
}

// TruncateLog truncates the WAL file to 0 bytes.
func (s *SqlBackend) TruncateLog() error {
	if _, err := s.db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"); err != nil {
		return err
	}
	return nil
}