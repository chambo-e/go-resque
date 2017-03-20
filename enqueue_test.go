package resque

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Enqueue(t *testing.T) {
	tests := []struct {
		job           Job
		queue         string
		shouldSuccess bool
	}{
		{
			job:           NewJob("test", "hello", "world"),
			queue:         "test",
			shouldSuccess: true,
		},
		{
			job:           NewJob("tutu", 1, "world"),
			queue:         "tutu",
			shouldSuccess: true,
		},
		{
			job:           NewJob("tutu", 1, "world"),
			queue:         "",
			shouldSuccess: false,
		},
	}

	cli, err := New(Configuration{
		RedisURI: "redis://localhost:6379",
	})
	require.Nil(t, err, "should be nil")
	require.NotNil(t, cli, "should not be nil")

	for _, test := range tests {
		err := cli.Enqueue(test.queue, test.job)

		if test.shouldSuccess {
			require.Nil(t, err, "should be nil")

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
				int64(1),
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
