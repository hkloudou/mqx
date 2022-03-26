package redis

import (
	"context"
	"fmt"
	"strings"
	"unsafe"

	"github.com/go-redis/redis/v8"
	"github.com/hkloudou/mqx/face"
)

type keyType uint8

const (
	ClientKeyType keyType = iota
	TopicKeyType
)

type model struct {
	Server string
	// Pool     uint16
	Db             uint16
	Username       string
	Password       string
	Salt           string
	clientPrefix   string `ini:"-"`
	pubTopicPrefix string `ini:"-"`
	priTopicPrefix string `ini:"-"`
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
			clientPrefix:   "mqtt.sess/cli",
			pubTopicPrefix: "mqtt.sess/pub",
			priTopicPrefix: "mqtt.sess/pri",
			Server:         "127.0.0.1:6379",
			Db:             3,
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
	keys := make([]interface{}, 0)
	for i := 0; i < len(patterns); i++ {
		keys = append(keys, patterns[i])
		if err := m.client.SAdd(ctx, m.toRedisKey(TopicKeyType, patterns[i]), clientid).Err(); err != nil {
			return err
		}
	}
	err := m.client.SAdd(ctx, m.toRedisKey(ClientKeyType, clientid), keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *redisSessioner) Remove(ctx context.Context, clientid string, patterns ...string) error {
	keys := make([]interface{}, 0)
	for i := 0; i < len(patterns); i++ {
		keys = append(keys, patterns[i])
		if err := m.client.SRem(ctx, m.toRedisKey(TopicKeyType, patterns[i]), clientid).Err(); err != nil {
			return err
		}
	}
	err := m.client.SRem(ctx, m.toRedisKey(ClientKeyType, clientid), keys...).Err()
	if err != nil {
		return err
	}
	return nil
}
func (m *redisSessioner) Match(ctx context.Context, topic string) ([]string, error) {
	redisPatternKey, err := m.buildRedisSearchPattern(topic)
	// log.Println("search", redisPatternKey)
	if err != nil {
		return nil, err
	}
	r := m.client.Keys(ctx, redisPatternKey)
	if r.Err() != nil {
		return nil, r.Err()
	}
	patterns := r.Val()
	matched := make([]string, 0)
	for i := 0; i < len(patterns); i++ {
		if face.MatchTopic(m.toTopic(patterns[i]), topic) == nil {
			matched = append(matched, patterns[i])
		}
	}

	lists := make(map[string]bool, 0)
	for i := 0; i < len(matched); i++ {
		r := m.client.SMembers(ctx, matched[i])
		if r.Err() != nil {
			// log.Println("err", r.Err())
			return nil, r.Err()
		}
		tmps := r.Val()
		for i := 0; i < len(tmps); i++ {
			lists[tmps[i]] = true
		}
	}
	keys := make([]string, 0, len(lists))
	for k := range lists {
		keys = append(keys, k)
	}
	return keys, nil
}
func (m *redisSessioner) Clear(ctx context.Context, clientid string) error {
	r := m.client.SMembers(ctx, m.toRedisKey(ClientKeyType, clientid))
	if r.Err() != nil {
		return r.Err()
	}
	patterns := r.Val()
	for i := 0; i < len(patterns); i++ {
		if err := m.client.SRem(ctx, m.toRedisKey(TopicKeyType, patterns[i]), clientid).Err(); err != nil {
			return err
		}
	}
	return m.client.Del(ctx, m.toRedisKey(ClientKeyType, clientid)).Err()
}

func (m *redisSessioner) getTopicPrefix(topic string) string {
	if face.IsPrivateTopic(topic) {
		return m.conf.priTopicPrefix
	}
	return m.conf.pubTopicPrefix
}

func (m *redisSessioner) toRedisKey(tp keyType, topic string) string {
	prefix := ""
	if tp == TopicKeyType {
		prefix = m.getTopicPrefix(topic)
	} else {
		prefix = m.conf.clientPrefix
	}

	if strings.HasPrefix(topic, prefix+"/") {
		return topic
	}
	return prefix + "/" + topic
}

func (m *redisSessioner) toTopic(key string) string {
	prefixs := []string{m.conf.pubTopicPrefix, m.conf.priTopicPrefix, m.conf.clientPrefix}
	for i := 0; i < len(prefixs); i++ {
		prefix := prefixs[i]
		if strings.HasPrefix(key, prefix+"/") {
			return strings.TrimPrefix(key, prefix+"/")
		}
	}
	return key
}

func (m *redisSessioner) buildRedisSearchPattern(topic string) (string, error) {
	// prefix := m.getTopicPrefix(topic)
	topicspec := strings.Split(topic, "/")

	redisPatternKey := ""
	if face.IsPrivateTopic(topicspec[0]) {
		if len(topicspec) < 3 {
			return "", fmt.Errorf("$usr,$cid topic should have topic  $usr/parame/topic")
		}
		redisPatternKey = m.conf.priTopicPrefix + "/" + topicspec[0] + "/" + topicspec[1] + "/*"
	} else {
		redisPatternKey = m.conf.pubTopicPrefix + "/*"
	}
	return redisPatternKey, nil
}

func (m *redisSessioner) ClientPatterns(ctx context.Context, cliendid string) ([]string, error) {
	// log.Println("c", m.conf.clientPrefix+"/"+cliendid+"/*")
	r := m.client.SMembers(ctx, m.conf.clientPrefix+"/"+cliendid)
	if r.Err() != nil {
		return nil, r.Err()
	}
	return r.Val(), nil
}
