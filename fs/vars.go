package fs

// Constant MBID Values
const errStartValueString = "errored"
const errStartValueInt = -1

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

// HasAttachables is any type that implements the HasAttachables method. 
// Right now, that's just Artists and Albums.
// TODO: Find a better way to do this.
type HasAttachables interface {
	HasAttachables()
}

// Task is an interface that defines a filesystem task
type Task interface {
	Folders() (string, string)
	SetFolders(string, string)
	Scan(string, string, chan struct{}) (int, error)
	Verbose(bool)
	WhoAmI() string
}