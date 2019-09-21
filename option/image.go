package option

import (
	"go.lsl.digital/lardwaz/upload"
)

var (
	defaultImageOptions = &OptsImage{
		minWidth:  NoLimit,
		minHeight: NoLimit,
	}
)

// OptsImage is an implementation of OptionsImage
type OptsImage struct {
	OptsENV
	minWidth  int
	minHeight int
	formats   upload.OptionsFormats
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
func (o OptsImage) Formats() upload.OptionsFormats {
	return o.formats
}

// SetFormats set Formats
func (o *OptsImage) SetFormats(opts upload.OptionsFormats) {
	o.formats = opts
}

// EvaluateImageOptions returns optionsImage
func EvaluateImageOptions(opts ...func(upload.OptionsImage)) upload.OptionsImage {
	optCopy := &OptsImage{}
	*optCopy = *defaultImageOptions
	optCopy.formats = NewOptionsFormats()
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// MinWidth returns a function to modify MinWidth option image
func MinWidth(w int) func(upload.OptionsImage) {
	return func(o upload.OptionsImage) {
		o.SetMinWidth(w)
	}
}

// MinHeight returns a function to modify MinHeight option image
func MinHeight(h int) func(upload.OptionsImage) {
	return func(o upload.OptionsImage) {
		o.SetMinHeight(h)
	}
}

// Formats returns a function to add Format option image
func Formats(opts ...func(upload.OptionsFormat)) func(upload.OptionsImage) {
	return func(o upload.OptionsImage) {
		format := EvaluateFormatOptions(opts...)

		formats := o.Formats()

		formats.Set(format)

		o.SetFormats(formats)
	}
}

// PROD returns a function to modify ENV
func PROD() func(upload.OptionsImage) {
	return func(o upload.OptionsImage) {
		o.SetPROD(true)
	}
}
