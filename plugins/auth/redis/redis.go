package redis

import (
	"context"
	"encoding/base64"
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
	obj := &redisAuther{
		opts: &opts,
	}
	if err := obj.GlobalConfig(); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *redisAuther) Update(ctx context.Context, req *face.AuthRequest, ttl time.Duration) error {
	if err := m.checkConfig(); err != nil {
		return err
	}
	return m.opts.client.Set(ctx, m.parseKeyFromReq(req, m.opts.authTmpl), m.encode([]byte(req.PassWord)), ttl).Err()
}

func (m *redisAuther) encode(datas []byte) string {
	var _s = make([]byte, len(datas))
	copy(_s, datas)
	for i := 0; i < len(_s); i++ {
		_s[i] = _s[i] ^ byte(i&0xFF)
	}
	return base64.RawURLEncoding.EncodeToString(_s)
}

// func (m *redisAuther) decode(datas []byte) []byte {
// 	return m.encode(datas)
// }

func (m *redisAuther) Check(ctx context.Context, req *face.AuthRequest, ttl time.Duration) (bool, error) {
	if err := m.checkConfig(); err != nil {
		return false, err
	}
	var r *redis.StringCmd
	if ttl == 0 {
		r = m.opts.client.Get(ctx, m.parseKeyFromReq(req, m.opts.authTmpl))
	} else {
		r = m.opts.client.GetEx(ctx, m.parseKeyFromReq(req, m.opts.authTmpl), ttl)
	}

	if r.Err() != nil {
		return false, r.Err()
	}
	if r.Val() != m.encode([]byte(req.PassWord)) {
		return false, face.ErrAuthInvalidPassword
	}
	return true, nil
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
