package upload

import (
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

// Supported file types by file upload
// Alias of types.Type
var (
	TypeJPEG 	= matchers.TypeJpeg
	TypeJPEG2   = matchers.TypeJpeg2000
	TypePNG 	= matchers.TypePng
	TypeGIF 	= matchers.TypeGif
	TypeHEIF 	= matchers.TypeHeif
	TypeMP3 	= matchers.TypeMp3
	TypeAAC 	= matchers.TypeAac
	TypeDOC 	= matchers.TypeDoc
	TypeDOCX 	= matchers.TypeDocx
	TypeXLS 	= matchers.TypeXls
	TypeXLSX 	= matchers.TypeXlsx
	TypePPT 	= matchers.TypePpt
	TypePPTX 	= matchers.TypePptx
	TypePDF 	= matchers.TypePdf
	TypeZIP 	= matchers.TypeZip
	TypeRAR 	= matchers.TypeRar
	Type7Z 	    = matchers.Type7z
	TypeMP4 	= matchers.TypeMp4
	TypeMOV 	= matchers.TypeMov
	TypeAVI 	= matchers.TypeAvi
	TypeWMV 	= matchers.TypeWmv
	TypeWEBM 	= matchers.TypeWebm
)

// SupportedTypes is a map of supported files by file upload
var SupportedTypes = matchers.Map{
	// Image
	TypeJPEG:     matchers.Jpeg,
	TypeJPEG2:    matchers.Jpeg2000,
	TypePNG:      matchers.Png,
	TypeGIF:      matchers.Gif,
	TypeHEIF:     matchers.Heif,
	// Audio
	TypeMP3:  	  matchers.Mp3,
	TypeAAC:      matchers.Aac,
	// Document
	TypeDOC:      matchers.Doc,
	TypeDOCX:     matchers.Docx,
	TypeXLS:      matchers.Xls,
	TypeXLSX:     matchers.Xlsx,
	TypePPT:      matchers.Ppt,
	TypePPTX:     matchers.Pptx,
	TypePDF:      matchers.Pdf,
	// Archive
	TypeZIP:      matchers.Zip,
	TypeRAR:      matchers.Rar,
	Type7Z:       matchers.SevenZ,
	// Video
	TypeMP4:      matchers.Mp4,	
	TypeMOV:      matchers.Mov,
	TypeAVI:      matchers.Avi,
	TypeWMV:      matchers.Wmv,	
	TypeWEBM:     matchers.Webm,
}

// isValidType checks if type supported by file upload
func isValidType(t types.Type) bool {
	_, valid := SupportedTypes[t]
	return valid
}

// isValidFile checks if file is supported by file upload
func isValidFile(content []byte) bool {
	kind := filetype.MatchMap(content, SupportedTypes)
	return kind != types.Unknown 
}

// isValidImage checks if file is an image supported by file upload
func isValidImage(content []byte) bool {
	return ( matchers.Jpeg(content) ||
		matchers.Jpeg2000(content) ||
		matchers.Png(content) ||
		matchers.Gif(content) )
}
