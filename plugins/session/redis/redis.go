package redis

import (
	"context"
	"strings"
	"unsafe"

	"github.com/go-redis/redis/v8"
	"github.com/hkloudou/mqx/face"
)

type model struct {
	Server string
	// Pool     uint16
	Db       uint16
	Username string
	Password string
	Salt     string
	prefix   string `ini:"-"`
}

//https://github.com/go-redis/redis/blob/master/internal/util/unsafe.go
// stringToBytes converts string to byte slice.
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

type redisSessioner struct {
	conf   model
	client *redis.Client
	// opts *Options
}

func init() {
	face.RegisterPugin("redis", MustNew)
}

func MustNew(conf face.Conf) face.Session {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Session, error) {
	obj := &redisSessioner{
		conf: model{
			prefix: "mqtt.session",
			Server: "127.0.0.1:6379",
			Db:     3,
		},
	}
	if conf != nil {
		if err := conf.MapTo("session.plugin.redis", &obj.conf); err != nil {
			return nil, err
		}
	}
	obj.client = redis.NewClient(&redis.Options{
		Addr:     obj.conf.Server,
		Password: obj.conf.Password,
		Username: obj.conf.Username,
		DB:       int(obj.conf.Db),
	})
	if err := obj.client.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *redisSessioner) Add(ctx context.Context, clientid string, patterns ...string) error {
	// m.client.SAdd(ctx, clientid, patterns)
	return nil
}
func (m *redisSessioner) Match(ctx context.Context, pattern string) ([]string, error) {
	return nil, nil
}
func (m *redisSessioner) Clear(ctx context.Context, cliendid string) {}

// func (m *redisSessioner) Store(ctx context.Context, data *mqtt.PublishPacket) error {
// 	if err := face.ValidateTopic(data.TopicName); err != nil {
// 		return err
// 	}
// 	if len(data.Payload) == 0 {
// 		return m.client.Del(ctx, m.toRedisKey(data.TopicName)).Err()
// 	}
// 	var buf bytes.Buffer
// 	if err := data.Write(&buf); err != nil {
// 		return err
// 	}
// 	_s := buf.Bytes()
// 	for i := 0; i < len(_s); i++ {
// 		_s[i] = _s[i] ^ byte(i&0xFF)
// 	}
// 	return m.client.Set(ctx, m.toRedisKey(data.TopicName), _s, 0).Err()
// }

// func (m *redisSessioner) Check(ctx context.Context, pattern string) ([]*mqtt.PublishPacket, error) {
// 	msgs := make([]*mqtt.PublishPacket, 0)
// 	if err := face.ValidatePattern(pattern); err != nil {
// 		return nil, err
// 	}
// 	//fuzzy query redis keys
// 	r := m.client.Keys(ctx, m.toRedisKey(strings.ReplaceAll(strings.ReplaceAll(pattern, "#", "*"), "+", "*")))
// 	if r.Err() != nil {
// 		return nil, r.Err()
// 	}
// 	matched := r.Val()
// 	matched2 := make([]string, 0)
// 	for i := 0; i < len(matched); i++ {
// 		if face.MatchTopic(pattern, m.toTopic(matched[i])) == nil {
// 			matched2 = append(matched2, matched[i])
// 		}
// 	}
// 	if len(matched2) == 0 {
// 		return msgs, nil
// 	}

// 	r2 := m.client.MGet(ctx, matched2...)
// 	// r2 := m.opts.client.Get(ctx, m.opts.prefix+"/"+matched2[0])
// 	if r2.Err() != nil {
// 		return nil, r2.Err()
// 	}
// 	for _, v := range r2.Val() {
// 		_s := stringToBytes(v.(string))
// 		for i := 0; i < len(_s); i++ {
// 			_s[i] = _s[i] ^ byte(i&0xFF)
// 		}
// 		obj, err := mqtt.ReadPacket(bytes.NewBuffer(_s))
// 		if err != nil {
// 			return nil, err
// 		}
// 		if obj, ok := obj.(*mqtt.PublishPacket); ok {
// 			msgs = append(msgs, obj)
// 		}
// 	}
// 	return msgs, nil
// }

// func (m *redisSessioner) Keys(ctx context.Context) ([]string, error) {
// 	r := m.client.Keys(context.TODO(), m.conf.prefix+"/*")
// 	if r.Err() != nil {
// 		return nil, r.Err()
// 	}
// 	keys := make([]string, 0)
// 	for _, v := range r.Val() {
// 		keys = append(keys, m.toTopic(v))
// 	}
// 	return keys, nil
// }

func (m *redisSessioner) toRedisKey(topic string) string {
	if strings.HasPrefix(topic, m.conf.prefix+"/") {
		return topic
	}
	return m.conf.prefix + "/" + topic
}

func (m *redisSessioner) toTopic(key string) string {
	if strings.HasPrefix(key, m.conf.prefix+"/") {
		return strings.TrimPrefix(key, m.conf.prefix+"/")
	}
	return key
}
