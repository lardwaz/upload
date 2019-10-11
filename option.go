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
	SetDir(string) Options
	Destination() string
	SetDestination(string) Options
	MediaPrefixURL() string
	SetMediaPrefixURL(string) Options
	FileType() []types.Type
	AddFileType(types.Type) Options
	MaxSize() int
	SetMaxSize(sz int) Options
	ConvertTo(t types.Type) types.Type
	SetConvertTo(old types.Type, new types.Type) Options
	FileTypeExist(t types.Type) bool
}

// OptionsImage represents a set of image processing options
type OptionsImage interface {
	OptionsENV
	MinWidth() int
	SetMinWidth(w int) OptionsImage
	MinHeight() int
	SetMinHeight(h int) OptionsImage
	Formats() OptionsFormats
	SetFormats(opts OptionsFormats) OptionsImage
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
	SetName(n string) OptionsFormat
	Width() int
	SetWidth(w int) OptionsFormat
	Height() int
	SetHeight(h int) OptionsFormat
	Backdrop() OptionsBackdrop
	SetBackdrop(opts ...func(OptionsBackdrop)) OptionsFormat
	Watermark() OptionsWatermark
	SetWatermark(opts ...func(OptionsWatermark)) OptionsFormat
}

// OptionsBackdrop represents a set of backdrop processing options
type OptionsBackdrop interface {
	Path() string
	SetPath(p string) OptionsBackdrop
}

// OptionsWatermark represents a set of watermark processing options
type OptionsWatermark interface {
	Path() string
	SetPath(p string) OptionsWatermark
	Horizontal() int
	SetHorizontal(h int) OptionsWatermark
	Vertical() int
	SetVertical(v int) OptionsWatermark
	OffsetX() int
	SetOffsetX(x int) OptionsWatermark
	OffsetY() int
	SetOffsetY(y int) OptionsWatermark
}
