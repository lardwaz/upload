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
	options
	codec int
}

func EvaluateVideoOptions(opts []OptionVideo) *optionsVideo {
	optCopy := &optionsVideo{}
	*optCopy = *defaultVideoOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type OptionVideo func(*optionsVideo)

func Codec(codec int) OptionVideo {
	return func(o *optionsVideo) {
		o.codec = codec
	}
}
