package upload

// Basic imports
import (
	"path/filepath"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/lsldigital/gocipe-upload/core"
)

type mockAssetBoxer struct{}

func NewMockAssetBoxer() *mockAssetBoxer {
	return &mockAssetBoxer{}
}

func (m *mockAssetBoxer) Open(name string) (*os.File, error) {
	return os.Open(name)
}

type imageProcessTest struct {
	name                 string
	prod 				 bool
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
	WatermarkImage(filepath.Join(testDataFolder, "watermarks", "test_watermark.png"))
	BackdropImage(filepath.Join(testDataFolder, "backdrops", "test_backdrop.jpg"))

	// Set asset box
	AssetBox(NewMockAssetBoxer())

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		{"Normal No Format", false, "normal.jpg", "noformat_normal_out.jpg", false, NewImageProcessor()},
		{"Normal No Format PNG", false, "normal.png", "noformat_normal_out.png", false, NewImageProcessor()},
		{"Normal Format", false, "normal.jpg", "format_normal_out.jpg", false, NewImageProcessor(Format("thumb", 200, 200, false))},
		{"Normal Format Negative Width & Height", false, "normal.jpg", "format_normal_out.jpg", false, NewImageProcessor(Format("neg", -1, -1, false))},
		{"PROD Normal Format", true, "normal.jpg", "format_prod_normal_out.jpg", false, NewImageProcessor(Format("thumb", 200, 200, false))},
		{"Normal Format PNG", false, "normal.png", "format_normal_out.png", false, NewImageProcessor(Format("thumb", 200, 200, false))},
		{"PROD Normal Format PNG", true, "normal.png", "format_prod_normal_out.png", false, NewImageProcessor(Format("thumb", 200, 200, false))},
		{"Normal Height Zero", false, "normal.jpg", "aspect_normal_out.jpg", false, NewImageProcessor(Format("hzero", 200, 0, false))},
		{"Normal Width Zero", false, "normal.jpg", "aspect_normal_out.jpg", false, NewImageProcessor(Format("wzero", 0, 200, false))},
		{"Normal Upscale", false, "normal.jpg", "upscale_normal_out.jpg", false, NewImageProcessor(Format("upscale", 500, 500, false))},
		{"Small Width", false, "normal.jpg", "min_normal_out.jpg", true, NewImageProcessor(MinWidth(500))},
		{"Small Height", false, "normal.jpg", "min_normal_out.jpg", true, NewImageProcessor(MinHeight(500))},
		{"Invalid File Type", false, "damaged.jpg", "invalid_normal_out.jpg", true, NewImageProcessor()},
		{"Invalid Image Type", false, "normal.gif", "invalid_normal_out.gif", true, NewImageProcessor()},
		{"Watermark Top Left", false, "normal.jpg", "watermarked_tl_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Left), WatermarkVertical(Top)))},
		{"Watermark Top Center", false, "normal.jpg", "watermarked_tc_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Center), WatermarkVertical(Top)))},
		{"Watermark Top Right", false, "normal.jpg", "watermarked_tr_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Right), WatermarkVertical(Top)))},
		{"Watermark Bottom Left", false, "normal.jpg", "watermarked_bl_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Left), WatermarkVertical(Bottom)))},
		{"Watermark Bottom Center", false, "normal.jpg", "watermarked_bc_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Center), WatermarkVertical(Bottom)))},
		{"Watermark Bottom Right", false, "normal.jpg", "watermarked_br_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Right), WatermarkVertical(Bottom)))},
		{"Watermark Center Left", false, "normal.jpg", "watermarked_cl_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Left), WatermarkVertical(Center)))},
		{"Watermark Center Center", false, "normal.jpg", "watermarked_cc_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Center), WatermarkVertical(Center)))},
		{"Watermark Center Right", false, "normal.jpg", "watermarked_cr_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(Right), WatermarkVertical(Center)))},
		{"Watermark Bad Pos", false, "normal.jpg", "watermarked_bad_prod_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(10), WatermarkVertical(10)))},
		{"PROD Watermark Bad Pos", true, "normal.jpg", "watermarked_bad_normal_out.jpg", false, NewImageProcessor(Format("water", 400, 400, false, WatermarkHorizontal(10), WatermarkVertical(10)))},
		{"Watermark Bad Pos", false, "normal.jpg", "watermarked_normal_out.jpg", false, NewImageProcessor(Format("damaged", 400, 400, false, WatermarkHorizontal(Center), WatermarkVertical(Center)))},
		{"Backdrop Landscape", false, "normal.jpg", "backdropped_normal_out.jpg", false, NewImageProcessor(Format("back", 200, 200, true))},
		{"PROD Backdrop Landscape", true, "normal.jpg", "backdropped_prod_normal_out.jpg", false, NewImageProcessor(Format("back", 200, 200, true))},
		{"Backdrop Portrait", false, "portrait.jpg", "backdropped_portrait_out.jpg", false, NewImageProcessor(Format("back", 200, 200, true))},
		{"PROD Backdrop Portrait", true, "portrait.jpg", "backdropped_prod_portrait_out.jpg", false, NewImageProcessor(Format("back", 200, 200, true))},
		{"Backdrop Damaged", false, "portrait.jpg", "backdropped_portrait_out.jpg", false, NewImageProcessor(Format("damaged", 200, 200, true))},
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
		oldEnv := core.Env
		// Adjust environment
		defer func(){
			core.Env = oldEnv
		}()
		if tt.prod {
			core.Env = core.EnvironmentPROD
		} else {
			core.Env = core.EnvironmentDEV
		}

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
