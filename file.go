package upload

// Uploaded represents the uploaded file
type Uploaded interface {
	URLPath() string
	DiskPath() string
	Content() []byte
	Save([]byte, bool) error
	Delete() error
	ChangeExt(string) error
}
