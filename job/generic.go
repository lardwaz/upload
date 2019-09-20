package job

import sdk "go.lsl.digital/lardwaz/sdk/upload"

// Generic represents current image file being processed
type Generic struct {
	file   sdk.Uploaded
	done   chan struct{}
	failed chan error
}

// NewGeneric returns a new Generic
func NewGeneric(file sdk.Uploaded) *Generic {
	return &Generic{
		file:   file,
		done:   make(chan struct{}),
		failed: make(chan error),
	}
}

// File returns the file sdk.Uploaded
func (j Generic) File() sdk.Uploaded {
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
