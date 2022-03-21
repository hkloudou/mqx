package redis

import (
	"strings"

	"github.com/go-redis/redis/v8"
)

type Option func(*Options) error

type Options struct {
	prefix   string
	client   *redis.Client
	authTmpl string
	listTmpl string
}

func WithRedisClient(options *redis.Options) Option {
	return func(o *Options) error {
		o.client = redis.NewClient(options)
		return nil
	}
}

func WithPrefix(prefix string) Option {
	return func(o *Options) error {
		o.prefix = strings.Trim(strings.Trim(prefix, " "), ".")
		return nil
	}
}
