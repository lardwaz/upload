package upload

// Basic imports
import (
	"path/filepath"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type imageProcessTest struct {
	name                 string
	outputFile			 string
	expectedProcessError bool
	processor            *ImageProcessor
}

type ProcessorTestSuite struct {
	suite.Suite
	uploadedFile 	  UploadedFile
	imageProcessTests []imageProcessTest
}

func (s *ProcessorTestSuite) SetupSuite() {
	const (
		testImage = "normal.jpg"
	)

	inputContent, err := ioutil.ReadFile(filepath.Join(testFolder, testImage))
	if err != nil {
		s.FailNowf("Cannot open input golden file", "Setup suite: %v", err)
	}

	// Common upload configurations
	common := []Option{
		Dir(testFolder),
		Destination("tmp"),
		MediaPrefixURL("/"+testFolder+"/"),
		FileType(TypeImage),
	}

	commonOpts := EvaluateOptions(common...)
	uploader := NewImageUploader(commonOpts)

	s.uploadedFile, err = uploader.Upload(testImage, inputContent)
	if err != nil {
		s.FailNowf("Cannot upload", "Setup suite: %v", err)
	}

	// TODO: Set Watermark and backdrop assets

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		{"Normal", "processed_normal_out.jpg", false, NewImageProcessor()},
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	for _, tt := range s.imageProcessTests {
		job, err := tt.processor.Process(s.uploadedFile, true)
		if tt.expectedProcessError && err != nil {
			// No problemo; we anticipated!
			return
		} else if err != nil {
			s.FailNowf("Cannot process file", "Case: \"%s\": %v", tt.name, err)
		}

		select {
		case <-time.After(3 * time.Second):
			// We timed out!
			if !tt.expectedProcessError {
				s.FailNowf("Cannot process file", "Case: \"%s\": Timed out!", tt.name)
			}
		case <-job.Done:
			// Job done! We are good!
		}

		
	}
}

func (s *ProcessorTestSuite) TearDownSuite() {
	// Cleanup
	if err := s.uploadedFile.Delete(); err != nil {
		s.Errorf(err, "Cannot delete uploaded file")
	}
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
