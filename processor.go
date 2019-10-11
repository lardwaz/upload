package upload

// Processor represents a generic file processor (SMI)
type Processor interface {
	// Upload accepts an uploaded file and
	// returns a channel and error
	Process(Uploaded, bool) (Job, error)
}

// ImageProcessor represents an image processor
type ImageProcessor interface {
	Processor
	Options() OptionsImage
}
