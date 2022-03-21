package redis

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hkloudou/mqx/face"
)

type redisAuther struct {
	// gOpt   *face.AuthOptions
	// face.AuthOptionConfiger
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
	// if err := obj.GlobalConfig(); err != nil {
	// 	return nil, err
	// }
	return obj, nil
}

func (m *redisAuther) Update(ctx context.Context, req *face.AuthRequest, options ...face.AuthRequestOption) error {
	var opts = face.DefaultAuthRequestOptions()
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return err
			}
		}
	}
	if err := m.checkConfig(); err != nil {
		return err
	}
	key := m.parseKeyFromReq(req, m.opts.authTmpl)
	model := &tokenModel{
		CreateAt:      uint64(time.Now().UnixNano()),
		TokenPassword: req.PassWord,
	}
	ttl := time.Duration(0)
	if opts.UseTtl {
		ttl = opts.Ttl
	}
	// set and check expired
	if ttl >= 0 {
		err := m.opts.client.Set(ctx, key, model, ttl).Err()
		if err != nil {
			return err
		}
		err = m.expiredBeforeConnection(req, opts.MaxTokens, opts.Discard)
		if err != nil {
			return err
		}
	}
	return m.opts.client.Del(ctx, key).Err()
}

func (m *redisAuther) Check(ctx context.Context, req *face.AuthRequest, options ...face.AuthRequestOption) (bool, error) {
	var opts = face.DefaultAuthRequestOptions()
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return false, err
			}
		}
	}
	if err := m.checkConfig(); err != nil {
		return false, face.ErrAuthServiceUnviable
	}
	var r *redis.StringCmd
	if !opts.UseTtl {
		r = m.opts.client.Get(ctx, m.parseKeyFromReq(req, m.opts.authTmpl))
	} else {
		r = m.opts.client.GetEx(ctx, m.parseKeyFromReq(req, m.opts.authTmpl), opts.Ttl)
	}

	if r.Err() != nil {
		if r.Err().Error() == redis.Nil.Error() {
			return false, face.ErrAuthInvalidUserNamePassword
		}
		return false, face.ErrAuthServiceUnviable
	}

	var obj tokenModel
	err := r.Scan(&obj)
	if err != nil {
		return false, err
	}
	if obj.TokenPassword != req.PassWord {
		return false, face.ErrAuthServiceUnviable
	}
	err = m.expiredBeforeConnection(req, opts.MaxTokens, opts.Discard)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *redisAuther) checkConfig() error {
	if m.opts == nil {
		return errors.New("MQX: auth.redis please init struct from New()")
	}
	if m.opts.client == nil {
		return errors.New("MQX: auth.redis please init struct from New()")
	}
	return nil
}

func (m *redisAuther) expiredBeforeConnection(req *face.AuthRequest, maxTokens uint64, discartPlolicy face.AuthDiscardPolicy) error {
	if maxTokens <= 0 {
		return nil
	}
	if err := m.checkConfig(); err != nil {
		return face.ErrAuthServiceUnviable
	}
	r := m.opts.client.Keys(context.TODO(), m.parseKeyFromReq(req, m.opts.listTmpl))
	if r.Err() != nil {
		return r.Err()
	}
	keys := r.Val()
	if discartPlolicy == face.AuthDiscardNew {
		if len(keys) >= int(maxTokens) {
			return face.ErrAuthInvalidTooManyTokens
		}
		return nil
	}
	items := make(tokenModelByCreateAd, 0)
	for i := 0; i < len(keys); i++ {
		var item tokenModel
		if v := m.opts.client.Get(context.TODO(), keys[i]); v.Err() != nil || v.Scan(&item) != nil {
			return face.ErrAuthServiceUnviable
		}
		item.Key = keys[i]
		items = append(items, item)
	}
	sort.Sort(sort.Reverse(items))
	keeped := uint64(0)
	discardedKeys := make([]string, 0)
	for i := 0; i < len(items); i++ {
		if keeped >= maxTokens {
			discardedKeys = append(discardedKeys, items[i].Key)
		}
		keeped++
	}
	if len(discardedKeys) > 0 {
		return m.opts.client.Del(context.TODO(), discardedKeys...).Err()
	}
	return nil
}
