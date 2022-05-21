package nats

import (
	"fmt"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/nrpc"
)

type model struct {
	Patterns []string `delim:","`
	Stream   bool
}

func init() {
	face.RegisterPugin[face.Bridge]("nats", MustNew)
	// face.DefaultAuths["redis"] = MustNew
}

func MustNew(conf face.Conf) face.Bridge {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Bridge, error) {
	obj := &natsBridge{}
	if conf == nil {
		return nil, fmt.Errorf("Invalid conf")
	}
	if err := conf.MapTo("bridge.plugin.nats", &obj.cfg); err != nil {
		return nil, err
	}
	conn, err := nrpc.Connect(obj.cfg.Server)
	if err != nil {
		return nil, err
	}
	obj.js, err = conn.JetStream()
	if err != nil {
		return nil, err
	}
	loopRead := func(fm string) []model {
		i := 0
		items := make([]model, 0)
		for {
			var item model
			if err := conf.MapTo(fmt.Sprintf(fm, i), &item); err != nil {
				break
			}
			if len(item.Patterns) == 0 {
				break
			}
			items = append(items, item)
			i++
		}
		return items
	}
	obj.models = make([]model, 0)
	obj.models = loopRead("bridge.plugin.motion.%d")
	return obj, nil
}
