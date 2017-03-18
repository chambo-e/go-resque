package resque

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch_Enqueue(t *testing.T) {
	tests := []struct {
		jobs  []Job
		queue string
	}{
		{
			jobs:  []Job{},
			queue: "test",
		},
		{
			jobs: []Job{
				NewJob("test", "hello", "world"),
			},
			queue: "test",
		},
		{
			jobs: []Job{
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
			},
			queue: "test",
		},
		{
			jobs: []Job{
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
			},
			queue: "test",
		},
	}

	cli, err := New(Configuration{
		RedisURI: "redis://localhost:6379",
	})
	assert.Nil(t, err, "should be nil")
	assert.NotNil(t, cli, "should not be nil")

	for _, test := range tests {
		batch, err := cli.NewBatch(test.queue)
		assert.Nil(t, err, "should be nil")

		batch.Enqueue(test.jobs...)

		assert.Equal(t, test.queue, batch.queue, "should be of same queue")
		assert.Equal(t, len(test.jobs), len(batch.jobs), "should be of same size")
	}
}

func TestBatch_Execute(t *testing.T) {
	tests := []struct {
		jobs          []Job
		queue         string
		shouldSuccess bool
	}{
		{
			jobs:          []Job{},
			queue:         "test",
			shouldSuccess: false,
		},
		{
			jobs: []Job{
				NewJob("test", "hello", "world"),
			},
			queue:         "test",
			shouldSuccess: true,
		},
		{
			jobs: []Job{
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
				NewJob("tata", "hello", true, 0.34),
			},
			queue:         "test",
			shouldSuccess: true,
		},
		{
			jobs: []Job{
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
				NewJob("tutu", 1, "world"),
			},
			queue:         "test",
			shouldSuccess: true,
		},
	}

	cli, err := New(Configuration{
		RedisURI: "redis://localhost:6379",
	})
	assert.Nil(t, err, "should be nil")
	assert.NotNil(t, cli, "should not be nil")

	for _, test := range tests {
		batch, err := cli.NewBatch(test.queue)
		assert.Nil(t, err, "should be nil")

		batch.Enqueue(test.jobs...)

		err = batch.Execute()

		if test.shouldSuccess {
			assert.Nil(t, err, "should be nil")
			assert.Empty(t, batch.jobs, "should be emptied")

			assert.Contains(
				t,
				cli.redisClient.SMembers(
					fmt.Sprintf("%s:queues", cli.namespace),
				).Val(),
				test.queue,
				"resque queues should contain +\""+test.queue+"\"",
			)

			assert.Equal(
				t,
				int64(len(test.jobs)),
				cli.redisClient.LLen(
					fmt.Sprintf("%s:queue:%s", cli.namespace, test.queue),
				).Val(),
				"resque queue length should be the same of the amount of jobs",
			)
		} else {
			assert.NotNil(t, err, "should not be nil")
		}

		assert.Nil(t, cli.redisClient.FlushDb().Err(), "FLUSHDB should have succedeed")
	}
}
