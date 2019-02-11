package upload

// Supported file types by file upload
const (
	TypeInvalid Type = iota
	TypeImage
	TypeVideo
	TypeAudio
	TypeDocument
	TypeSheet
	TypeCSV
	TypePDF
)

// Type represents a filetype
type Type uint8