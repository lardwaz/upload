package upload

const (
	// typeImageJPG denotes image of file type jpg
	typeImageJPG = "jpg"
	// typeImageJPEG denotes image of file type jpeg
	typeImageJPEG = "jpeg"
	// typeImagePNG denotes image of file type png
	typeImagePNG = "png"
)

type imageProcessor struct{}

// Process adds a job to process an image based on specific options
func (v *imageProcessor) Process(file UploadedFile, opts optionsImage) error {
	panic("Not implemented")
	return nil
}
