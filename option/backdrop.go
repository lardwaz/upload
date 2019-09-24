package option

import "go.lsl.digital/lardwaz/upload"

// OptsBackdrop is an implementation of OptionsBackdrop
type OptsBackdrop struct {
	path string
}

// NewBackdrop returns a new OptionsBackdrop
func NewBackdrop() upload.OptionsBackdrop {
	return &OptsBackdrop{}
}

// Path returns Path
func (o *OptsBackdrop) Path() string {
	return o.path
}

// SetPath sets Path
func (o *OptsBackdrop) SetPath(p string) upload.OptionsBackdrop {
	o.path = p

	return o
}

// EvaluateBackdropOptions returns OptionsBackdrop
func EvaluateBackdropOptions(opts ...func(upload.OptionsBackdrop)) upload.OptionsBackdrop {
	optCopy := NewBackdrop()
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
