package resque

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_NewBatch(t *testing.T) {
	tests := []struct {
		queue         string
		shouldSuccess bool
	}{
		{
			queue:         "test",
			shouldSuccess: true,
		},
		{
			queue:         "tutu",
			shouldSuccess: true,
		},
		{
			queue:         "tata",
			shouldSuccess: true,
		},
		{
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
		batch, err := cli.NewBatch(test.queue)
		if test.shouldSuccess {
			require.Nil(t, err, "should be nil")
			require.Equal(t, test.queue, batch.queue, "should be of same queue")
		} else {
			require.NotNil(t, err, "should not be nil")
		}
	}
}
