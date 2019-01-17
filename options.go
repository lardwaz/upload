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

func EvaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

func Dir(d string) Option {
	return func(o *options) {
		o.dir = d
	}
}

func Destination(d string) Option {
	return func(o *options) {
		o.destination = d
	}
}

func MediaPrefixURL(u string) Option {
	return func(o *options) {
		o.mediaPrefixURL = u
	}
}

func FileType(t uint8) Option {
	return func(o *options) {
		o.fileType = t
	}
}

func MaxSize(s int) Option {
	return func(o *options) {
		o.maxSize = s
	}
}

func ConvertTo(t string) Option {
	return func(o *options) {
		o.convertTo = t
	}
}
