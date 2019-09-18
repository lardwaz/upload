package upload

// Processor represents a file Processor (SMI)
type Processor interface {
	// Upload accepts an uploaded file and
	// returns a channel and error
	Process(Uploaded, bool) (Job, error)
}
