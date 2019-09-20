package option

import (
	sdk "go.lsl.digital/lardwaz/sdk/upload"
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
	formats   sdk.OptionsFormats
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
func (o OptsImage) Formats() sdk.OptionsFormats {
	return o.formats
}

// SetFormats set Formats
func (o *OptsImage) SetFormats(opts sdk.OptionsFormats) {
	o.formats = opts
}

// EvaluateImageOptions returns optionsImage
func EvaluateImageOptions(opts ...func(sdk.OptionsImage)) sdk.OptionsImage {
	optCopy := &OptsImage{}
	*optCopy = *defaultImageOptions
	optCopy.formats = NewOptionsFormats()
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// MinWidth returns a function to modify MinWidth option image
func MinWidth(w int) func(sdk.OptionsImage) {
	return func(o sdk.OptionsImage) {
		o.SetMinWidth(w)
	}
}

// MinHeight returns a function to modify MinHeight option image
func MinHeight(h int) func(sdk.OptionsImage) {
	return func(o sdk.OptionsImage) {
		o.SetMinHeight(h)
	}
}

// Formats returns a function to add Format option image
func Formats(opts ...func(sdk.OptionsFormat)) func(sdk.OptionsImage) {
	return func(o sdk.OptionsImage) {
		format := EvaluateFormatOptions(opts...)

		formats := o.Formats()

		formats.Set(format)

		o.SetFormats(formats)
	}
}

// PROD returns a function to modify ENV
func PROD() func(sdk.OptionsImage) {
	return func(o sdk.OptionsImage) {
		o.SetPROD(true)
	}
}
