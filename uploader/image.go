package uploader

import (
	"fmt"

	"github.com/h2non/filetype"
	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/file"
	"go.lsl.digital/lardwaz/upload/option"
	utypes "go.lsl.digital/lardwaz/upload/types"
)

// Image is an image uploader
type Image struct {
	Options sdk.Options
}

// NewImage returns Image
func NewImage(opts ...func(sdk.Options)) *Image {
	options := option.EvaluateOptions(opts...)
	return &Image{Options: options}
}

// Upload method to satisfy uploader interface
func (u *Image) Upload(name string, content []byte) (sdk.Uploaded, error) {
	if !utypes.IsValidImage(content) {
		return nil, fmt.Errorf("Not a valid image")
	}

	uploadedFile := file.NewUploaded(name, u.Options)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	fileType, err := filetype.MatchFile(uploadedFile.DiskPath())
	if err != nil {
		return nil, fmt.Errorf("Error retrieving file type: %v", err)
	}

	newType := u.Options.ConvertTo(fileType)
	if err := uploadedFile.ChangeExt(newType.Extension); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
