package upload_test

// Basic imports
import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/lsldigital/gocipe-upload"

)

type genericUploadTest struct {
	name                 string
	inputFile            string
	expectedFile         string
	expectedUploadError  bool
	expectedContentError bool
	uploader             *upload.GenericUploader
}

type GenericUploaderTestSuite struct {
	suite.Suite
	genericUploadTests []genericUploadTest
}

func (s *GenericUploaderTestSuite) SetupSuite() {
	// Common upload configurations
	common := []upload.Option{
		upload.Dir(testDataFolder),
		upload.Destination("tmp"),
		upload.MediaPrefixURL("/"+testDataFolder+"/"),
		upload.FileType(upload.TypePDF),
		upload.FileType(upload.TypeMP3),
		upload.FileType(upload.TypeMP4),
		upload.FileType(upload.TypeZIP),
	}

	commonOpts := upload.EvaluateOptions(common...)
	commonMaxSizeOpts := upload.EvaluateOptions(append(common, upload.MaxSize(300))...)
	commonPDFMP3Opts := upload.EvaluateOptions(append(common, upload.ConvertTo(upload.TypePDF, upload.TypeMP3))...)

	// Test cases
	s.genericUploadTests = []genericUploadTest{
		{"PDF", "normal.pdf", "normal_out.pdf", false, false, upload.NewGenericUploader(commonOpts)},
		{"PDF to MP3", "normal.pdf", "normal_convert_out.mp3", false, false, upload.NewGenericUploader(commonPDFMP3Opts)},
		{"PDF", "normal.pdf", "normal_out.pdf", false, false, upload.NewGenericUploader(commonOpts)},
		{"MP3", "normal.mp3", "normal_out.mp3", false, false, upload.NewGenericUploader(commonOpts)},
		{"MP4", "normal.mp4", "normal_out.mp4", false, false, upload.NewGenericUploader(commonOpts)},
		{"ZIP", "normal.zip", "normal_out.zip", false, false, upload.NewGenericUploader(commonOpts)},
		{"ZIP (Max Size)", "normal.zip", "normal_out.zip", true, false, upload.NewGenericUploader(commonMaxSizeOpts)},
		{"TXT (invalid)", "normal.txt", "normal_out.txt", true, false, upload.NewGenericUploader(commonOpts)},
		{"JS (invalid)", "normal.js", "normal_out.js", true, false, upload.NewGenericUploader(commonOpts)},
		{"PHP (invalid)", "normal.php", "normal_out.php", true, false, upload.NewGenericUploader(commonOpts)},
		{"JPG (invalid + damaged)", "damaged.jpg", "damaged_out.php", true, false, upload.NewGenericUploader(commonOpts)},
	}
}

func (s *GenericUploaderTestSuite) TestGenericUpload() {
	for _, tt := range s.genericUploadTests {
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

func TestGenericUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(GenericUploaderTestSuite))
}
