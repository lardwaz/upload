package file

import (
	"io/ioutil"
	"path"
	"path/filepath"

	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

// MockUploaded is a mock implementation of Uploaded
type MockUploaded struct {
	url      string
	diskPath string
	content  []byte
	options  sdk.Options
}

// NewMockUploaded returns a new MockUploaded (used for testing image processing so far)
func NewMockUploaded(name string, opts ...func(sdk.Options)) *MockUploaded {
	options := option.EvaluateOptions(opts...)
	dirPath := path.Join(options.Dir(), options.Destination())
	urlPath := path.Join(options.MediaPrefixURL(), options.Destination(), name)
	diskPath := filepath.Join(dirPath, name)

	content, err := ioutil.ReadFile(diskPath)
	if err != nil {
		// Nothing too bad. We are mocking! ;)
	}

	return &MockUploaded{
		url:      urlPath,
		diskPath: diskPath,
		content:  content,
	}
}

// URLPath returns the URLPath
func (m *MockUploaded) URLPath() string {
	return m.url
}

// DiskPath returns the DiskPath
func (m *MockUploaded) DiskPath() string {
	return m.diskPath
}

// Content returns the Content
func (m *MockUploaded) Content() []byte {
	return m.content
}

// Save returns the Save
func (m *MockUploaded) Save(content []byte, overwrite bool) error {
	// Don't need an actual implementation
	return nil
}

// Delete returns the Delete
func (m *MockUploaded) Delete() error {
	// Don't need an actual implementation
	return nil
}

// ChangeExt returns the ChangeExt
func (m *MockUploaded) ChangeExt(string) error {
	// Don't need an actual implementation
	return nil
}
