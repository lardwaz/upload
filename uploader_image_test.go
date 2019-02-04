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
	testDataFolder = "testdata"
)

var update = flag.Bool("update", false, "update golden files")

type imageUploadTest struct {
	name                 string
	inputFile            string
	expectedFile         string
	expectedUploadError  bool
	expectedContentError bool
	uploader             *ImageUploader
}

type UploaderTestSuite struct {
	suite.Suite
	imageUploadTests []imageUploadTest
}

func (s *UploaderTestSuite) SetupSuite() {
	// Common upload configurations
	common := []Option{
		Dir(testDataFolder),
		Destination("tmp"),
		MediaPrefixURL("/"+testDataFolder+"/"),
		FileType(TypeImage),
	}
	commonJPEG := EvaluateOptions(append(common, ConvertTo(typeImageJPEG))...)
	commonPNG := EvaluateOptions(append(common, ConvertTo(typeImagePNG))...)

	// Test cases
	s.imageUploadTests = []imageUploadTest{
		{"Normal JPG", "normal.jpg", "normal_out.jpg", false, false, NewImageUploader(commonJPEG)},
		{"Normal PNG", "normal.png", "normal_out.png", false, false, NewImageUploader(commonPNG)},
		{"Transparent PNG", "transparent.png", "transparent_out.png", false, false, NewImageUploader(commonPNG)},
		{"Malformed JPG", "malformed.jpg", "malformed_out.jpg", false, false, NewImageUploader(commonJPEG)},
		{"Malformed PNG", "malformed.png", "malformed_out.png", false, false, NewImageUploader(commonPNG)},
		{"Damaged JPG", "damaged.jpg", "damaged_out.jpg", false, false, NewImageUploader(commonJPEG)},
		{"Damaged PNG", "damaged.png", "damaged_out.png", false, false, NewImageUploader(commonPNG)},
	}
}

func (s *UploaderTestSuite) TestImageUpload() {
	for _, tt := range s.imageUploadTests {
		s.Run(tt.name, func(){
			inputContent, err := ioutil.ReadFile(filepath.Join(testDataFolder, tt.inputFile))
			if err != nil {
				s.Failf("Cannot open input golden file", "%s: %v", tt.inputFile, err)
				return
			}

			uploaded, err := tt.uploader.Upload(tt.inputFile, inputContent)
			if tt.expectedUploadError && err != nil {
				// No problemo; we anticipated!
				return
			} else if err != nil {
				s.Failf("Cannot upload", "%s: %v", tt.inputFile, err)
				return
			}

			defer func() {
				// Cleanup
				if err = uploaded.Delete(); err != nil {
					s.Failf("Cannot delete uploaded file", "%s: %v", uploaded.DiskPath(), err)
				}
			}()

			content, err := ioutil.ReadFile(uploaded.DiskPath())
			if tt.expectedContentError && err != nil {
				// No problemo; we anticipated!
				return
			} else if err != nil {
				s.Failf("Cannot open uploaded file", "%s: %v", uploaded.DiskPath(), err)
				return
			}

			if *update {
				if err = ioutil.WriteFile(filepath.Join(testDataFolder, tt.expectedFile), content, 0644); err != nil {
					s.Failf("Cannot update golden file", "%s: %v", tt.expectedFile, err)
					return
				}
			}

			expectedContent, err := ioutil.ReadFile(filepath.Join(testDataFolder, tt.expectedFile))
			if err != nil {
				s.Failf("Cannot open output golden file", "%s: %v", tt.expectedFile, err)
				return
			}

			// Check if file content valid
			s.Equalf(expectedContent, content, "Uploaded content invalid")
		})
	}
}

func TestUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(UploaderTestSuite))
}
