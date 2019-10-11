package job

import "go.lsl.digital/lardwaz/upload"

// Generic represents current image file being processed
type Generic struct {
	file   upload.Uploaded
	done   chan struct{}
	failed chan error
}

// NewGeneric returns a new Generic
func NewGeneric(file upload.Uploaded) *Generic {
	return &Generic{
		file:   file,
		done:   make(chan struct{}),
		failed: make(chan error),
	}
}

// File returns the file upload.Uploaded
func (j Generic) File() upload.Uploaded {
	return j.file
}

// Done returns a channel indicating if job is done
func (j Generic) Done() <-chan struct{} {
	return j.done
}

// SetDone sets the job as completed
func (j *Generic) SetDone() {
	j.done <- struct{}{}
}

// Failed returns a channel indicating if job has failed
func (j Generic) Failed() <-chan error {
	return j.failed
}

// SetFailed sets the job as failed
func (j *Generic) SetFailed(err error) {
	j.failed <- err
}
