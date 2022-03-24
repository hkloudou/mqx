package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

type redisAuther struct {
	conf   model
	client *redis.Client
}

// [auth.redis]
// server = 127.0.0.1:6379
// pool = 8
// database = 3
// password =
// salt = dbsalt
type model struct {
	Server   string
	Pool     uint16
	Db       uint16
	Username string
	Password string
	Salt     string
	authTmpl string `ini:"-"`
	listTmpl string `ini:"-"`
	prefix   string `ini:"-"`
}

func MustNew(conf face.Conf) face.Auth {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Auth, error) {
	obj := &redisAuther{conf: model{
		prefix:   "mqtt.auth",
		authTmpl: "$p/$u/$c",
		listTmpl: "$p/$u/*",
		Server:   "127.0.0.1:6379",
		Db:       3,
	}}
	if conf != nil {
		if err := conf.MapTo("auth.plugin.redis", &obj.conf); err != nil {
			return nil, err
		}
	}
	obj.client = redis.NewClient(&redis.Options{
		Addr:     obj.conf.Server,
		Password: obj.conf.Password,
		Username: obj.conf.Username,
		DB:       int(obj.conf.Db),
	})
	return obj, nil
}

func (m *redisAuther) Update(ctx context.Context, req *face.AuthRequest, options ...face.AuthRequestOption) error {
	/*
		http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.pdf
		   The ClientId MUST be a UTF-8 encoded string as defined in Section 1.5.3 [MQTT-3.1.3-4].
		   578
		   579 The Server MUST allow ClientIds which are between 1 and 23 UTF-8 encoded bytes in length, and that
		   580 contain only the characters
		   581 "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" [
	*/

	if b, err := regexp.MatchString(`^[A-Za-z0-9-_]{1,36}$`, req.ClientId); !b || err != nil {
		return face.ErrAuthInvalidClientId
	}
	if b, err := regexp.MatchString(`^[A-Za-z0-9-_@]{1,36}$`, req.UserName); !b || err != nil {
		return face.ErrAuthInvalidUserNamePassword
	}

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
	//check Expired Before Auth if discard policy is discard New
	if opts.Discard == face.AuthDiscardNew {
		err := m.expiredBeforeConnection(req, opts.MaxTokens, opts.Discard)
		if err != nil {
			return err
		}
	}

	key := m.parseKeyFromReq(req, m.conf.authTmpl)
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
		err := m.client.Set(ctx, key, model, ttl).Err()
		if err != nil {
			return err
		}
		// if discardold
		if opts.Discard == face.AuthDiscardOld {
			err = m.expiredBeforeConnection(req, opts.MaxTokens, opts.Discard)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// return nil
	return m.client.Del(ctx, key).Err()
}

func (m *redisAuther) Check(ctx context.Context, req *face.AuthRequest, options ...face.AuthRequestOption) mqtt.ConnackReturnCode {
	var opts = face.DefaultAuthRequestOptions()
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return mqtt.ErrRefusedServerUnavailable
			}
		}
	}
	if err := m.checkConfig(); err != nil {
		return mqtt.ErrRefusedServerUnavailable
	}
	var r *redis.StringCmd
	if !opts.UseTtl {
		r = m.client.Get(ctx, m.parseKeyFromReq(req, m.conf.authTmpl))
	} else {
		r = m.client.GetEx(ctx, m.parseKeyFromReq(req, m.conf.authTmpl), opts.Ttl)
	}

	if r.Err() != nil {
		if r.Err().Error() == redis.Nil.Error() {
			return mqtt.ErrRefusedBadUsernameOrPassword
		}
		return mqtt.ErrRefusedServerUnavailable
	}

	var obj tokenModel
	err := r.Scan(&obj)
	if err != nil {
		return mqtt.ErrRefusedServerUnavailable
	}
	if obj.TokenPassword != req.PassWord {
		return mqtt.ErrRefusedBadUsernameOrPassword
	}
	err = m.expiredBeforeConnection(req, opts.MaxTokens, opts.Discard)
	if err != nil {
		return mqtt.ErrRefusedServerUnavailable
	}
	return mqtt.Accepted
}

func (m *redisAuther) checkConfig() error {
	// if m.opts == nil {
	// 	return errors.New("MQX: auth.redis please init struct from New()")
	// }
	if m.client == nil {
		return errors.New("MQX: auth.redis please init struct from New()")
	}
	return nil
}

func (m *redisAuther) expiredBeforeConnection(req *face.AuthRequest, maxTokens uint64, discartPlolicy face.AuthDiscardPolicy) error {
	if maxTokens <= 0 {
		return nil
	}
	r := m.client.Keys(context.TODO(), m.parseKeyFromReq(req, m.conf.listTmpl))
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
		if v := m.client.Get(context.TODO(), keys[i]); v.Err() != nil || v.Scan(&item) != nil {
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
		return m.client.Del(context.TODO(), discardedKeys...).Err()
	}
	return nil
}

func (m *redisAuther) MotionExpired(fc func(userName, clientId string) error) error {
	for {
		err := m.motionExpired(fc)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (m *redisAuther) motionExpired(fc func(userName, clientId string) error) (reerr error) {
	defer func() {
		if r := recover(); r != nil {
			reerr = fmt.Errorf("%v", r)
		}
	}()
	if err := m.client.ConfigSet(context.TODO(), "notify-keyspace-events", "$Kxeg").Err(); err != nil {
		return err
	}
	//"+m.opts.prefix+".*"
	str := fmt.Sprintf("__keyspace@%d__:"+m.conf.prefix+"/*", m.client.Options().DB)

	// log.Println("motion", str)
	pubsub := m.client.PSubscribe(context.TODO(),
		str,
	)
	defer pubsub.Close()
	for {
		data, err := pubsub.ReceiveMessage(context.TODO())
		if err != nil {
			panic(err)
		}
		// log.Println("data.Pattern", data.Pattern, "data.Channel", data.Channel, "p", data.Payload)
		if string(data.Payload) == "set" || string(data.Payload) == "del" || string(data.Payload) == "expired" || string(data.Payload) == "evict" {
			go func() {
				for i := 0; i < 10; i++ {
					remain := strings.TrimPrefix(data.Channel, fmt.Sprintf("__keyspace@%d__:"+m.conf.prefix+".", m.client.Options().DB))
					// println("remain", remain)
					arr := strings.Split(remain, "/")
					if len(arr) != 2 {
						return
					}
					if fc(arr[0], arr[1]) == nil {
						break
					}
					time.Sleep(1 * time.Second)
				}
			}()
		}
	}
}
