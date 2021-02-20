package fs

import "github.com/eHoward1996/aiomst/db"

// Constant MBID Values
const errMBIDStartValue       = "errored"
const errNoClient             = "400:Client Not Available"
const errNoInfo               = "500:Info Not Available"
const errUnexpectedTermLength = "600:Unexpected Number of Terms"
const errConfidenceTooLow     = "200:Confidence Scores Too Low"

// The name of the default metadata file
const metadataFile = "metadata"

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

// AttachesArt is either a db.Album or db.Artist. Both implement methods
// to get an Art object by ID and set the objects ArtId. 
type AttachesArt interface {
	GetArt()	(*db.Art, error)
	SetArtID(int) error
}

// HasMetadata is either a db.Album or db.Artist. Both implement methods to get 
// and set Metadata (by ID).
type HasMetadata interface {
	GetMetadata() (*db.Metadata, error)
	SetMetadataID(int) error
}

// Task is an interface that defines a filesystem task
type Task interface {
	Folders() (string, string)
	SetFolders(string, string)
	Scan(string, string, chan struct{}) (int, error)
	Verbose(bool)
	WhoAmI() string
}