package upload

type VideoUpload struct {
	fileOpts  options
	videoOpts optionsVideo
}

func NewVideoUploader(common options, opts ...OptionVideo) VideoUpload {
	options := EvaluateVideoOptions(opts)
	return VideoUpload{fileOpts: common, videoOpts: *options}
}

func (u *VideoUpload) Upload(name string, content []byte) (UploadedFile, error) {
	uploadedFile := NewUploadedFile(name, u.fileOpts)

	if err := uploadedFile.Save(content, true); err != nil {
		return uploadedFile, err
	}

	if err := uploadedFile.ChangeExt(u.fileOpts.convertTo); err != nil {
		return uploadedFile, err
	}

	// Process video

	return uploadedFile, nil
}
