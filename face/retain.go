package face

import (
	"context"

	"github.com/hkloudou/xtransport/packets/mqtt"
)

type Retain interface {
	Watch(func(data *mqtt.PublishPacket))
	Store(ctx context.Context, data *mqtt.PublishPacket) error
	Check(ctx context.Context, pattern string) (*mqtt.PublishPacket, error)
	Keys(ctx context.Context) ([]string, error)
}
