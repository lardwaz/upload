package option

import (
	"go.lsl.digital/lardwaz/upload"
)

// OptsFormats is an implementation of OptsFormats
type OptsFormats map[string]upload.OptionsFormat

// NewOptionsFormats returns a new OptionsFormats
func NewOptionsFormats() upload.OptionsFormats {
	return make(OptsFormats)
}

// Filter returns a Schematics collection without elements filtered by fn (returning false)
func (o OptsFormats) Filter(fn func(name string, item upload.OptionsFormat) bool) upload.OptionsFormats {
	s := NewOptionsFormats()
	for name, item := range o {
		if fn(name, item) {
			s.Set(item)
		}
	}

	return s
}

// Each loops over each item in the collection
func (o OptsFormats) Each(fn func(name string, item upload.OptionsFormat)) {
	for name, item := range o {
		fn(name, item)
	}
}

// Set schematic in collection
func (o OptsFormats) Set(item upload.OptionsFormat) {
	var name string

	if item == nil {
		name = ""
	} else {
		name = item.Name()
	}

	o[name] = item
}

// Get a single item by name, if present
func (o OptsFormats) Get(name string) (upload.OptionsFormat, bool) {
	if i, ok := o[name]; ok {
		return i, true
	}

	return nil, false
}

// Length returns the number of items in the collection
func (o OptsFormats) Length() int {
	return len(o)
}
