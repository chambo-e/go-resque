package resque

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.Nil(t, err, "should be nil")
	require.NotNil(t, cli, "should not be nil")

	for _, test := range tests {
		batch, err := cli.NewBatch(test.queue)
		require.Nil(t, err, "should be nil")

		batch.Enqueue(test.jobs...)

		require.Equal(t, test.queue, batch.queue, "should be of same queue")
		require.Equal(t, len(test.jobs), len(batch.jobs), "should be of same size")
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
	require.Nil(t, err, "should be nil")
	require.NotNil(t, cli, "should not be nil")

	for _, test := range tests {
		batch, err := cli.NewBatch(test.queue)
		require.Nil(t, err, "should be nil")

		batch.Enqueue(test.jobs...)

		err = batch.Execute()

		if test.shouldSuccess {
			require.Nil(t, err, "should be nil")
			require.Empty(t, batch.jobs, "should be emptied")

			require.Contains(
				t,
				cli.redisClient.SMembers(
					fmt.Sprintf("%s:queues", cli.namespace),
				).Val(),
				test.queue,
				"resque queues should contain +\""+test.queue+"\"",
			)

			require.Equal(
				t,
				int64(len(test.jobs)),
				cli.redisClient.LLen(
					fmt.Sprintf("%s:queue:%s", cli.namespace, test.queue),
				).Val(),
				"resque queue length should be the same of the amount of jobs",
			)
		} else {
			require.NotNil(t, err, "should not be nil")
		}

		require.Nil(t, cli.redisClient.FlushDb().Err(), "FLUSHDB should have succedeed")
	}
}
