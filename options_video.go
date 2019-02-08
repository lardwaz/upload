package upload

var (
	defaultVideoOptions = &optionsVideo{
		codec: mpeg,
	}
)

const (
	mpeg = iota
)

type optionsVideo struct {
	Options
	codec int
}

// EvaluateVideoOptions returns list of video options
func EvaluateVideoOptions(opts []OptionVideo) *optionsVideo {
	optCopy := &optionsVideo{}
	*optCopy = *defaultVideoOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// OptionVideo is a function to modify options video
type OptionVideo func(*optionsVideo)

// Codec modifies codec video options
func Codec(codec int) OptionVideo {
	return func(o *optionsVideo) {
		o.codec = codec
	}
}
