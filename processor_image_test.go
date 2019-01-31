package upload

// Basic imports
import (
	"path/filepath"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type imageProcessTest struct {
	name                 string
	expectedFile		 string
	expectedProcessError bool
	processor            *ImageProcessor
}

type ProcessorTestSuite struct {
	suite.Suite
	uploadedFile 	  Uploaded
	imageProcessTests []imageProcessTest
}

func (s *ProcessorTestSuite) SetupSuite() {
	const (
		testImage = "normal.jpg"
	)

	inputContent, err := ioutil.ReadFile(filepath.Join(testDataFolder, testImage))
	if err != nil {
		s.FailNowf("Cannot open input golden file", "Setup suite: %v", err)
	}

	// Common upload configurations
	common := []Option{
		Dir(testDataFolder),
		Destination("tmp"),
		MediaPrefixURL("/"+testDataFolder+"/"),
		FileType(TypeImage),
	}

	commonOpts := EvaluateOptions(common...)
	uploader := NewImageUploader(commonOpts)

	s.uploadedFile, err = uploader.Upload(testImage, inputContent)
	if err != nil {
		s.FailNowf("Cannot upload", "Setup suite: %v", err)
	}

	// Set Watermark and backdrop assets
	WatermarkImage(filepath.Join(testDataFolder, "watermarks", "test-watermark.png"))
	BackdropImage(filepath.Join(testDataFolder, "backdrops", "test-backdrop.jpg"))

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		{"Normal", "processed_normal_out.jpg", false, NewImageProcessor()},
		{"Normal Thumb", "processed_normal_out.jpg", false, NewImageProcessor(Format("thumb", 200, 200, false))},
		{"Normal Height Zero", "processed_normal_out.jpg", false, NewImageProcessor(Format("hzero", 200, 0, false))},
		{"Normal Width Zero", "processed_normal_out.jpg", false, NewImageProcessor(Format("wzero", 0, 200, false))},
		{"Normal Upscale", "processed_normal_out.jpg", false, NewImageProcessor(Format("upscale", 500, 500, false))},
		{"Small Width", "processed_normal_out.jpg", true, NewImageProcessor(MinWidth(500))},
		{"Small Height", "processed_normal_out.jpg", true, NewImageProcessor(MinHeight(500))},
		{"Watermark", "watermarked_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Center), WatermarkVertical(Center)))},
		{"Backdrop Landscape", "backdropped_normal_out.jpg", false, NewImageProcessor(Format("back", 200, 200, true))},
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	for _, tt := range s.imageProcessTests {
		job, err := tt.processor.Process(s.uploadedFile, true)
		if tt.expectedProcessError && err != nil {
			// No problemo; we anticipated!
			continue
		} else if err != nil {
			s.Failf("Cannot process file", "Case: \"%s\": %v", tt.name, err)
			continue
		}

		select {
		case <-time.After(3 * time.Second):
			// We timed out!
			if !tt.expectedProcessError {
				s.Failf("Cannot process file", "Case: \"%s\": Timed out!", tt.name)
				continue
			}
		case <-job.Done:
			// Job done! We are good!
		}
		for _, format := range tt.processor.options.formats {
			fileDiskPath := job.File.DiskPath()+":"+format.name
			content, err := ioutil.ReadFile(fileDiskPath)
			if err != nil {
				s.Failf("Cannot open processed file", "Case: \"%s\". %s: %v", tt.name, fileDiskPath, err)
				continue
			}

			defer func(){
				// Cleanup
				if err = os.Remove(fileDiskPath); err != nil {
					// Not a problem!
				}
			}()
	
			expectedFileDiskPath := tt.expectedFile+":"+format.name
			if *update {
				if err = ioutil.WriteFile(filepath.Join(testDataFolder, expectedFileDiskPath), content, 0644); err != nil {
					s.Failf("Cannot update golden file", "Case: \"%s\". %s: %v", tt.name, expectedFileDiskPath, err)
					continue
				}
			}
	
			expectedContent, err := ioutil.ReadFile(filepath.Join(testDataFolder, expectedFileDiskPath))
			if err != nil {
				s.Failf("Cannot open output golden file", "Case: \"%s\". %s: %v", tt.name, expectedFileDiskPath, err)
				continue
			}
	
			// Check if file content valid
			s.Equalf(expectedContent, content, "Case: \"%s\". Uploaded content invalid", tt.name)
		}
	}
}

func (s *ProcessorTestSuite) TearDownSuite() {
	// Cleanup
	if err := s.uploadedFile.Delete(); err != nil {
		// Not a problem!
	}
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
