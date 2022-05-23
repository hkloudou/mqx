package face

import (
	"context"

	"github.com/hkloudou/xtransport/packets/mqtt"
)

/*
	Acl - Access Control List
*/
type AclCode uint8

const (
	AclCodeAllow AclCode = iota
	AclCodeDeny
	AclCodeIgnore
	AclCodeError
)

type Acl interface {
	// context.Context removed
	GetSub(meta *MetaInfo, pattern string) AclCode
	GetPub(meta *MetaInfo, topic string) AclCode
}

type Retain interface {
	Store(ctx context.Context, data *mqtt.PublishPacket) error
	Check(ctx context.Context, pattern string) ([]*mqtt.PublishPacket, error)
	Keys(ctx context.Context) ([]string, error)
}

type Session interface {
	Add(ctx context.Context, clientid string, patterns ...string) error
	Remove(ctx context.Context, clientid string, patterns ...string) error
	Clear(ctx context.Context, cliendid string) error
	Match(ctx context.Context, topic string) ([]string, error)
	ClientPatterns(ctx context.Context, cliendid string) ([]string, error)
}

type Bridge interface {
	Publish(pack *mqtt.PublishPacket) error
	Motion(cb func(obj *mqtt.PublishPacket)) error
}
