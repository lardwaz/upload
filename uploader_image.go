package upload

type ImageUploader struct {
	Options  *Options
	Processor *ImageProcessor
}

// NewImageUploader returns ImageUploader
func NewImageUploader(common *Options, opts ...OptionImage) *ImageUploader {
	processor := NewImageProcessor(opts...)
	return &ImageUploader{Options: common, Processor: processor}
}

func (u *ImageUploader) Upload(name string, content []byte) (*UploadedFile, error) {
	uploadedFile := NewUploadedFile(name, *u.Options)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	if err := uploadedFile.ChangeExt(u.Options.ConvertTo()); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
