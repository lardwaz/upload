package upload

// GenericJob represents current image file being processed
type GenericJob struct {
	file   Uploaded
	done   chan struct{}
	failed chan error
}

// NewGenericJob returns a new GenericJob
func NewGenericJob(file Uploaded) *GenericJob {
	return &GenericJob{
		file:   file,
		done:   make(chan struct{}),
		failed: make(chan error),
	}
}

// File returns the file Uploaded
func (j GenericJob) File() Uploaded {
	return j.file
}

// Done returns a channel indicating if job is done
func (j GenericJob) Done() <-chan struct{} {
	return j.done
}

// SetDone sets the job as completed
func (j *GenericJob) SetDone() {
	j.done <- struct{}{}
}

// Failed returns a channel indicating if job has failed
func (j GenericJob) Failed() <-chan error {
	return j.failed
}

// SetFailed sets the job as failed
func (j *GenericJob) SetFailed(err error) {
	j.failed <- err
}
