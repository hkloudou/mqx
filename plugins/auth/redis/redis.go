package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hkloudou/mqx/face"
)

type token struct {
	// SET	mqtt.$u.$c
	Req *face.AuthReply
}

type redisAuther struct {
	// gOpt   *face.AuthOptions
	face.AuthOptionConfiger
	opts *Options
}

func New(options ...Option) (face.Auth, error) {
	opts := Options{
		prefix:   "mqtt.auth",
		authTmpl: "$p.$u.$c",
		listTmpl: "$p.$u.*",
	}
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
	}
	if opts.client == nil {
		opts.client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       1,
		})
	}

	return &redisAuther{
		opts: &opts,
	}, nil
}

func (m *redisAuther) Update(ctx context.Context, req *face.AuthRequest, ttl time.Duration) error {
	if err := m.checkConfig(); err != nil {
		return err
	}
	return m.opts.client.Set(ctx, m.parseKeyFromReq(req, m.opts.authTmpl), req.PassWord, ttl).Err()
}

func (m *redisAuther) checkConfig() error {
	if m.GOpt == nil {
		return errors.New("MQX: auth.redis please call Config()")
	}
	if m.opts == nil {
		return errors.New("MQX: auth.redis please init struct from New()")
	}
	if m.opts.client == nil {
		return errors.New("MQX: auth.redis please init struct from New()")
	}
	return nil
}
