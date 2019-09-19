package option

import sdk "go.lsl.digital/lardwaz/sdk/upload"

var (
	defaultWatermarkOptions = &OptsWatermark{}
)

// OptsWatermark is an implementation of OptionsWatermark
type OptsWatermark struct {
	path       string
	horizontal int
	vertical   int
	offsetX    int
	offsetY    int
}

// Path returns Path
func (o *OptsWatermark) Path() string {
	return o.path
}

// SetPath sets Path
func (o *OptsWatermark) SetPath(p string) {
	o.path = p
}

// Horizontal returns Horizontal
func (o *OptsWatermark) Horizontal() int {
	return o.horizontal
}

// SetHorizontal sets Horizontal
func (o *OptsWatermark) SetHorizontal(h int) {
	o.horizontal = h
}

// Vertical returns Vertical
func (o *OptsWatermark) Vertical() int {
	return o.vertical
}

// SetVertical sets Vertical
func (o *OptsWatermark) SetVertical(v int) {
	o.vertical = v
}

// OffsetX returns OffsetX
func (o *OptsWatermark) OffsetX() int {
	return o.offsetX
}

// SetOffsetX sets OffsetX
func (o *OptsWatermark) SetOffsetX(x int) {
	o.offsetX = x
}

// OffsetY returns OffsetY
func (o *OptsWatermark) OffsetY() int {
	return o.offsetY
}

// SetOffsetY sets OffsetY
func (o *OptsWatermark) SetOffsetY(y int) {
	o.offsetY = y
}

// EvaluateWatermarkOptions returns OptionsWatermark
func EvaluateWatermarkOptions(opts ...func(sdk.OptionsWatermark)) sdk.OptionsWatermark {
	optCopy := &OptsWatermark{}
	*optCopy = *defaultWatermarkOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// WatermarkPath returns OptionWatermark to modify WatermarkPath
func WatermarkPath(p string) func(sdk.OptionsWatermark) {
	return func(o sdk.OptionsWatermark) {
		o.SetPath(p)
	}
}

// WatermarkHorizontal returns OptionWatermark to modify WatermarkHorizontal
func WatermarkHorizontal(h int) func(sdk.OptionsWatermark) {
	return func(o sdk.OptionsWatermark) {
		o.SetHorizontal(h)
	}
}

// WatermarkVertical returns OptionWatermark to modify WatermarkVertical
func WatermarkVertical(v int) func(sdk.OptionsWatermark) {
	return func(o sdk.OptionsWatermark) {
		o.SetVertical(v)
	}
}

// WatermarkOffsetX returns OptionWatermark to modify WatermarkOffsetX
func WatermarkOffsetX(x int) func(sdk.OptionsWatermark) {
	return func(o sdk.OptionsWatermark) {
		o.SetOffsetX(x)
	}
}

// WatermarkOffsetY returns OptionWatermark to modify WatermarkOffsetY
func WatermarkOffsetY(y int) func(sdk.OptionsWatermark) {
	return func(o sdk.OptionsWatermark) {
		o.SetOffsetY(y)
	}
}
