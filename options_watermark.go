package upload

var (
	defaultWatermarkOptions = &OptionsWatermark{}
)

// OptionsWatermark holds the watermark position
type OptionsWatermark struct {
	horizontal int
	vertical   int
	offsetX    int
	offsetY    int
}

// EvaluateWatermarkOptions returns OptionsWatermark
func EvaluateWatermarkOptions(opts ...OptionWatermark) *OptionsWatermark {
	optCopy := &OptionsWatermark{}
	*optCopy = *defaultWatermarkOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// OptionWatermark is a function to modify watermark options
type OptionWatermark func(*OptionsWatermark)

// WatermarkHorizontal returns OptionWatermark to modify WatermarkHorizontal
func WatermarkHorizontal(d int) OptionWatermark {
	return func(o *OptionsWatermark) {
		o.horizontal = d
	}
}

// WatermarkVertical returns OptionWatermark to modify WatermarkVertical
func WatermarkVertical(d int) OptionWatermark {
	return func(o *OptionsWatermark) {
		o.vertical = d
	}
}

// WatermarkOffsetX returns OptionWatermark to modify WatermarkOffsetX
func WatermarkOffsetX(d int) OptionWatermark {
	return func(o *OptionsWatermark) {
		o.offsetX = d
	}
}

// WatermarkOffsetY returns OptionWatermark to modify WatermarkOffsetY
func WatermarkOffsetY(d int) OptionWatermark {
	return func(o *OptionsWatermark) {
		o.offsetY = d
	}
}
