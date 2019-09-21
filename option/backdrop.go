package option

import "go.lsl.digital/lardwaz/upload"

var (
	defaultBackdropOptions = &OptsBackdrop{}
)

// OptsBackdrop is an implementation of OptionsBackdrop
type OptsBackdrop struct {
	path string
}

// Path returns Path
func (o *OptsBackdrop) Path() string {
	return o.path
}

// SetPath sets Path
func (o *OptsBackdrop) SetPath(p string) {
	o.path = p
}

// EvaluateBackdropOptions returns OptionsBackdrop
func EvaluateBackdropOptions(opts ...func(upload.OptionsBackdrop)) upload.OptionsBackdrop {
	optCopy := &OptsBackdrop{}
	*optCopy = *defaultBackdropOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// BackdropPath returns OptionWatermark to modify BackdropPath
func BackdropPath(p string) func(upload.OptionsBackdrop) {
	return func(o upload.OptionsBackdrop) {
		o.SetPath(p)
	}
}
