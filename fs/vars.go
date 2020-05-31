package fs

// imgType is a set of valid file extensions for images
var imgType = map[string]bool	{
	".jpg"  : true,
	".jpeg" : true,
	".png"  : true,
}

// audioType is a set of valid file extensions for audio
var audioType = map[string]bool {
	".ape" : true,
	".flac": true,
	".m4a" : true,
	".mp3" : true,
	".mpc" : true,
	".ogg" : true,
	".wma" : true,
	".wv"  : true,
}

// AttachesArt is either a db.Album or db.Artist. Both implement a method
// to update a row in the database with an ArtID. 
type AttachesArt interface {
	GetArtID()    int
	SetArtID(int) error
}

// Task is an interface that defines a filesystem task
type Task interface {
	Folders() (string, string)
	SetFolders(string, string)
	Scan(string, string, chan struct{}) (int, error)
	Verbose(bool)
	WhoAmI() string
}