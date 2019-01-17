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

func EvaluateWatermarkOptions(opts []OptionWatermark) *optionsWatermark {
	optCopy := &optionsWatermark{}
	*optCopy = *defaultWatermarkOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type OptionWatermark func(*optionsWatermark)

func WatermarkHorizontal(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.horizontal = d
	}
}

func WatermarkVertical(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.vertical = d
	}
}

func WatermarkOffsetX(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.offsetX = d
	}
}

func WatermarkOffsetY(d int) OptionWatermark {
	return func(o *optionsWatermark) {
		o.offsetY = d
	}
}
