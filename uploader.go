package upload

// Uploader represents a file uploader (SMI)
type Uploader interface {
	// Upload accepts a filename, content and
	// returns a file disk path, file url path and error
	Upload(filename string, content []byte) (UploadedFile, error)
}
