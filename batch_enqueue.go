package resque

import (
	"errors"
	"fmt"

	redis "gopkg.in/redis.v5"
)

// Enqueue add one or more jobs in a Batch action
func (b *Batch) Enqueue(jobs ...Job) {
	b.jobs = append(b.jobs, jobs...)
}

// Execute enqueue every jobs currently inside the Batch action
// Batch is reseted on success and can be reused for further enqueue
func (b *Batch) Execute() error {
	if len(b.jobs) == 0 {
		return errors.New("No job provided")
	}

	if _, err := b.redisClient.Pipelined(func(pipe *redis.Pipeline) error {
		pipe.SAdd(fmt.Sprintf("%s:queues", b.namespace), b.queue)

		for _, job := range b.jobs {
			buffer, err := job.Marshal()
			if err != nil {
				return err
			}

			if err := pipe.RPush(
				fmt.Sprintf("%s:queue:%s", b.namespace, b.queue),
				buffer,
			).Err(); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	b.jobs = b.jobs[:0]

	return nil
}
