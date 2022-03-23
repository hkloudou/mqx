package redis

import (
	"strings"

	"github.com/go-redis/redis/v8"
)

type Option func(*Options) error

type Options struct {
	prefix string
	addr   string
	db     int
	client *redis.Client
}

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
		trims := []string{" ", ":", "."}
		o.prefix = prefix
		for i := 0; i < len(trims); i++ {
			o.prefix = strings.Trim(o.prefix, trims[i])
		}
		return nil
	}
}
