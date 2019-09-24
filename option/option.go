package option

import (
	"github.com/h2non/filetype/types"
	"go.lsl.digital/lardwaz/upload"
	utypes "go.lsl.digital/lardwaz/upload/types"
)

// Opts is an implementation of Options
type Opts struct {
	dir            string
	destination    string
	mediaPrefixURL string
	fileType       []types.Type
	maxSize        int
	convertTo      map[types.Type]types.Type
}

// NewUpload return a new options
func NewUpload() upload.Options {
	return &Opts{
		dir:            "media",
		mediaPrefixURL: "/media/",
		maxSize:        NoLimit,
		convertTo:      make(map[types.Type]types.Type),
	}
}

// Dir returns Dir
func (o Opts) Dir() string {
	return o.dir
}

// SetDir sets the Dir
func (o *Opts) SetDir(dir string) upload.Options {
	o.dir = dir

	return o
}

// Destination returns Destination
func (o Opts) Destination() string {
	return o.destination
}

// SetDestination sets the Destination
func (o *Opts) SetDestination(dest string) upload.Options {
	o.destination = dest

	return o
}

// MediaPrefixURL returns MediaPrefixURL
func (o Opts) MediaPrefixURL() string {
	return o.mediaPrefixURL
}

// SetMediaPrefixURL sets the MediaPrefixURL
func (o *Opts) SetMediaPrefixURL(url string) upload.Options {
	o.mediaPrefixURL = url

	return o
}

// FileType returns FileType
func (o Opts) FileType() []types.Type {
	return o.fileType
}

// AddFileType adds another the FileType
func (o *Opts) AddFileType(t types.Type) upload.Options {
	if utypes.IsValidType(t) && !o.FileTypeExist(t) {
		o.fileType = append(o.fileType, t)
	}

	return o
}

// MaxSize returns MaxSize
func (o Opts) MaxSize() int {
	return o.maxSize
}

// SetMaxSize sets the MaxSize
func (o *Opts) SetMaxSize(sz int) upload.Options {
	o.maxSize = sz

	return o
}

// ConvertTo returns ConvertTo
func (o Opts) ConvertTo(t types.Type) types.Type {
	return o.convertTo[t]
}

// SetConvertTo sets the MaxSize
func (o *Opts) SetConvertTo(old types.Type, new types.Type) upload.Options {
	o.convertTo[old] = new

	return o
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

// EvaluateOptions returns list of options
func EvaluateOptions(opts ...func(upload.Options)) upload.Options {
	optCopy := NewUpload()
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Dir returns a function to change Dir
func Dir(dir string) func(upload.Options) {
	return func(o upload.Options) {
		o.SetDir(dir)
	}
}

// Destination returns a function to change Destination
func Destination(dest string) func(upload.Options) {
	return func(o upload.Options) {
		o.SetDestination(dest)
	}
}

// MediaPrefixURL returns a function to change MediaPrefixURL
func MediaPrefixURL(url string) func(upload.Options) {
	return func(o upload.Options) {
		o.SetMediaPrefixURL(url)
	}
}

// FileType returns a function to change FileType
func FileType(t types.Type) func(upload.Options) {
	return func(o upload.Options) {
		o.AddFileType(t)
	}
}

// MaxSize returns a function to change MaxSize
func MaxSize(sz int) func(upload.Options) {
	return func(o upload.Options) {
		o.SetMaxSize(sz)
	}
}

// ConvertTo returns a function to change ConvertTo
func ConvertTo(old, new types.Type) func(upload.Options) {
	return func(o upload.Options) {
		o.SetConvertTo(old, new)
	}
}
