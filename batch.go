package resque

import (
	"errors"

	redis "gopkg.in/redis.v5"
)

// Batch represents a batch enqueue action
type Batch struct {
	queue       string
	jobs        []Job
	redisClient *redis.Client
	namespace   string
}

// NewBatch prepares a Batch action
func (r *Client) NewBatch(queue string) (*Batch, error) {
	if queue == "" {
		return nil, errors.New("invalid queue name")
	}

	return &Batch{
		queue:       queue,
		redisClient: r.redisClient,
		namespace:   r.namespace,
	}, nil
}
