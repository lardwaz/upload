package option

import sdk "go.lsl.digital/lardwaz/sdk/upload"

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
func EvaluateBackdropOptions(opts ...func(sdk.OptionsBackdrop)) sdk.OptionsBackdrop {
	optCopy := &OptsBackdrop{}
	*optCopy = *defaultBackdropOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// BackdropPath returns OptionWatermark to modify BackdropPath
func BackdropPath(p string) func(sdk.OptionsBackdrop) {
	return func(o sdk.OptionsBackdrop) {
		o.SetPath(p)
	}
}
