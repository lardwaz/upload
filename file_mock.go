package upload

import (
	"io/ioutil"
	"path"
	"path/filepath"

	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

type MockUploadedFile struct {
	url      string
	diskPath string
	content  []byte
	options  sdk.Options
}

// NewMockUploadedFile returns a new MockUploadedFile (used for testing image processing so far)
func NewMockUploadedFile(name string, opts ...func(sdk.Options)) *MockUploadedFile {
	options := option.EvaluateOptions(opts...)
	dirPath := path.Join(options.Dir(), options.Destination())
	urlPath := path.Join(options.MediaPrefixURL(), options.Destination(), name)
	diskPath := filepath.Join(dirPath, name)

	content, err := ioutil.ReadFile(diskPath)
	if err != nil {
		// Nothing too bad. We are mocking! ;)
	}

	return &MockUploadedFile{
		url:      urlPath,
		diskPath: diskPath,
		content:  content,
	}
}

func (m *MockUploadedFile) URLPath() string {
	return m.url
}

func (m *MockUploadedFile) DiskPath() string {
	return m.diskPath
}

func (m *MockUploadedFile) Content() []byte {
	return m.content
}

func (m *MockUploadedFile) Save(content []byte, overwrite bool) error {
	// Don't need an actual implementation
	return nil
}

func (m *MockUploadedFile) Delete() error {
	// Don't need an actual implementation
	return nil
}

func (m *MockUploadedFile) ChangeExt(string) error {
	// Don't need an actual implementation
	return nil
}
