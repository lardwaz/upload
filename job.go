package upload

// Job represents current file being processed
type Job interface {
	File() Uploaded
	Done() <-chan struct{}
	SetDone()
	Failed() <-chan error
	SetFailed(error)
}
