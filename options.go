package upload

import (
	"github.com/lsldigital/gocipe-upload/core"
	"github.com/h2non/filetype/types"
)

var (
	defaultOptions = &Options{
		dir:            "media",
		mediaPrefixURL: "/media/",
		maxSize:        core.NoLimit,
	}
)

// Options represents a list of common options
type Options struct {
	dir            string
	destination    string
	mediaPrefixURL string
	fileType       []types.Type
	maxSize        int
	convertTo      string
}

// Dir returns Dir
func(o Options) Dir() string {
	return o.dir
}

// Destination returns Destination
func(o Options) Destination() string {
	return o.destination
}

// MediaPrefixURL returns MediaPrefixURL
func(o Options) MediaPrefixURL() string {
	return o.mediaPrefixURL
}

// FileType returns FileType
func(o Options) FileType() []types.Type {
	return o.fileType
}

// MaxSize returns MaxSize
func(o Options) MaxSize() int {
	return o.maxSize
}

// ConvertTo returns ConvertTo
func(o Options) ConvertTo() string {
	return o.convertTo
}

// FileTypeExist checks if filetype exists
func(o Options) FileTypeExist(t types.Type) bool {
	for _, fileType := range o.fileType {
		if fileType == t {
			return true
		}
	}

	return false
}

// EvaluateOptions returns list of options
func EvaluateOptions(opts ...Option) *Options {
	optCopy := &Options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option used to modify an option
type Option func(*Options)

// Dir returns a function to change Dir
func Dir(d string) Option {
	return func(o *Options) {
		o.dir = d
	}
}

// Destination returns a function to change Destination
func Destination(d string) Option {
	return func(o *Options) {
		o.destination = d
	}
}

// MediaPrefixURL returns a function to change MediaPrefixURL
func MediaPrefixURL(u string) Option {
	return func(o *Options) {
		o.mediaPrefixURL = u
	}
}

// FileType returns a function to change FileType
func FileType(t types.Type) Option {
	return func(o *Options) {
		if isValidType(t) && !o.FileTypeExist(t) {
			o.fileType = append(o.fileType, t)
		}
	}
}

// MaxSize returns a function to change MaxSize
func MaxSize(s int) Option {
	return func(o *Options) {
		o.maxSize = s
	}
}

// ConvertTo returns a function to change ConvertTo
func ConvertTo(t string) Option {
	return func(o *Options) {
		o.convertTo = t
	}
}