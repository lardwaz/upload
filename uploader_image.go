package upload

type ImageUploader struct {
	options  *options
	processor *ImageProcessor
}

// NewImageUploader returns ImageUploader
func NewImageUploader(common *options, opts ...OptionImage) *ImageUploader {
	processor := NewImageProcessor(opts...)
	return &ImageUploader{options: common, processor: processor}
}

func (u *ImageUploader) Upload(name string, content []byte) (*UploadedFile, error) {
	uploadedFile := NewUploadedFile(name, *u.options)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	if err := uploadedFile.ChangeExt(u.options.convertTo); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
