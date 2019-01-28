package upload

var (
	defaultWatermarkOptions = &optionsWatermark{}
)

// optionsWatermark holds the watermark position
type optionsWatermark struct {
	horizontal int
	vertical   int
	offsetX    int
	offsetY    int
}

// EvaluateWatermarkOptions returns optionsWatermark
func EvaluateWatermarkOptions(opts ...OptionWatermark) *optionsWatermark {
	optCopy := &optionsWatermark{}
	*optCopy = *defaultWatermarkOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// OptionWatermark is a function to modify watermark options
type OptionWatermark func(*optionsWatermark)

// WatermarkHorizontal returns OptionWatermark to modify WatermarkHorizontal
func WatermarkHorizontal(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.horizontal = d
	}
}

// WatermarkVertical returns OptionWatermark to modify WatermarkVertical
func WatermarkVertical(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.vertical = d
	}
}

// WatermarkOffsetX returns OptionWatermark to modify WatermarkOffsetX
func WatermarkOffsetX(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.offsetX = d
	}
}

// WatermarkOffsetY returns OptionWatermark to modify WatermarkOffsetY
func WatermarkOffsetY(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.offsetY = d
	}
}
