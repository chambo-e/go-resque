package resque

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		cfg           Configuration
		shouldSuccess bool
	}{
		{
			cfg:           Configuration{},
			shouldSuccess: false,
		},
		{
			cfg: Configuration{
				RedisURI: "lolo",
			},
			shouldSuccess: false,
		},
		{
			cfg: Configuration{
				RedisURI: "redis://localhost:6379",
			},
			shouldSuccess: true,
		},
		{
			cfg: Configuration{
				RedisURI:  "redis://localhost:6379",
				Namespace: "test",
			},
			shouldSuccess: true,
		},
		{
			cfg: Configuration{
				Redis: RedisOptions{
					Addr:        "localhost:6379",
					DB:          1,
					DialTimeout: time.Second,
				},
				Namespace: "test",
			},
			shouldSuccess: true,
		},
		{
			cfg: Configuration{
				Redis: RedisOptions{
					Addr:        "localhost:6378",
					DB:          1,
					DialTimeout: 100 * time.Millisecond,
				},
				Namespace: "test",
			},
			shouldSuccess: false,
		},
	}

	for _, test := range tests {
		cli, err := New(test.cfg)

		if test.shouldSuccess {
			require.Nil(t, err, "should be nil")
			require.NotNil(t, cli, "should not be nil")

			if test.cfg.Namespace != "" {
				require.Equal(t, cli.namespace, test.cfg.Namespace, "should have set namespace")
			} else {
				require.Equal(t, "resque", cli.namespace, "should have set default namespace")
			}
		} else {
			require.NotNil(t, err, "should not be nil")
			require.Nil(t, cli, "should be nil")
		}

	}
}
