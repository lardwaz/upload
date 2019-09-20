package uploader_test

// Basic imports
import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/option"
	utypes "go.lsl.digital/lardwaz/upload/types"
	"go.lsl.digital/lardwaz/upload/uploader"
)

type genericUploadTest struct {
	name                 string
	inputFile            string
	expectedFile         string
	expectedUploadError  bool
	expectedContentError bool
	uploader             *uploader.Generic
}

type GenericUploaderTestSuite struct {
	suite.Suite
	genericUploadTests []genericUploadTest
}

func (s *GenericUploaderTestSuite) SetupSuite() {
	// Common upload configurations
	common := []func(sdk.Options){
		option.Dir(testDataFolder),
		option.Destination("tmp"),
		option.MediaPrefixURL("/" + testDataFolder + "/"),
		option.FileType(utypes.TypePDF),
		option.FileType(utypes.TypeMP3),
		option.FileType(utypes.TypeMP4),
		option.FileType(utypes.TypeZIP),
	}

	commonMaxSizeOpts := append(common, option.MaxSize(300))
	commonPDFMP3Opts := append(common, option.ConvertTo(utypes.TypePDF, utypes.TypeMP3))

	// Test cases
	s.genericUploadTests = []genericUploadTest{
		{"PDF", "normal.pdf", "normal_out.pdf", false, false, uploader.NewGeneric(common...)},
		{"PDF to MP3", "normal.pdf", "normal_convert_out.mp3", false, false, uploader.NewGeneric(commonPDFMP3Opts...)},
		{"PDF", "normal.pdf", "normal_out.pdf", false, false, uploader.NewGeneric(common...)},
		{"MP3", "normal.mp3", "normal_out.mp3", false, false, uploader.NewGeneric(common...)},
		{"MP4", "normal.mp4", "normal_out.mp4", false, false, uploader.NewGeneric(common...)},
		{"ZIP", "normal.zip", "normal_out.zip", false, false, uploader.NewGeneric(common...)},
		{"ZIP (Max Size)", "normal.zip", "normal_out.zip", true, false, uploader.NewGeneric(commonMaxSizeOpts...)},
		{"TXT (invalid)", "normal.txt", "normal_out.txt", true, false, uploader.NewGeneric(common...)},
		{"JS (invalid)", "normal.js", "normal_out.js", true, false, uploader.NewGeneric(common...)},
		{"PHP (invalid)", "normal.php", "normal_out.php", true, false, uploader.NewGeneric(common...)},
		{"JPG (invalid + damaged)", "damaged.jpg", "damaged_out.php", true, false, uploader.NewGeneric(common...)},
	}
}

func (s *GenericUploaderTestSuite) TestGenericUpload() {
	for _, tt := range s.genericUploadTests {
		s.Run(tt.name, func() {
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
			s.Equalf(expectedContent, content, "sdk.Uploaded content invalid")
		})
	}
}

func TestGenericUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(GenericUploaderTestSuite))
}
