package upload

var (
	defaultFormatOptions = &OptsFormat{}
)

// OptionsFormat represents a set of format processing options
type OptionsFormat interface {
	Name() string
	SetName(n string)
	Width() int
	SetWidth(w int)
	Height() int
	SetHeight(h int)
	Backdrop() bool
	SetBackdrop(b bool)
	Watermark() OptionsWatermark
	SetWatermark(opts ...func(OptionsWatermark))
}

// OptsFormat holds dimensions options for Format
type OptsFormat struct {
	name      string
	width     int
	height    int
	backdrop  bool             // (default: false) If true, will add a backdrop
	watermark OptionsWatermark // (default: nil) If not nil, will overlay an image as watermark at X,Y pos +-OffsetX,OffsetY
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
func (o OptsFormat) Backdrop() bool {
	return o.backdrop
}

// SetBackdrop sets the Backdrop
func (o *OptsFormat) SetBackdrop(b bool) {
	o.backdrop = b
}

// Watermark returns Watermark
func (o OptsFormat) Watermark() OptionsWatermark {
	return o.watermark
}

// SetWatermark sets the watermark
func (o *OptsFormat) SetWatermark(opts ...func(OptionsWatermark)) {
	o.watermark = evaluateWatermarkOptions(opts...)
}

// evaluateFormatOptions returns optionsImage
func evaluateFormatOptions(opts ...func(OptionsFormat)) OptionsFormat {
	optCopy := &OptsFormat{}
	*optCopy = *defaultFormatOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// FormatName returns a function to modify format Name
func FormatName(n string) func(OptionsFormat) {
	return func(o OptionsFormat) {
		o.SetName(n)
	}
}

// FormatWidth returns a function to modify format Width
func FormatWidth(w int) func(OptionsFormat) {
	return func(o OptionsFormat) {
		o.SetWidth(w)
	}
}

// FormatHeight returns a function to modify format Height
func FormatHeight(h int) func(OptionsFormat) {
	return func(o OptionsFormat) {
		o.SetHeight(h)
	}
}

// FormatBackdrop returns a function to modify format Backdrop
func FormatBackdrop(b bool) func(OptionsFormat) {
	return func(o OptionsFormat) {
		o.SetBackdrop(b)
	}
}

// FormatWatermark returns a function to modify format watermark
func FormatWatermark(opts ...func(OptionsWatermark)) func(OptionsFormat) {
	return func(o OptionsFormat) {
		o.SetWatermark(opts...)
	}
}
