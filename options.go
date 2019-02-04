package upload

import (
	"github.com/lsldigital/gocipe-upload/core"
)

var (
	defaultOptions = &options{
		dir:            "media",
		mediaPrefixURL: "/media/",
		fileType:       TypeImage,
		maxSize:        core.NoLimit,
		convertTo:      typeImageJPG,
	}
)

type options struct {
	dir            string
	destination    string
	mediaPrefixURL string
	fileType       uint8
	maxSize        int
	convertTo      string
}

// EvaluateOptions returns list of options
func EvaluateOptions(opts ...Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option used to modify an option
type Option func(*options)

// Dir returns a function to change Dir
func Dir(d string) Option {
	return func(o *options) {
		o.dir = d
	}
}

// Destination returns a function to change Destination
func Destination(d string) Option {
	return func(o *options) {
		o.destination = d
	}
}

// MediaPrefixURL returns a function to change MediaPrefixURL
func MediaPrefixURL(u string) Option {
	return func(o *options) {
		o.mediaPrefixURL = u
	}
}

// FileType returns a function to change FileType
func FileType(t uint8) Option {
	return func(o *options) {
		o.fileType = t
	}
}

// MaxSize returns a function to change MaxSize
func MaxSize(s int) Option {
	return func(o *options) {
		o.maxSize = s
	}
}

// ConvertTo returns a function to change ConvertTo
func ConvertTo(t string) Option {
	return func(o *options) {
		o.convertTo = t
	}
}
