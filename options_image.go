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

func EvaluateImageOptions(opts ...OptionImage) *optionsImage {
	optCopy := &optionsImage{}
	*optCopy = *defaultImageOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type OptionImage func(*optionsImage)

func MinWidth(d int) OptionImage {
	return func(o *optionsImage) {
		o.minWidth = d
	}
}

func MinHeight(d int) OptionImage {
	return func(o *optionsImage) {
		o.minHeight = d
	}
}

func Format(name string, width int, height int, backdrop bool, opts ...OptionWatermark) OptionImage {
	return func(o *optionsImage) {
		imageFormat := format{
			name:      name,
			width:     width,
			height:    height,
			backdrop:  backdrop,
			watermark: EvaluateWatermarkOptions(opts...),
		}
		o.formats = append(o.formats, imageFormat)
	}
}
