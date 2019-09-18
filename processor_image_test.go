package upload_test

// Basic imports
import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/core"
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
	prod                 bool
	inputFile            string
	expectedFile         string
	expectedProcessError bool
	processor            upload.ImageProcessor
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
		{"Normal Format", false, "normal.jpg", "format_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("thumb"), upload.FormatWidth(200), upload.FormatHeight(200)))},
		{"Normal Format Negative Width & Height", false, "normal.jpg", "format_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("neg"), upload.FormatWidth(-1), upload.FormatHeight(-1)))},
		{"PROD Normal Format", true, "normal.jpg", "format_prod_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("thumb"), upload.FormatWidth(200), upload.FormatHeight(200)))},
		{"Normal Format PNG", false, "normal.png", "format_normal_out.png", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("thumb"), upload.FormatWidth(200), upload.FormatHeight(200)))},
		{"PROD Normal Format PNG", true, "normal.png", "format_prod_normal_out.png", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("thumb"), upload.FormatWidth(200), upload.FormatHeight(200)))},
		{"Normal Height Zero", false, "normal.jpg", "aspect_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("hzero"), upload.FormatWidth(200)))},
		{"Normal Width Zero", false, "normal.jpg", "aspect_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("wzero"), upload.FormatHeight(200)))},
		{"Normal Upscale", false, "normal.jpg", "upscale_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("upscale"), upload.FormatWidth(500), upload.FormatHeight(500)))},
		{"Small Width", false, "normal.jpg", "min_normal_out.jpg", true, upload.NewImageProcessor(upload.MinWidth(500))},
		{"Small Height", false, "normal.jpg", "min_normal_out.jpg", true, upload.NewImageProcessor(upload.MinHeight(500))},
		{"Invalid File Type", false, "damaged.jpg", "invalid_normal_out.jpg", true, upload.NewImageProcessor()},
		{"Invalid Image Type", false, "normal.gif", "invalid_normal_out.gif", true, upload.NewImageProcessor()},
		{"Watermark Top Left", false, "normal.jpg", "watermarked_tl_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Left), upload.WatermarkVertical(upload.Top))))},
		{"Watermark Top Center", false, "normal.jpg", "watermarked_tc_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Top))))},
		{"Watermark Top Right", false, "normal.jpg", "watermarked_tr_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Right), upload.WatermarkVertical(upload.Top))))},
		{"Watermark Bottom Left", false, "normal.jpg", "watermarked_bl_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Left), upload.WatermarkVertical(upload.Bottom))))},
		{"Watermark Bottom Center", false, "normal.jpg", "watermarked_bc_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Bottom))))},
		{"Watermark Bottom Right", false, "normal.jpg", "watermarked_br_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Right), upload.WatermarkVertical(upload.Bottom))))},
		{"Watermark Center Left", false, "normal.jpg", "watermarked_cl_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Left), upload.WatermarkVertical(upload.Center))))},
		{"Watermark Center Center", false, "normal.jpg", "watermarked_cc_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Center))))},
		{"Watermark Center Right", false, "normal.jpg", "watermarked_cr_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Right), upload.WatermarkVertical(upload.Center))))},
		{"Watermark Bad Pos", false, "normal.jpg", "watermarked_bad_prod_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(10), upload.WatermarkVertical(10))))},
		{"PROD Watermark Bad Pos", true, "normal.jpg", "watermarked_bad_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("water"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(10), upload.WatermarkVertical(10))))},
		{"Watermark Bad Pos", false, "normal.jpg", "watermarked_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("damaged"), upload.FormatWidth(400), upload.FormatHeight(400), upload.FormatWatermark(upload.WatermarkHorizontal(upload.Center), upload.WatermarkVertical(upload.Center))))},
		{"Backdrop Landscape", false, "normal.jpg", "backdropped_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("back"), upload.FormatWidth(200), upload.FormatHeight(200), upload.FormatBackdrop(true)))},
		{"PROD Backdrop Landscape", true, "normal.jpg", "backdropped_prod_normal_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("back"), upload.FormatWidth(200), upload.FormatHeight(200), upload.FormatBackdrop(true)))},
		{"Backdrop Portrait", false, "portrait.jpg", "backdropped_portrait_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("back"), upload.FormatWidth(200), upload.FormatHeight(200), upload.FormatBackdrop(true)))},
		{"PROD Backdrop Portrait", true, "portrait.jpg", "backdropped_prod_portrait_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("back"), upload.FormatWidth(200), upload.FormatHeight(200), upload.FormatBackdrop(true)))},
		{"Backdrop Damaged", false, "portrait.jpg", "backdropped_portrait_out.jpg", false, upload.NewImageProcessor(upload.Formats(upload.FormatName("damaged"), upload.FormatWidth(200), upload.FormatHeight(200), upload.FormatBackdrop(true)))},
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	// Common upload configurations
	commonOpts := []func(upload.Options){
		upload.Dir(testDataFolder),
		upload.MediaPrefixURL("/" + testDataFolder + "/"),
		upload.FileType(upload.TypeJPEG),
		upload.FileType(upload.TypeJPEG2),
		upload.FileType(upload.TypePNG),
		upload.FileType(upload.TypeGIF),
		upload.FileType(upload.TypeHEIF),
	}

	for _, tt := range s.imageProcessTests {
		s.Run(tt.name, func() {
			oldEnv := core.Env
			// Adjust environment
			defer func() {
				core.Env = oldEnv
			}()
			if tt.prod {
				core.Env = core.EnvironmentPROD
			} else {
				core.Env = core.EnvironmentDEV
			}

			uploadedFile := upload.NewMockUploadedFile(tt.inputFile, commonOpts...)
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
				// Job timed out! Did we expect?
				if !tt.expectedProcessError {
					s.Failf("Cannot process file", "%s: Timed out!", job.File().DiskPath())
					return
				}
			case <-job.Done():
			// Job done! We are good!

			case err = <-job.Failed():
				// Job failed! Did we expect?
				if !tt.expectedProcessError {
					s.Failf("Cannot process file", "%s: %v", job.File().DiskPath(), err)
					return
				}
			}

			formats := tt.processor.Options().Formats()

			formats.Each(func(name string, format upload.OptionsFormat) {
				fileDiskPath := job.File().DiskPath() + "-" + format.Name()
				content, err := ioutil.ReadFile(fileDiskPath)
				if err != nil {
					s.Failf("Cannot open processed file", "%s: %v", fileDiskPath, err)
					return
				}

				defer func() {
					// Cleanup
					if err = os.Remove(fileDiskPath); err != nil {
						// Not a problem!
					}
				}()

				expectedFileDiskPath := tt.expectedFile + "-" + format.Name()
				if *update {
					if err = ioutil.WriteFile(filepath.Join(testDataFolder, expectedFileDiskPath), content, 0644); err != nil {
						s.Failf("Cannot update golden file", "%s: %v", expectedFileDiskPath, err)
						return
					}
				}

				expectedContent, err := ioutil.ReadFile(filepath.Join(testDataFolder, expectedFileDiskPath))
				if err != nil {
					s.Failf("Cannot open output golden file", "%s: %v", expectedFileDiskPath, err)
					return
				}

				// Check if file content valid
				s.Equalf(expectedContent, content, "Uploaded content invalid")
			})
		})
	}
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
