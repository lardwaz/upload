package upload

var (
	defaultWatermarkOptions = &OptsWatermark{}
)

// OptionsWatermark represents a set of watermark processing options
type OptionsWatermark interface {
	Horizontal() int
	SetHorizontal(h int)
	Vertical() int
	SetVertical(v int)
	OffsetX() int
	SetOffsetX(x int)
	OffsetY() int
	SetOffsetY(y int)
}

// OptsWatermark is an implementation of OptionsWatermark
type OptsWatermark struct {
	horizontal int
	vertical   int
	offsetX    int
	offsetY    int
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

// evaluateWatermarkOptions returns OptionsWatermark
func evaluateWatermarkOptions(opts ...func(OptionsWatermark)) OptionsWatermark {
	optCopy := &OptsWatermark{}
	*optCopy = *defaultWatermarkOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// WatermarkHorizontal returns OptionWatermark to modify WatermarkHorizontal
func WatermarkHorizontal(h int) func(OptionsWatermark) {
	return func(o OptionsWatermark) {
		o.SetHorizontal(h)
	}
}

// WatermarkVertical returns OptionWatermark to modify WatermarkVertical
func WatermarkVertical(v int) func(OptionsWatermark) {
	return func(o OptionsWatermark) {
		o.SetVertical(v)
	}
}

// WatermarkOffsetX returns OptionWatermark to modify WatermarkOffsetX
func WatermarkOffsetX(x int) func(OptionsWatermark) {
	return func(o OptionsWatermark) {
		o.SetOffsetX(x)
	}
}

// WatermarkOffsetY returns OptionWatermark to modify WatermarkOffsetY
func WatermarkOffsetY(y int) func(OptionsWatermark) {
	return func(o OptionsWatermark) {
		o.SetOffsetY(y)
	}
}
