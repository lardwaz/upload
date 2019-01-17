package upload

type ImageUpload struct {
	fileOpts  options
	imageOpts optionsImage
}

func NewImageUploader(common options, opts ...OptionImage) ImageUpload {
	options := EvaluateImageOptions(opts)
	return ImageUpload{fileOpts: common, imageOpts: *options}
}

func (u *ImageUpload) Upload(name string, content []byte) (UploadedFile, error) {
	uploadedFile := NewUploadedFile(name, u.fileOpts)

	if err := uploadedFile.Save(content, true); err != nil {
		return uploadedFile, err
	}

	if err := uploadedFile.ChangeExt(u.fileOpts.convertTo); err != nil {
		return uploadedFile, err
	}

	// Process image

	return uploadedFile, nil
}
