package redis

import (
	"strings"

	"github.com/go-redis/redis/v8"
)

type Option func(*Options) error

type Options struct {
	prefix   string
	addr     string
	db       int
	client   *redis.Client
	authTmpl string
	listTmpl string
}

// func WithRedisClient(options *redis.Options) Option {
// 	return func(o *Options) error {
// 		o.client = redis.NewClient(&redis.Options{})
// 		return nil
// 	}
// }

func WithRedisAddr(addr string) Option {
	return func(o *Options) error {
		o.addr = addr
		return nil
	}
}
func WithRedisDB(db int) Option {
	return func(o *Options) error {
		o.db = db
		return nil
	}
}

func WithPrefix(prefix string) Option {
	return func(o *Options) error {
		o.prefix = strings.Trim(strings.Trim(prefix, " "), ".")
		return nil
	}
}
