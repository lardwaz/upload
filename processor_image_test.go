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
	inputFile 			 string
	expectedFile		 string
	expectedProcessError bool
	processor            *ImageProcessor
}

type ProcessorTestSuite struct {
	suite.Suite
	imageProcessTests []imageProcessTest
}

func (s *ProcessorTestSuite) SetupSuite() {
	// Set Watermark and backdrop assets
	WatermarkImage(filepath.Join(testDataFolder, "watermarks", "test-watermark.png"))
	BackdropImage(filepath.Join(testDataFolder, "backdrops", "test-backdrop.jpg"))

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		{"Normal", "normal.jpg", "processed_normal_out.jpg", false, NewImageProcessor()},
		{"Normal Thumb", "normal.jpg", "processed_normal_out.jpg", false, NewImageProcessor(Format("thumb", 200, 200, false))},
		{"Normal Height Zero", "normal.jpg", "processed_normal_out.jpg", false, NewImageProcessor(Format("hzero", 200, 0, false))},
		{"Normal Width Zero", "normal.jpg", "processed_normal_out.jpg", false, NewImageProcessor(Format("wzero", 0, 200, false))},
		{"Normal Upscale", "normal.jpg", "processed_normal_out.jpg", false, NewImageProcessor(Format("upscale", 500, 500, false))},
		{"Small Width", "normal.jpg", "processed_normal_out.jpg", true, NewImageProcessor(MinWidth(500))},
		{"Small Height", "normal.jpg", "processed_normal_out.jpg", true, NewImageProcessor(MinHeight(500))},
		{"Watermark", "normal.jpg", "watermarked_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Center), WatermarkVertical(Center)))},
		{"Backdrop Landscape", "normal.jpg", "backdropped_normal_out.jpg", false, NewImageProcessor(Format("back", 200, 200, true))},
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	// Common upload configurations
	commonOpts := EvaluateOptions(
		Dir(testDataFolder),
		MediaPrefixURL("/"+testDataFolder+"/"),
		FileType(TypeImage),
	)

	for _, tt := range s.imageProcessTests {
		uploadedFile := NewMockUploadedFile(tt.inputFile, *commonOpts)
		job, err := tt.processor.Process(uploadedFile, true)
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

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
