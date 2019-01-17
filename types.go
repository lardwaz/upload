package upload

// Supported file types by file upload
const (
	TypeInvalid uint8 = iota
	TypeImage
	TypeVideo
	TypeAudio
	TypeDocument
	TypeSheet
	TypeCSV
	TypePDF
)
