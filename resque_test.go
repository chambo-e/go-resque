package resque

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			assert.Nil(t, err, "should be nil")
			assert.NotNil(t, cli, "should not be nil")

			if test.cfg.Namespace != "" {
				assert.Equal(t, cli.namespace, test.cfg.Namespace, "should have set namespace")
			} else {
				assert.Equal(t, "resque", cli.namespace, "should have set default namespace")
			}
		} else {
			assert.NotNil(t, err, "should not be nil")
			assert.Nil(t, cli, "should be nil")
		}

	}
}
