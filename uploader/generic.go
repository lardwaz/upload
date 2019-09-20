package uploader

import (
	"fmt"

	"github.com/h2non/filetype"
	sdk "go.lsl.digital/lardwaz/sdk/upload"
	"go.lsl.digital/lardwaz/upload/file"
	"go.lsl.digital/lardwaz/upload/option"
)

// Generic is a Generic uploader
type Generic struct {
	Options sdk.Options
}

// NewGeneric returns Generic
func NewGeneric(opts ...func(sdk.Options)) *Generic {
	options := option.EvaluateOptions(opts...)
	return &Generic{Options: options}
}

// Upload method to satisfy uploader interface
func (u *Generic) Upload(name string, content []byte) (sdk.Uploaded, error) {
	fileType, err := filetype.Match(content)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving file type: %v", err)
	}

	if !u.Options.FileTypeExist(fileType) {
		return nil, fmt.Errorf("Unknown file type")
	}

	uploadedFile := file.NewGeneric(name, u.Options)

	if err := uploadedFile.Save(content, true); err != nil {
		return nil, err
	}

	newType := u.Options.ConvertTo(fileType)
	if err := uploadedFile.ChangeExt(newType.Extension); err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
