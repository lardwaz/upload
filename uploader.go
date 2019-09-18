package upload // import "go.lsl.digital/lardwaz/upload"

// Uploader represents a file uploader (SMI)
type Uploader interface {
	// Upload accepts a filename, content and
	// returns a file disk path, file url path and error
	Upload(filename string, content []byte) (Uploaded, error)
}
