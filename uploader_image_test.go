package upload

// Basic imports
import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testFolder = "testdata"
)

var update = flag.Bool("update", false, "update golden files")

type imageUploadTest struct {
	name          string
	inputFile     string
	expectedFile  string
	expectedError bool
	uploader      ImageUpload
}

type UploaderTestSuite struct {
	suite.Suite
	imageUploadTests []imageUploadTest
}

func (s *UploaderTestSuite) SetupTest() {
	// Test cases
	s.imageUploadTests = []imageUploadTest{
		{"Normal JPG", "normal.jpg", "normal_out.jpg", false,
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
		{"Normal PNG", "normal.png", "normal_out.png", false,
			NewImageUploader(
				EvaluateOptions(
					Dir(testFolder),
					Destination("uploaded"),
					MediaPrefixURL("/testdata/"),
					FileType(TypeImage),
					ConvertTo(typeImagePNG),
				),
			),
		},
		{"Transparent PNG", "transparent.png", "transparent_out.png", false,
			NewImageUploader(
				EvaluateOptions(
					Dir(testFolder),
					Destination("uploaded"),
					MediaPrefixURL("/testdata/"),
					FileType(TypeImage),
					ConvertTo(typeImagePNG),
				),
			),
		},
	}
}

func (s *UploaderTestSuite) TestImageUpload() {
	for _, tt := range s.imageUploadTests {
		inputContent, err := ioutil.ReadFile(filepath.Join(testFolder, tt.inputFile))
		if err != nil {
			s.FailNowf("Cannot open input golden file", "Case: \"%s\". %s: %v", tt.name, tt.inputFile, err)
		}

		uploaded, err := tt.uploader.Upload(tt.inputFile, inputContent)
		if tt.expectedError && err != nil {
			// No problemo; we anticipated!
			return
		} else if err != nil {
			s.FailNowf("Cannot upload", "Case: \"%s\". %s: %v", tt.name, tt.inputFile, err)
		}

		content, err := ioutil.ReadFile(uploaded.DiskPath())
		if err != nil {
			s.FailNowf("Cannot open uploaded file", "Case: \"%s\". %s: %v", tt.name, uploaded.DiskPath(), err)
		}

		if *update {
			if err = ioutil.WriteFile(filepath.Join(testFolder, tt.expectedFile), content, 0644); err != nil {
				s.FailNowf("Cannot update golden file", "Case: \"%s\". %s: %v", tt.name, tt.expectedFile, err)
			}
		}

		expectedContent, err := ioutil.ReadFile(filepath.Join(testFolder, tt.expectedFile))
		if err != nil {
			s.FailNowf("Cannot open output golden file", "Case: \"%s\". %s: %v", tt.name, tt.expectedFile, err)
		}

		// Check if file content valid
		s.Equalf(expectedContent, content, "Case: \"%s\". Uploaded content invalid", tt.name)

		// Cleanup
		if err = uploaded.Delete(); err != nil {
			s.FailNowf("Cannot delete uploaded file", "Case: \"%s\". %s: %v", tt.name, uploaded.DiskPath(), err)
		}
	}
}

func TestUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(UploaderTestSuite))
}
