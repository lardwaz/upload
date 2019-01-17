package upload

// Processor represents a file Processor (SMI)
type Processor interface {
	Process(string, []byte, ...Option) error
}
