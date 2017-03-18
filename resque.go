package resque

import redis "gopkg.in/redis.v5"

// Client holds an active resque connection
type Client struct {
	redisClient *redis.Client
	namespace   string
}

// RedisOptions is a wrapper of https://godoc.org/gopkg.in/redis.v5#Options
type RedisOptions redis.Options

// Configuration stores the required configuration to create a resque client
type Configuration struct {
	// Redis URI connection string
	RedisURI string

	Redis RedisOptions
	// Resque namespace (default: "resque")
	Namespace string
}

// New returns a new resque client.
// Goroutine safe
func New(cfg Configuration) (*Client, error) {
	var opts redis.Options

	if cfg.RedisURI == "" {
		opts = redis.Options(cfg.Redis)
	} else {
		optsNew, err := redis.ParseURL(cfg.RedisURI)
		if err != nil {
			return nil, err
		}

		opts = *optsNew
	}

	client := redis.NewClient(&opts)
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	if cfg.Namespace == "" {
		cfg.Namespace = "resque"
	}

	return &Client{
		redisClient: client,
		namespace:   cfg.Namespace,
	}, nil
}
