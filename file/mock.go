package file

import (
	"io/ioutil"
	"path"
	"path/filepath"

	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

// MockGeneric is a mock implementation of Generic
type MockGeneric struct {
	url      string
	diskPath string
	content  []byte
	options  upload.Options
}

// NewMockGeneric returns a new MockGeneric (used for testing image processing so far)
func NewMockGeneric(name string, opts ...func(upload.Options)) *MockGeneric {
	options := option.EvaluateOptions(opts...)
	dirPath := path.Join(options.Dir(), options.Destination())
	urlPath := path.Join(options.MediaPrefixURL(), options.Destination(), name)
	diskPath := filepath.Join(dirPath, name)

	content, err := ioutil.ReadFile(diskPath)
	if err != nil {
		// Nothing too bad. We are mocking! ;)
	}

	return &MockGeneric{
		url:      urlPath,
		diskPath: diskPath,
		content:  content,
	}
}

// URLPath returns the URLPath
func (m *MockGeneric) URLPath() string {
	return m.url
}

// DiskPath returns the DiskPath
func (m *MockGeneric) DiskPath() string {
	return m.diskPath
}

// Content returns the Content
func (m *MockGeneric) Content() []byte {
	return m.content
}

// Save returns the Save
func (m *MockGeneric) Save(content []byte, overwrite bool) error {
	// Don't need an actual implementation
	return nil
}

// Delete returns the Delete
func (m *MockGeneric) Delete() error {
	// Don't need an actual implementation
	return nil
}

// ChangeExt returns the ChangeExt
func (m *MockGeneric) ChangeExt(string) error {
	// Don't need an actual implementation
	return nil
}
