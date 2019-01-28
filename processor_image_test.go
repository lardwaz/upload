package upload

// Basic imports
import (
	"filepath"
	"log"
	"io/ioutil"
	"testing"

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
	// Common upload configurations
	common := []Option{
		Dir(testFolder),
		Destination("uploaded"),
		MediaPrefixURL("/testdata/"),
		FileType(TypeImage),
	}

	uploader := NewImageUploader(EvaluateOptions(common...))

	inputContent, err := ioutil.ReadFile(filepath.Join(testFolder, "normal.jpg"))
	if err != nil {
		s.FailNowf("Cannot open input golden file", "Setup suite: %v", err)
	}

	s.uploadedFile, err := uploader.Upload("normal.jpg", inputContent)
	if err != nil {
		s.FailNowf("Cannot upload", "Setup suite: %v", err)
	}

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		// TODO: input test cases
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	for _, tt := range s.imageProcessTests {
		log.Println(tt)
	}
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
