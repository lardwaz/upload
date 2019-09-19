package option

import sdk "go.lsl.digital/lardwaz/sdk/upload"

var (
	defaultFormatOptions = &OptsFormat{}
)

// OptsFormat holds dimensions options for Format
type OptsFormat struct {
	name      string
	width     int
	height    int
	backdrop  sdk.OptionsBackdrop  // (default: nil) If not nil, will add a backdrop
	watermark sdk.OptionsWatermark // (default: nil) If not nil, will overlay an image as watermark at X,Y pos +-OffsetX,OffsetY
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
func (o OptsFormat) Backdrop() sdk.OptionsBackdrop {
	return o.backdrop
}

// SetBackdrop sets the Backdrop
func (o *OptsFormat) SetBackdrop(opts ...func(sdk.OptionsBackdrop)) {
	o.backdrop = EvaluateBackdropOptions(opts...)
}

// Watermark returns Watermark
func (o OptsFormat) Watermark() sdk.OptionsWatermark {
	return o.watermark
}

// SetWatermark sets the watermark
func (o *OptsFormat) SetWatermark(opts ...func(sdk.OptionsWatermark)) {
	o.watermark = EvaluateWatermarkOptions(opts...)
}

// EvaluateFormatOptions returns optionsImage
func EvaluateFormatOptions(opts ...func(sdk.OptionsFormat)) sdk.OptionsFormat {
	optCopy := &OptsFormat{}
	*optCopy = *defaultFormatOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// FormatName returns a function to modify format Name
func FormatName(n string) func(sdk.OptionsFormat) {
	return func(o sdk.OptionsFormat) {
		o.SetName(n)
	}
}

// FormatWidth returns a function to modify format Width
func FormatWidth(w int) func(sdk.OptionsFormat) {
	return func(o sdk.OptionsFormat) {
		o.SetWidth(w)
	}
}

// FormatHeight returns a function to modify format Height
func FormatHeight(h int) func(sdk.OptionsFormat) {
	return func(o sdk.OptionsFormat) {
		o.SetHeight(h)
	}
}

// FormatBackdrop returns a function to modify format Backdrop
func FormatBackdrop(opts ...func(sdk.OptionsBackdrop)) func(sdk.OptionsFormat) {
	return func(o sdk.OptionsFormat) {
		o.SetBackdrop(opts...)
	}
}

// FormatWatermark returns a function to modify format watermark
func FormatWatermark(opts ...func(sdk.OptionsWatermark)) func(sdk.OptionsFormat) {
	return func(o sdk.OptionsFormat) {
		o.SetWatermark(opts...)
	}
}
