package upload

import (
	"fmt"

	"github.com/h2non/filetype"
	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

// GenericUploader is a generic uploader
type GenericUploader struct {
	Options sdk.Options
}

// NewGenericUploader returns GenericUploader
func NewGenericUploader(opts ...func(sdk.Options)) *GenericUploader {
	options := option.EvaluateOptions(opts...)
	return &GenericUploader{Options: options}
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

	uploadedFile := NewUploadedFile(name, u.Options)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	newType := u.Options.ConvertTo(fileType)
	if err := uploadedFile.ChangeExt(newType.Extension); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
