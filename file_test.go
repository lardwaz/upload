package upload

import (
	"io/ioutil"
	"path"
	"path/filepath"
)

type mockUploadedFile struct {
	url      string
	diskPath string
	content  []byte
	options  Options
}

// NewMockUploadedFile returns a new mockUploadedFile (used for testing image processing so far)
func NewMockUploadedFile(name string, opts ...Option) *mockUploadedFile {
	options := EvaluateOptions(opts...)
	dirPath := path.Join(options.dir, options.destination)
	urlPath := path.Join(options.mediaPrefixURL, options.destination, name)
	diskPath := filepath.Join(dirPath, name)

	content, err := ioutil.ReadFile(diskPath)
	if err != nil {
		// Nothing too bad. We are mocking! ;)
	}

	return &mockUploadedFile{
		url:      urlPath,
		diskPath: diskPath,
		content:  content,
	}
}

func (m *mockUploadedFile) URLPath() string {
	return m.url
}

func (m *mockUploadedFile) DiskPath() string {
	return m.diskPath
}

func (m *mockUploadedFile) Content() []byte {
	return m.content
}

func (m *mockUploadedFile) Save(content []byte, overwrite bool) error {
	// Don't need an actual implementation
	return nil
}

func (m *mockUploadedFile) Delete() error {
	// Don't need an actual implementation
	return nil
}

func (m *mockUploadedFile) ChangeExt(string) error {
	// Don't need an actual implementation
	return nil
}
