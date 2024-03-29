package processor_test

// Basic imports
import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/file"
	"go.lsl.digital/lardwaz/upload/option"
	"go.lsl.digital/lardwaz/upload/processor"
	"go.lsl.digital/lardwaz/upload/processor/box"
	"go.lsl.digital/lardwaz/upload/processor/position"
	utypes "go.lsl.digital/lardwaz/upload/types"
)

const (
	testDataFolder = "../testdata"
)

var update = flag.Bool("update", false, "update golden files")

type imageProcessTest struct {
	name                 string
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
	watermarkOptPath := option.WatermarkPath(filepath.Join(testDataFolder, "watermarks", "test_watermark.png"))
	backdropOptPath := option.BackdropPath(filepath.Join(testDataFolder, "backdrops", "test_backdrop.jpg"))

	// Set asset box
	box.Set(box.NewMockAsset())

	// Test cases
	s.imageProcessTests = []imageProcessTest{
		{"Normal No Format", "normal.jpg", "noformat_normal_out.jpg", false, processor.NewImage()},
		{"Normal No Format PNG", "normal.png", "noformat_normal_out.png", false, processor.NewImage()},
		{"Normal Format", "normal.jpg", "format_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("thumb"), option.FormatWidth(200), option.FormatHeight(200)))},
		{"Normal Format Negative Width & Height", "normal.jpg", "format_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("neg"), option.FormatWidth(-1), option.FormatHeight(-1)))},
		{"PROD Normal Format", "normal.jpg", "format_prod_normal_out.jpg", false, processor.NewImage(option.PROD(), option.Formats(option.FormatName("thumb"), option.FormatWidth(200), option.FormatHeight(200)))},
		{"Normal Format PNG", "normal.png", "format_normal_out.png", false, processor.NewImage(option.Formats(option.FormatName("thumb"), option.FormatWidth(200), option.FormatHeight(200)))},
		{"PROD Normal Format PNG", "normal.png", "format_prod_normal_out.png", false, processor.NewImage(option.PROD(), option.Formats(option.FormatName("thumb"), option.FormatWidth(200), option.FormatHeight(200)))},
		{"Normal Height Zero", "normal.jpg", "aspect_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("hzero"), option.FormatWidth(200)))},
		{"Normal Width Zero", "normal.jpg", "aspect_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("wzero"), option.FormatHeight(200)))},
		{"Normal Upscale", "normal.jpg", "upscale_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("upscale"), option.FormatWidth(500), option.FormatHeight(500)))},
		{"Small Width", "normal.jpg", "min_normal_out.jpg", true, processor.NewImage(option.MinWidth(500))},
		{"Small Height", "normal.jpg", "min_normal_out.jpg", true, processor.NewImage(option.MinHeight(500))},
		{"Invalid File Type", "damaged.jpg", "invalid_normal_out.jpg", true, processor.NewImage()},
		{"Invalid Image Type", "normal.gif", "invalid_normal_out.gif", true, processor.NewImage()},
		{"Watermark Top Left", "normal.jpg", "watermarked_tl_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Left), option.WatermarkVertical(position.Top))))},
		{"Watermark Top Center", "normal.jpg", "watermarked_tc_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Center), option.WatermarkVertical(position.Top))))},
		{"Watermark Top Right", "normal.jpg", "watermarked_tr_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Right), option.WatermarkVertical(position.Top))))},
		{"Watermark Bottom Left", "normal.jpg", "watermarked_bl_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Left), option.WatermarkVertical(position.Bottom))))},
		{"Watermark Bottom Center", "normal.jpg", "watermarked_bc_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Center), option.WatermarkVertical(position.Bottom))))},
		{"Watermark Bottom Right", "normal.jpg", "watermarked_br_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Right), option.WatermarkVertical(position.Bottom))))},
		{"Watermark Center Left", "normal.jpg", "watermarked_cl_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Left), option.WatermarkVertical(position.Center))))},
		{"Watermark Center Center", "normal.jpg", "watermarked_cc_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Center), option.WatermarkVertical(position.Center))))},
		{"Watermark Center Right", "normal.jpg", "watermarked_cr_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Right), option.WatermarkVertical(position.Center))))},
		{"Watermark Bad Pos", "normal.jpg", "watermarked_bad_prod_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(10), option.WatermarkVertical(10))))},
		{"PROD Watermark Bad Pos", "normal.jpg", "watermarked_bad_normal_out.jpg", false, processor.NewImage(option.PROD(), option.Formats(option.FormatName("water"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(10), option.WatermarkVertical(10))))},
		{"Watermark Bad Pos", "normal.jpg", "watermarked_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("damaged"), option.FormatWidth(400), option.FormatHeight(400), option.FormatWatermark(watermarkOptPath, option.WatermarkHorizontal(position.Center), option.WatermarkVertical(position.Center))))},
		{"Backdrop Landscape", "normal.jpg", "backdropped_normal_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("back"), option.FormatWidth(200), option.FormatHeight(200), option.FormatBackdrop(backdropOptPath)))},
		{"PROD Backdrop Landscape", "normal.jpg", "backdropped_prod_normal_out.jpg", false, processor.NewImage(option.PROD(), option.Formats(option.FormatName("back"), option.FormatWidth(200), option.FormatHeight(200), option.FormatBackdrop(backdropOptPath)))},
		{"Backdrop Portrait", "portrait.jpg", "backdropped_portrait_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("back"), option.FormatWidth(200), option.FormatHeight(200), option.FormatBackdrop(backdropOptPath)))},
		{"PROD Backdrop Portrait", "portrait.jpg", "backdropped_prod_portrait_out.jpg", false, processor.NewImage(option.PROD(), option.Formats(option.FormatName("back"), option.FormatWidth(200), option.FormatHeight(200), option.FormatBackdrop(backdropOptPath)))},
		{"Backdrop Damaged", "portrait.jpg", "backdropped_portrait_out.jpg", false, processor.NewImage(option.Formats(option.FormatName("damaged"), option.FormatWidth(200), option.FormatHeight(200), option.FormatBackdrop(backdropOptPath)))},
	}
}

func (s *ProcessorTestSuite) TestImageProcess() {
	// Common upload configurations
	commonOpts := []func(upload.Options){
		option.Dir(testDataFolder),
		option.MediaPrefixURL("/" + testDataFolder + "/"),
		option.FileType(utypes.TypeJPEG),
		option.FileType(utypes.TypeJPEG2),
		option.FileType(utypes.TypePNG),
		option.FileType(utypes.TypeGIF),
		option.FileType(utypes.TypeHEIF),
	}

	for _, tt := range s.imageProcessTests {
		s.Run(tt.name, func() {
			uploadedFile := file.NewMockGeneric(tt.inputFile, commonOpts...)
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
				// upload.Job timed out! Did we expect?
				if !tt.expectedProcessError {
					s.Failf("Cannot process file", "%s: Timed out!", job.File().DiskPath())
					return
				}
			case <-job.Done():
			// upload.Job done! We are good!

			case err = <-job.Failed():
				// upload.Job failed! Did we expect?
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
				s.Equalf(expectedContent, content, "upload.Uploaded content invalid")
			})
		})
	}
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}
