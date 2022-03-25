package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

type memoryAuth struct {
	authTmpl  string `ini:"-"`
	listTmpl  string `ini:"-"`
	prefix    string `ini:"-"`
	public    modelPublic
	ttl       time.Duration
	discard   uint
	maxTokens uint64
}

type modelPublic struct {
	Enable   bool
	Username string
	Password string
}

func init() {
	face.RegisterPugin[face.Auth]("memory", MustNew)
	// face.DefaultAuths["redis"] = MustNew
}

func MustNew(conf face.Conf) face.Auth {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Auth, error) {
	obj := &memoryAuth{
		authTmpl: "$p/$u/$c",
		listTmpl: "$p/$u/*",
		public: modelPublic{
			Enable: false,
		},
	}
	if conf == nil {
		return nil, fmt.Errorf("Invalid conf")
	}
	if err := conf.MapTo("auth.public", &obj.public); err != nil {
		return nil, err
	}
	// ignore config
	// default maxTokens and ttl
	obj.maxTokens = uint64(conf.MustUint("auth", "max_tokens", uint(obj.maxTokens)))
	obj.ttl = conf.MustDuration("auth", "ttl", obj.ttl)
	obj.discard = conf.MustUint("auth", "discard", 0)
	if obj.discard > 1 { //only 0,1 is valid value
		obj.discard = 0
	}
	return obj, nil
}

func (m *memoryAuth) Update(ctx context.Context, req *face.AuthRequest, options ...face.AuthRequestOption) error {
	return nil
}

func (m *memoryAuth) Check(ctx context.Context, req *face.AuthRequest, options ...face.AuthRequestOption) mqtt.ConnackReturnCode {

	// public account
	if m.public.Enable && m.public.Username == req.UserName && m.public.Password == req.PassWord {
		return mqtt.Accepted
	}
	return mqtt.ErrRefusedServerUnavailable
}

func (m *memoryAuth) MotionExpired(fc func(userName, clientId string) error) error {
	for {
		time.Sleep(1 * time.Second)
	}
}
