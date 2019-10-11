package option

import "go.lsl.digital/lardwaz/upload"

// OptsWatermark is an implementation of OptionsWatermark
type OptsWatermark struct {
	path       string
	horizontal int
	vertical   int
	offsetX    int
	offsetY    int
}

// NewWatermark returns a new OptionsWatermark
func NewWatermark() upload.OptionsWatermark {
	return &OptsWatermark{}
}

// Path returns Path
func (o *OptsWatermark) Path() string {
	return o.path
}

// SetPath sets Path
func (o *OptsWatermark) SetPath(p string) upload.OptionsWatermark {
	o.path = p

	return o
}

// Horizontal returns Horizontal
func (o *OptsWatermark) Horizontal() int {
	return o.horizontal
}

// SetHorizontal sets Horizontal
func (o *OptsWatermark) SetHorizontal(h int) upload.OptionsWatermark {
	o.horizontal = h

	return o
}

// Vertical returns Vertical
func (o *OptsWatermark) Vertical() int {
	return o.vertical
}

// SetVertical sets Vertical
func (o *OptsWatermark) SetVertical(v int) upload.OptionsWatermark {
	o.vertical = v

	return o
}

// OffsetX returns OffsetX
func (o *OptsWatermark) OffsetX() int {
	return o.offsetX
}

// SetOffsetX sets OffsetX
func (o *OptsWatermark) SetOffsetX(x int) upload.OptionsWatermark {
	o.offsetX = x

	return o
}

// OffsetY returns OffsetY
func (o *OptsWatermark) OffsetY() int {
	return o.offsetY
}

// SetOffsetY sets OffsetY
func (o *OptsWatermark) SetOffsetY(y int) upload.OptionsWatermark {
	o.offsetY = y

	return o
}

// EvaluateWatermarkOptions returns OptionsWatermark
func EvaluateWatermarkOptions(opts ...func(upload.OptionsWatermark)) upload.OptionsWatermark {
	optCopy := NewWatermark()
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// WatermarkPath returns OptionWatermark to modify WatermarkPath
func WatermarkPath(p string) func(upload.OptionsWatermark) {
	return func(o upload.OptionsWatermark) {
		o.SetPath(p)
	}
}

// WatermarkHorizontal returns OptionWatermark to modify WatermarkHorizontal
func WatermarkHorizontal(h int) func(upload.OptionsWatermark) {
	return func(o upload.OptionsWatermark) {
		o.SetHorizontal(h)
	}
}

// WatermarkVertical returns OptionWatermark to modify WatermarkVertical
func WatermarkVertical(v int) func(upload.OptionsWatermark) {
	return func(o upload.OptionsWatermark) {
		o.SetVertical(v)
	}
}

// WatermarkOffsetX returns OptionWatermark to modify WatermarkOffsetX
func WatermarkOffsetX(x int) func(upload.OptionsWatermark) {
	return func(o upload.OptionsWatermark) {
		o.SetOffsetX(x)
	}
}

// WatermarkOffsetY returns OptionWatermark to modify WatermarkOffsetY
func WatermarkOffsetY(y int) func(upload.OptionsWatermark) {
	return func(o upload.OptionsWatermark) {
		o.SetOffsetY(y)
	}
}
