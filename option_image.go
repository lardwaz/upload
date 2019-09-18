package upload

import (
	"log"

	"go.lsl.digital/lardwaz/upload/core"
)

var (
	defaultImageOptions = &OptsImage{
		minWidth:  core.NoLimit,
		minHeight: core.NoLimit,
	}
)

// OptionsImage represents a set of image processing options
type OptionsImage interface {
	MinWidth() int
	SetMinWidth(w int)
	MinHeight() int
	SetMinHeight(h int)
	Formats() OptionsFormats
	SetFormats(opts OptionsFormats)
}

// OptsImage is an implementation of OptionsImage
type OptsImage struct {
	minWidth  int
	minHeight int
	formats   OptionsFormats
}

// MinWidth returns MinWidth
func (o OptsImage) MinWidth() int {
	return o.minWidth
}

// SetMinWidth sets MinWidth
func (o *OptsImage) SetMinWidth(w int) {
	o.minWidth = w
}

// MinHeight returns MinHeight
func (o OptsImage) MinHeight() int {
	return o.minHeight
}

// SetMinHeight sets MinHeight
func (o *OptsImage) SetMinHeight(h int) {
	o.minHeight = h
}

// Formats returns Formats
func (o OptsImage) Formats() OptionsFormats {
	return o.formats
}

// SetFormats set Formats
func (o *OptsImage) SetFormats(opts OptionsFormats) {
	o.formats = opts
}

// evaluateImageOptions returns optionsImage
func evaluateImageOptions(opts ...func(OptionsImage)) OptionsImage {
	optCopy := &OptsImage{}
	*optCopy = *defaultImageOptions
	optCopy.formats = NewOptionsFormats()
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// MinWidth returns a function to modify MinWidth option image
func MinWidth(w int) func(OptionsImage) {
	return func(o OptionsImage) {
		o.SetMinWidth(w)
	}
}

// MinHeight returns a function to modify MinHeight option image
func MinHeight(h int) func(OptionsImage) {
	return func(o OptionsImage) {
		o.SetMinHeight(h)
	}
}

// Formats returns a function to add Format option image
func Formats(opts ...func(OptionsFormat)) func(OptionsImage) {
	return func(o OptionsImage) {
		format := evaluateFormatOptions(opts...)

		formats := o.Formats()

		log.Println("Formats", formats)

		formats.Set(format)

		o.SetFormats(formats)
	}
}
