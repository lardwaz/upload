package upload

import (
	"fmt"

	"github.com/h2non/filetype"
)

// GenericUploader is a generic uploader
// Does not have a processor
type GenericUploader struct {
	Options *Options
}

// NewGenericUploader returns GenericUploader
// Does not have a processor
func NewGenericUploader(common *Options) *GenericUploader {
	return &GenericUploader{Options: common}
}

// Upload method to satisfy uploader interface
func (u *GenericUploader) Upload(name string, content []byte) (Uploaded, error) {
	fileType, err := filetype.Match(content)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving file type: %v", err)
	}

	if !u.Options.FileTypeExist(fileType) {
		return nil, fmt.Errorf("Unknown file type")
	}

	uploadedFile := NewUploadedFile(name, *u.Options)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	newType := u.Options.ConvertTo(fileType)
	if err := uploadedFile.ChangeExt(newType.Extension); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
