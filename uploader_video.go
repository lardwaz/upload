package upload

type VideoUpload struct {
	fileOpts  Options
	videoOpts optionsVideo
}

func NewVideoUploader(common Options, opts ...OptionVideo) VideoUpload {
	options := EvaluateVideoOptions(opts)
	return VideoUpload{fileOpts: common, videoOpts: *options}
}

func (u *VideoUpload) Upload(name string, content []byte) (*UploadedFile, error) {
	uploadedFile := NewUploadedFile(name, u.fileOpts)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	if err := uploadedFile.ChangeExt(u.fileOpts.convertTo); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
