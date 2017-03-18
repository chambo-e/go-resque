package resque

import "encoding/json"

// Job describe a resque job
type Job struct {
	Class string        `json:"class"`
	Args  []interface{} `json:"args"`
}

// NewJob is a helper to create a Job
func NewJob(klass string, args ...interface{}) Job {
	return Job{
		Class: klass,
		Args:  args,
	}
}

// Marshal returns json marshaled job buffer
func (job *Job) Marshal() ([]byte, error) {
	return json.Marshal(job)
}
