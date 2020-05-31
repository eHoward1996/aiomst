package db

const dbSchema = `
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

/* folders */
CREATE TABLE IF NOT EXISTS "folders" (
	"id"        INTEGER PRIMARY KEY AUTOINCREMENT,
	"parent_id" INTEGER,
	"title"     TEXT,
	"path"      TEXT NOT NULL UNIQUE
);
CREATE UNIQUE INDEX "folders_unique_path" ON "folders" ("path");

/* artists */
CREATE TABLE IF NOT EXISTS "artists" (
	"id"    		INTEGER PRIMARY KEY AUTOINCREMENT,
	"art_id" 		INTEGER NOT NULL,
	"folder_id" INTEGER NOT NULL,
	"title" 		TEXT NOT NULL UNIQUE
);
CREATE UNIQUE INDEX "artists_unique_title" ON "artists" ("title");

/* albums */
CREATE TABLE IF NOT EXISTS "albums" (
	"id"        INTEGER PRIMARY KEY AUTOINCREMENT,
	"art_id"		INTEGER NOT NULL,
	"artist_id" INTEGER NOT NULL,
	"folder_id" INTEGER NOT NULL,
	"title"     TEXT NOT NULL,
	"year"      INTEGER NOT NULL
);
CREATE UNIQUE INDEX "albums_unique_artistId_title" ON "albums" ("artist_id", "title");

/* art */
CREATE TABLE IF NOT EXISTS "art" (
	"id"            INTEGER PRIMARY KEY AUTOINCREMENT,
	"file_size"     INTEGER NOT NULL,
	"file_name"     TEXT NOT NULL UNIQUE,
	"last_modified" INTEGER NOT NULL
);
CREATE UNIQUE INDEX "art_unique_fileName" ON "art" ("file_name");

/* songs */
CREATE TABLE IF NOT EXISTS "songs" (
	"id"            INTEGER PRIMARY KEY AUTOINCREMENT,
	"album_id"      INTEGER NOT NULL,
	"artist_id"     INTEGER NOT NULL,
	"bitrate"       INTEGER NOT NULL,
	"channels"      INTEGER NOT NULL,
	"comment"       TEXT,
	"file_name"     TEXT NOT NULL UNIQUE,
	"file_size"     INTEGER NOT NULL,
	"file_type_id"  INTEGER NOT NULL,
	"folder_id"     INTEGER NOT NULL,
	"genre"         TEXT,
	"last_modified" INTEGER NOT NULL,
	"length"        INTEGER NOT NULL,
	"sample_rate"   INTEGER NOT NULL,
	"title"         TEXT,
	"track"         INTEGER,
	"year"          INTEGER
);
CREATE UNIQUE INDEX "songs_unique_fileName" ON "songs" ("file_name");
COMMIT;`

// /* sessions */
// CREATE TABLE "sessions" (
// 	"id"      INTEGER PRIMARY KEY AUTOINCREMENT,
// 	"user_id" INTEGER NOT NULL,
// 	"client"  TEXT,
// 	"expire"  INTEGER NOT NULL,
// 	"key"     TEXT
// );
// CREATE UNIQUE INDEX "sessions_unique_key" ON "sessions" ("key");
// /* users */
// CREATE TABLE "users" (
// 	"id"           INTEGER PRIMARY KEY AUTOINCREMENT,
// 	"username"     TEXT,
// 	"password"     TEXT,
// 	"role_id"      INTEGER,
// 	"lastfm_token" TEXT
// );
// CREATE UNIQUE INDEX "users_unique_username" ON "users" ("username");


func getSchema() string {
	return dbSchema
}