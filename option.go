package upload

import (
	"github.com/h2non/filetype/types"
	"go.lsl.digital/lardwaz/upload/core"
)

var (
	defaultOptions = &Opts{
		dir:            "media",
		mediaPrefixURL: "/media/",
		maxSize:        core.NoLimit,
		convertTo:      make(map[types.Type]types.Type),
	}
)

// Options represents a set of upload options
type Options interface {
	Dir() string
	SetDir(string)
	Destination() string
	SetDestination(string)
	MediaPrefixURL() string
	SetMediaPrefixURL(string)
	FileType() []types.Type
	AddFileType(types.Type)
	MaxSize() int
	SetMaxSize(sz int)
	ConvertTo(t types.Type) types.Type
	SetConvertTo(old types.Type, new types.Type)
	FileTypeExist(t types.Type) bool
}

// Opts is an implementation of Options
type Opts struct {
	dir            string
	destination    string
	mediaPrefixURL string
	fileType       []types.Type
	maxSize        int
	convertTo      map[types.Type]types.Type
}

// Dir returns Dir
func (o Opts) Dir() string {
	return o.dir
}

// SetDir sets the Dir
func (o *Opts) SetDir(dir string) {
	o.dir = dir
}

// Destination returns Destination
func (o Opts) Destination() string {
	return o.destination
}

// SetDestination sets the Destination
func (o *Opts) SetDestination(dest string) {
	o.destination = dest
}

// MediaPrefixURL returns MediaPrefixURL
func (o Opts) MediaPrefixURL() string {
	return o.mediaPrefixURL
}

// SetMediaPrefixURL sets the MediaPrefixURL
func (o *Opts) SetMediaPrefixURL(url string) {
	o.mediaPrefixURL = url
}

// FileType returns FileType
func (o Opts) FileType() []types.Type {
	return o.fileType
}

// AddFileType adds another the FileType
func (o *Opts) AddFileType(t types.Type) {
	if isValidType(t) && !o.FileTypeExist(t) {
		o.fileType = append(o.fileType, t)
	}
}

// MaxSize returns MaxSize
func (o Opts) MaxSize() int {
	return o.maxSize
}

// SetMaxSize sets the MaxSize
func (o *Opts) SetMaxSize(sz int) {
	o.maxSize = sz
}

// ConvertTo returns ConvertTo
func (o Opts) ConvertTo(t types.Type) types.Type {
	return o.convertTo[t]
}

// SetConvertTo sets the MaxSize
func (o *Opts) SetConvertTo(old types.Type, new types.Type) {
	o.convertTo[old] = new
}

// FileTypeExist checks if filetype exists
func (o Opts) FileTypeExist(t types.Type) bool {
	for _, fileType := range o.fileType {
		if fileType == t {
			return true
		}
	}

	return false
}

// evaluateOptions returns list of options
func evaluateOptions(opts ...func(Options)) Options {
	optCopy := &Opts{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Dir returns a function to change Dir
func Dir(dir string) func(Options) {
	return func(o Options) {
		o.SetDir(dir)
	}
}

// Destination returns a function to change Destination
func Destination(dest string) func(Options) {
	return func(o Options) {
		o.SetDestination(dest)
	}
}

// MediaPrefixURL returns a function to change MediaPrefixURL
func MediaPrefixURL(url string) func(Options) {
	return func(o Options) {
		o.SetMediaPrefixURL(url)
	}
}

// FileType returns a function to change FileType
func FileType(t types.Type) func(Options) {
	return func(o Options) {
		o.AddFileType(t)
	}
}

// MaxSize returns a function to change MaxSize
func MaxSize(sz int) func(Options) {
	return func(o Options) {
		o.SetMaxSize(sz)
	}
}

// ConvertTo returns a function to change ConvertTo
func ConvertTo(old, new types.Type) func(Options) {
	return func(o Options) {
		o.SetConvertTo(old, new)
	}
}
