package upload

import (
	"github.com/h2non/filetype/types"
)

// OptionsENV represents ENV(DEV / PROD) options
type OptionsENV interface {
	IsPROD() bool
	SetPROD(b bool)
}

// Options represents a set of upload options
type Options interface {
	Dir() string
	SetDir(string)
	Destination() string
	SetDestination(string)
	MediaPrefixURL() string
	SetMediaPrefixURL(string)
	FileType() []types.Type
	AddFileType(types.Type)
	MaxSize() int
	SetMaxSize(sz int)
	ConvertTo(t types.Type) types.Type
	SetConvertTo(old types.Type, new types.Type)
	FileTypeExist(t types.Type) bool
}

// OptionsImage represents a set of image processing options
type OptionsImage interface {
	OptionsENV
	MinWidth() int
	SetMinWidth(w int)
	MinHeight() int
	SetMinHeight(h int)
	Formats() OptionsFormats
	SetFormats(opts OptionsFormats)
}

// OptionsFormats represents a list of OptionsFormat
type OptionsFormats interface {
	// Filter returns a OptionsFormats collection without elements filtered by fn (returning false)
	Filter(fn func(name string, item OptionsFormat) bool) OptionsFormats

	// Each loops over each item in the collection
	Each(fn func(name string, item OptionsFormat))

	// Set OptionsFormat in collection
	Set(item OptionsFormat)

	// Get a single item by name, if present
	Get(name string) (OptionsFormat, bool)

	// Length returns the number of items in the collection
	Length() int
}

// OptionsFormat represents a set of format processing options
type OptionsFormat interface {
	Name() string
	SetName(n string)
	Width() int
	SetWidth(w int)
	Height() int
	SetHeight(h int)
	Backdrop() OptionsBackdrop
	SetBackdrop(opts ...func(OptionsBackdrop))
	Watermark() OptionsWatermark
	SetWatermark(opts ...func(OptionsWatermark))
}

// OptionsBackdrop represents a set of backdrop processing options
type OptionsBackdrop interface {
	Path() string
	SetPath(p string)
}

// OptionsWatermark represents a set of watermark processing options
type OptionsWatermark interface {
	Path() string
	SetPath(p string)
	Horizontal() int
	SetHorizontal(h int)
	Vertical() int
	SetVertical(v int)
	OffsetX() int
	SetOffsetX(x int)
	OffsetY() int
	SetOffsetY(y int)
}
