package option

import "go.lsl.digital/lardwaz/upload"

var (
	defaultFormatOptions = &OptsFormat{}
)

// OptsFormat holds dimensions options for Format
type OptsFormat struct {
	name      string
	width     int
	height    int
	backdrop  upload.OptionsBackdrop  // (default: nil) If not nil, will add a backdrop
	watermark upload.OptionsWatermark // (default: nil) If not nil, will overlay an image as watermark at X,Y pos +-OffsetX,OffsetY
}

// Name returns Name
func (o OptsFormat) Name() string {
	return o.name
}

// SetName sets the Name
func (o *OptsFormat) SetName(n string) {
	o.name = n
}

// Width returns Width
func (o OptsFormat) Width() int {
	return o.width
}

// SetWidth sets the Width
func (o *OptsFormat) SetWidth(w int) {
	o.width = w
}

// Height returns Height
func (o OptsFormat) Height() int {
	return o.height
}

// SetHeight sets the Height
func (o *OptsFormat) SetHeight(h int) {
	o.height = h
}

// Backdrop returns Backdrop
func (o OptsFormat) Backdrop() upload.OptionsBackdrop {
	return o.backdrop
}

// SetBackdrop sets the Backdrop
func (o *OptsFormat) SetBackdrop(opts ...func(upload.OptionsBackdrop)) {
	o.backdrop = EvaluateBackdropOptions(opts...)
}

// Watermark returns Watermark
func (o OptsFormat) Watermark() upload.OptionsWatermark {
	return o.watermark
}

// SetWatermark sets the watermark
func (o *OptsFormat) SetWatermark(opts ...func(upload.OptionsWatermark)) {
	o.watermark = EvaluateWatermarkOptions(opts...)
}

// EvaluateFormatOptions returns optionsImage
func EvaluateFormatOptions(opts ...func(upload.OptionsFormat)) upload.OptionsFormat {
	optCopy := &OptsFormat{}
	*optCopy = *defaultFormatOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// FormatName returns a function to modify format Name
func FormatName(n string) func(upload.OptionsFormat) {
	return func(o upload.OptionsFormat) {
		o.SetName(n)
	}
}

// FormatWidth returns a function to modify format Width
func FormatWidth(w int) func(upload.OptionsFormat) {
	return func(o upload.OptionsFormat) {
		o.SetWidth(w)
	}
}

// FormatHeight returns a function to modify format Height
func FormatHeight(h int) func(upload.OptionsFormat) {
	return func(o upload.OptionsFormat) {
		o.SetHeight(h)
	}
}

// FormatBackdrop returns a function to modify format Backdrop
func FormatBackdrop(opts ...func(upload.OptionsBackdrop)) func(upload.OptionsFormat) {
	return func(o upload.OptionsFormat) {
		o.SetBackdrop(opts...)
	}
}

// FormatWatermark returns a function to modify format watermark
func FormatWatermark(opts ...func(upload.OptionsWatermark)) func(upload.OptionsFormat) {
	return func(o upload.OptionsFormat) {
		o.SetWatermark(opts...)
	}
}
