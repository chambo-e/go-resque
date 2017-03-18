package resque

import (
	"errors"
	"fmt"
)

// Enqueue enqueue one or more jobs into specified queue
func (r *Client) Enqueue(queue string, job Job) error {
	if queue == "" {
		return errors.New("invalid queue name")
	}

	if err := r.redisClient.SAdd(
		fmt.Sprintf("%s:queues", r.namespace),
		queue,
	).Err(); err != nil {
		return err
	}

	buffer, err := job.Marshal()
	if err != nil {
		return err
	}

	if err := r.redisClient.RPush(
		fmt.Sprintf("%s:queue:%s", r.namespace, queue),
		buffer,
	).Err(); err != nil {
		return err
	}

	return nil
}
