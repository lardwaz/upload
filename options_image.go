package upload

import (
	"github.com/lsldigital/gocipe-upload/core"
)

var (
	defaultImageOptions = &optionsImage{
		minWidth:  core.NoLimit,
		minHeight: core.NoLimit,
	}
)

// format holds dimensions options for format
type format struct {
	name      string
	width     int
	height    int
	backdrop  bool              // (default: false) If true, will add a backdrop
	watermark *optionsWatermark // (default: nil) If not nil, will overlay an image as watermark at X,Y pos +-OffsetX,OffsetY
}

type optionsImage struct {
	minWidth  int
	minHeight int
	formats   []format
}

// EvaluateImageOptions returns optionsImage
func EvaluateImageOptions(opts ...OptionImage) *optionsImage {
	optCopy := &optionsImage{}
	*optCopy = *defaultImageOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// OptionImage is a function to modify options image
type OptionImage func(*optionsImage)

// MinWidth returns a function to modify MinWidth option image
func MinWidth(d int) OptionImage {
	return func(o *optionsImage) {
		o.minWidth = d
	}
}

// MinHeight returns a function to modify MinHeight option image
func MinHeight(d int) OptionImage {
	return func(o *optionsImage) {
		o.minHeight = d
	}
}

// Format returns a function to add Format option image
func Format(name string, width int, height int, backdrop bool, opts ...OptionWatermark) OptionImage {
	return func(o *optionsImage) {
		var watermarkOpts *optionsWatermark
		if len(opts) == 0 {
			watermarkOpts = nil
		} else {
			watermarkOpts = EvaluateWatermarkOptions(opts...)
		}

		imageFormat := format{
			name:      name,
			width:     width,
			height:    height,
			backdrop:  backdrop,
			watermark: watermarkOpts,
		}
		o.formats = append(o.formats, imageFormat)
	}
}
