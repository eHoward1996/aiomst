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
	"id"    		       INTEGER PRIMARY KEY AUTOINCREMENT,
	"mb_id"						 TEXT NOT NULL,
	"discogs_id"			 TEXT NOT NULL,
	"metadata_id"      INTEGER,
	"art_id" 		       INTEGER NOT NULL,
	"folder_id"        INTEGER NOT NULL,
	"title" 		       TEXT NOT NULL UNIQUE,
	"normalized_title" TEXT NOT NULL UNIQUE
);
CREATE UNIQUE INDEX "artists_unique_title" ON "artists" ("title");

/* albums */
CREATE TABLE IF NOT EXISTS "albums" (
	"id"               INTEGER PRIMARY KEY AUTOINCREMENT,
	"mb_id"						 TEXT NOT NULL,
	"discogs_id"			 INTEGER NOT NULL,
	"metadata_id"      INTEGER,
	"art_id"		       INTEGER NOT NULL,
	"artist_id"        INTEGER NOT NULL,
	"folder_id"        INTEGER NOT NULL,
	"title"            TEXT NOT NULL,
	"normalized_title" TEXT NOT NULL,
	"year"             INTEGER NOT NULL
);
CREATE UNIQUE INDEX "albums_unique_artist_id_title" ON "albums" ("artist_id", "title");

/* art */
CREATE TABLE IF NOT EXISTS "art" (
	"id"            INTEGER PRIMARY KEY AUTOINCREMENT,
	"file_size"     INTEGER NOT NULL,
	"path"          TEXT NOT NULL UNIQUE,
	"last_modified" INTEGER NOT NULL
);
CREATE UNIQUE INDEX "art_unique_path" ON "art" ("path");

/* songs */
CREATE TABLE IF NOT EXISTS "songs" (
	"id"                INTEGER PRIMARY KEY AUTOINCREMENT,
	"mb_id"						  TEXT NOT NULL,
	"album_id"          INTEGER NOT NULL,
	"artist_id"         INTEGER NOT NULL,
	"bitrate"           INTEGER NOT NULL,
	"channels"          INTEGER NOT NULL,
	"comment"           TEXT,
	"path"              TEXT NOT NULL UNIQUE,
	"file_size"         INTEGER NOT NULL,
	"file_type_id"      INTEGER NOT NULL,
	"folder_id"         INTEGER NOT NULL,
	"genre"             TEXT,
	"last_modified"     INTEGER NOT NULL,
	"length"            INTEGER NOT NULL,
	"sample_rate"       INTEGER NOT NULL,
	"title"             TEXT NOT NULL,
	"normalized_title"  TEXT NOT NULL,
	"track"             INTEGER,
	"year"              INTEGER
);
CREATE UNIQUE INDEX "songs_unique_path" ON "songs" ("path");

/* metadata */
CREATE TABLE IF NOT EXISTS "metadata" (
	"id"            INTEGER PRIMARY KEY AUTOINCREMENT,
	"folder_id"     INTEGER NOT NULL,
	"file_size"     INTEGER NOT NULL,
	"last_modified" INTEGER NOT NULL,
	"path"          TEXT NOT NULL UNIQUE
);
CREATE UNIQUE INDEX "metadata_unique_path" ON "metadata" ("path");
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