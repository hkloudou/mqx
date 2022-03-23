package redis

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

type redisRetainer struct {
	// gOpt   *face.AuthOptions
	// face.AuthOptionConfiger
	opts *Options
}

func New(options ...Option) (face.Retain, error) {
	opts := Options{
		addr:   "localhost:6379",
		db:     3,
		prefix: "mqtt.retain",
	}
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
	}
	opts.client = redis.NewClient(&redis.Options{
		Addr:     opts.addr,
		Password: "", // no password set
		DB:       opts.db,
	})
	obj := &redisRetainer{
		opts: &opts,
	}
	return obj, nil
}

func (m *redisRetainer) Watch(func(data *mqtt.PublishPacket)) {}
func (m *redisRetainer) Store(ctx context.Context, data *mqtt.PublishPacket) error {
	if err := face.ValidatePublishTopic(data.TopicName); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := data.Write(&buf); err != nil {
		return err
	}
	_s := buf.Bytes()
	for i := 0; i < len(_s); i++ {
		_s[i] = _s[i] ^ byte(i&0xFF)
	}

	return m.opts.client.Set(ctx, m.opts.prefix+":"+data.TopicName, _s, 0).Err()
}

func (m *redisRetainer) Check(ctx context.Context, pattern string) (*mqtt.PublishPacket, error) {
	if err := face.ValidateTopicPattern(pattern); err != nil {
		return nil, err
	}
	regStr := pattern
	regStr = strings.ReplaceAll(regStr, "#", "*")
	regStr = strings.ReplaceAll(regStr, "+", "*")
	r := m.opts.client.Keys(ctx, m.opts.prefix+":"+regStr)
	if r.Err() != nil {
		return nil, r.Err()
	}
	matched := r.Val()
	if len(matched) == 0 {
		return nil, nil
	}
	r2 := m.opts.client.Get(ctx, matched[0])
	if r2.Err() != nil {
		return nil, r2.Err()
	}
	_s, err := r2.Bytes()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(_s); i++ {
		_s[i] = _s[i] ^ byte(i&0xFF)
	}
	obj, err := mqtt.ReadPacket(bytes.NewBuffer(_s))
	if err != nil {
		return nil, err
	}
	if obj2, ok := obj.(*mqtt.PublishPacket); !ok {
		return nil, errors.New("not found")
	} else {
		return obj2, nil
	}
}

func (m *redisRetainer) Keys(ctx context.Context) ([]string, error) {
	r := m.opts.client.Keys(context.TODO(), m.opts.prefix+":*")
	if r.Err() != nil {
		return nil, r.Err()
	}
	return r.Val(), nil
}
