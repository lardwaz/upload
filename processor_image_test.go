package upload_test

// Basic imports
import (
	"path/filepath"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/lsldigital/gocipe-upload/core"
	"github.com/lsldigital/gocipe-upload"
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
	processor            *upload.ImageProcessor
}

type ProcessorTestSuite struct {
	suite.Suite
	imageProcessTests []imageProcessTest
}

func (s *ProcessorTestSuite) SetupSuite() {
	// Set Watermark and backdrop assets
	upload.WatermarkImage(filepath.Join(testDataFolder, "watermarks", "test_watermark.png"))
	upload.BackdropImage(filepath.Join(testDataFolder, "backdrops", "test_backdrop.jpg"))

	// Set asset box
	upload.AssetBox(NewMockAssetBoxer())

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		{"Normal No Format", false, "normal.jpg", "noformat_normal_out.jpg", false, upload.NewImageProcessor()},
		{"Normal No Format PNG", false, "normal.png", "noformat_normal_out.png", false, upload.NewImageProcessor()},
		{"Normal Format", false, "normal.jpg", "format_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("thumb", 200, 200, false))},
		{"Normal Format Negative Width & Height", false, "normal.jpg", "format_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("neg", -1, -1, false))},
		{"PROD Normal Format", true, "normal.jpg", "format_prod_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("thumb", 200, 200, false))},
		{"Normal Format PNG", false, "normal.png", "format_normal_out.png", false, upload.NewImageProcessor(upload.Formats("thumb", 200, 200, false))},
		{"PROD Normal Format PNG", true, "normal.png", "format_prod_normal_out.png", false, upload.NewImageProcessor(upload.Formats("thumb", 200, 200, false))},
		{"Normal Height Zero", false, "normal.jpg", "aspect_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("hzero", 200, 0, false))},
		{"Normal Width Zero", false, "normal.jpg", "aspect_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("wzero", 0, 200, false))},
		{"Normal Upscale", false, "normal.jpg", "upscale_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("upscale", 500, 500, false))},
		{"Small Width", false, "normal.jpg", "min_normal_out.jpg", true, upload.NewImageProcessor(upload.MinWidth(500))},
		{"Small Height", false, "normal.jpg", "min_normal_out.jpg", true, upload.NewImageProcessor(upload.MinHeight(500))},
		{"Invalid File Type", false, "damaged.jpg", "invalid_normal_out.jpg", true, upload.NewImageProcessor()},
		{"Invalid Image Type", false, "normal.gif", "invalid_normal_out.gif", true, upload.NewImageProcessor()},
		{"Watermark Top Left", false, "normal.jpg", "watermarked_tl_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Left), upload.WatermarkVertical(upload.Top)))},
		{"Watermark Top Center", false, "normal.jpg", "watermarked_tc_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Top)))},
		{"Watermark Top Right", false, "normal.jpg", "watermarked_tr_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Right), upload.WatermarkVertical(upload.Top)))},
		{"Watermark Bottom Left", false, "normal.jpg", "watermarked_bl_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Left), upload.WatermarkVertical(upload.Bottom)))},
		{"Watermark Bottom Center", false, "normal.jpg", "watermarked_bc_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Bottom)))},
		{"Watermark Bottom Right", false, "normal.jpg", "watermarked_br_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Right), upload.WatermarkVertical(upload.Bottom)))},
		{"Watermark Center Left", false, "normal.jpg", "watermarked_cl_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Left), upload.WatermarkVertical(upload.Center)))},
		{"Watermark Center Center", false, "normal.jpg", "watermarked_cc_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Center)))},
		{"Watermark Center Right", false, "normal.jpg", "watermarked_cr_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(upload.Right), upload.WatermarkVertical(upload.Center)))},
		{"Watermark Bad Pos", false, "normal.jpg", "watermarked_bad_prod_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(10), upload.WatermarkVertical(10)))},
		{"PROD Watermark Bad Pos", true, "normal.jpg", "watermarked_bad_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("water", 400, 400, false, upload.WatermarkHorizontal(10), upload.WatermarkVertical(10)))},
		{"Watermark Bad Pos", false, "normal.jpg", "watermarked_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("damaged", 400, 400, false, upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Center)))},
		{"Backdrop Landscape", false, "normal.jpg", "backdropped_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("back", 200, 200, true))},
		{"PROD Backdrop Landscape", true, "normal.jpg", "backdropped_prod_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats("back", 200, 200, true))},
		{"Backdrop Portrait", false, "portrait.jpg", "backdropped_portrait_out.jpg", false, upload.NewImageProcessor(upload.Formats("back", 200, 200, true))},
		{"PROD Backdrop Portrait", true, "portrait.jpg", "backdropped_prod_portrait_out.jpg", false, upload.NewImageProcessor(upload.Formats("back", 200, 200, true))},
		{"Backdrop Damaged", false, "portrait.jpg", "backdropped_portrait_out.jpg", false, upload.NewImageProcessor(upload.Formats("damaged", 200, 200, true))},
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	// Common upload configurations
	commonOpts := upload.EvaluateOptions(
		upload.Dir(testDataFolder),
		upload.MediaPrefixURL("/"+testDataFolder+"/"),
		upload.FileType(upload.TypeImage),
	)

	for _, tt := range s.imageProcessTests {
		s.Run(tt.name, func(){
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
	
			uploadedFile := upload.NewMockUploadedFile(tt.inputFile, *commonOpts)
			job, err := tt.processor.Process(uploadedFile, true)
			if tt.expectedProcessError && err != nil {
				// No problemo; we anticipated!
				return
			} else if err != nil {
				s.Failf("Cannot process file", "%v", err)
				return
			}
	
			select {
			case <-time.After(3 * time.Second):
				// We timed out!
				if !tt.expectedProcessError {
					s.Failf("Cannot process file", "Timed out!")
					return
				}
			case <-job.Done:
				// Job done! We are good!
			}
			for _, format := range tt.processor.Options().Formats() {
				fileDiskPath := job.File.DiskPath()+":"+format.Name()
				content, err := ioutil.ReadFile(fileDiskPath)
				if err != nil {
					s.Failf("Cannot open processed file", "%s: %v", fileDiskPath, err)
					return
				}
	
				defer func(){
					// Cleanup
					if err = os.Remove(fileDiskPath); err != nil {
						// Not a problem!
					}
				}()
		
				expectedFileDiskPath := tt.expectedFile+":"+format.Name()
				if *update {
					if err = ioutil.WriteFile(filepath.Join(testDataFolder, expectedFileDiskPath), content, 0644); err != nil {
						s.Failf("Cannot update golden file", "%s: %v", expectedFileDiskPath, err)
						continue
					}
				}
		
				expectedContent, err := ioutil.ReadFile(filepath.Join(testDataFolder, expectedFileDiskPath))
				if err != nil {
					s.Failf("Cannot open output golden file", "%s: %v", expectedFileDiskPath, err)
					continue
				}
		
				// Check if file content valid
				s.Equalf(expectedContent, content, "Uploaded content invalid")
			}
		})
	}
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
