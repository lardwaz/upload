package upload

// Basic imports
import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testFolder = "testdata"
)

type imageUploadTest struct {
	name         string
	inputFile    string
	expectedFile string
	uploader     ImageUpload
}

type UploaderTestSuite struct {
	suite.Suite
	imageUploadTests []imageUploadTest
}

func (s *UploaderTestSuite) SetupTest() {
	// Test cases
	s.imageUploadTests = []imageUploadTest{
		{"Normal JPG", "normal.jpg", "normal_out.jpg",
			NewImageUploader(
				EvaluateOptions(
					Dir(testFolder),
					Destination("uploaded"),
					MediaPrefixURL("/testdata/"),
					FileType(TypeImage),
					ConvertTo(typeImageJPEG),
				),
			),
		},
	}
}

func (s *UploaderTestSuite) TestUpload() {
	for _, tt := range s.imageUploadTests {
		inputContent, err := ioutil.ReadFile(filepath.Join(testFolder, tt.inputFile))
		if s.NotNil(err) {
			s.FailNowf("Cannot open input golden file", "%s: %v", tt.inputFile, err)
		}

		expectedContent, err := ioutil.ReadFile(filepath.Join(testFolder, tt.expectedFile))
		if s.NotNil(err) {
			s.FailNowf("Cannot open output golden file", "%s: %v", tt.expectedFile, err)
		}

		uploaded, err := tt.uploader.Upload(tt.inputFile, inputContent)
		if s.NotNil(err) {
			s.FailNowf("Cannot upload", "%s: %v", tt.inputFile, err)
		}

		content, err := ioutil.ReadFile(uploaded.DiskPath())
		if s.NotNil(err) {
			s.FailNowf("Cannot open uploaded file", "%s: %v", uploaded.DiskPath(), err)
		}

		// Check if file content valid
		s.Equal(expectedContent, content, "Uploaded content invalid")

		// Cleanup
		s.Errorf(uploaded.Delete(), "Cannot delete uploaded file %s: %v", uploaded.DiskPath(), err)
	}
}

func TestUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(UploaderTestSuite))
}
