package upload

// OptionsFormats represents a list of OptionsFormats
type OptionsFormats interface {
	// Filter returns a OptionsFormats collection without elements filtered by fn (returning false)
	Filter(fn func(name string, item OptionsFormat) bool) OptionsFormats

	// Each loops over each item in the collection
	Each(fn func(name string, item OptionsFormat))

	// Set OptionsFormat in collection
	Set(item OptionsFormat)

	// Get a single item by name, if present
	Get(name string) (OptionsFormat, bool)

	// Length returns the number of items in the collection
	Length() int
}

// OptsFormats is an implementation of OptsFormats
type OptsFormats map[string]OptionsFormat

// NewOptionsFormats returns a new OptionsFormats
func NewOptionsFormats() OptionsFormats {
	return make(OptsFormats)
}

// Filter returns a Schematics collection without elements filtered by fn (returning false)
func (o OptsFormats) Filter(fn func(name string, item OptionsFormat) bool) OptionsFormats {
	s := NewOptionsFormats()
	for name, item := range o {
		if fn(name, item) {
			s.Set(item)
		}
	}

	return s
}

// Each loops over each item in the collection
func (o OptsFormats) Each(fn func(name string, item OptionsFormat)) {
	for name, item := range o {
		fn(name, item)
	}
}

// Set schematic in collection
func (o OptsFormats) Set(item OptionsFormat) {
	var name string

	if item == nil {
		name = ""
	} else {
		name = item.Name()
	}

	o[name] = item
}

// Get a single item by name, if present
func (o OptsFormats) Get(name string) (OptionsFormat, bool) {
	if i, ok := o[name]; ok {
		return i, true
	}

	return nil, false
}

// Length returns the number of items in the collection
func (o OptsFormats) Length() int {
	return len(o)
}
