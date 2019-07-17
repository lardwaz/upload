package upload

import (
	"go.lsl.digital/lardwaz/upload/core"
)

var (
	defaultImageOptions = &OptionsImage{
		minWidth:  core.NoLimit,
		minHeight: core.NoLimit,
	}
)

// Format holds dimensions options for Format
type Format struct {
	name      string
	width     int
	height    int
	backdrop  bool              // (default: false) If true, will add a backdrop
	watermark *OptionsWatermark // (default: nil) If not nil, will overlay an image as watermark at X,Y pos +-OffsetX,OffsetY
}

// Name returns Name option format
func(o Format) Name() string {
	return o.name
}

// Width returns Width option format
func(o Format) Width() int {
	return o.width
}

// Height returns Height option format
func(o Format) Height() int {
	return o.height
}

// Backdrop returns Backdrop option format
func(o Format) Backdrop() bool {
	return o.backdrop
}

// Watermark returns Watermark option format
func(o Format) Watermark() OptionsWatermark {
	return *o.watermark
}

type OptionsImage struct {
	minWidth  int
	minHeight int
	formats   []Format
}

// EvaluateImageOptions returns optionsImage
func EvaluateImageOptions(opts ...OptionImage) *OptionsImage {
	optCopy := &OptionsImage{}
	*optCopy = *defaultImageOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// MinWidth returns MinWidth option image
func(o OptionsImage) MinWidth() int {
	return o.minWidth
}

// MinHeight returns MinHeight option image
func(o OptionsImage) MinHeight() int {
	return o.minHeight
}

// Formats returns Formats option image
func(o OptionsImage) Formats() []Format {
	return o.formats
}

// OptionImage is a function to modify options image
type OptionImage func(*OptionsImage)

// MinWidth returns a function to modify MinWidth option image
func MinWidth(d int) OptionImage {
	return func(o *OptionsImage) {
		o.minWidth = d
	}
}

// MinHeight returns a function to modify MinHeight option image
func MinHeight(d int) OptionImage {
	return func(o *OptionsImage) {
		o.minHeight = d
	}
}

// Formats returns a function to add Format option image
func Formats(name string, width int, height int, backdrop bool, opts ...OptionWatermark) OptionImage {
	return func(o *OptionsImage) {
		var watermarkOpts *OptionsWatermark
		if len(opts) == 0 {
			watermarkOpts = nil
		} else {
			watermarkOpts = EvaluateWatermarkOptions(opts...)
		}

		imageFormat := Format{
			name:      name,
			width:     width,
			height:    height,
			backdrop:  backdrop,
			watermark: watermarkOpts,
		}
		o.formats = append(o.formats, imageFormat)
	}
}
