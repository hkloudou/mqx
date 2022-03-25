package memory

import (
	"context"
	"sync"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

type memoryRetainer struct {
	datas sync.Map
}

func init() {
	face.RegisterPugin("memory", MustNew)
}

func MustNew(conf face.Conf) face.Retain {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Retain, error) {
	obj := &memoryRetainer{datas: sync.Map{}}
	return obj, nil
}

func (m *memoryRetainer) Store(ctx context.Context, data *mqtt.PublishPacket) error {
	if err := face.ValidateTopic(data.TopicName); err != nil {
		return err
	}
	if len(data.Payload) == 0 {
		m.datas.Delete(data.TopicName)
		return nil
	}
	m.datas.Store(data.TopicName, data)
	return nil
}

func (m *memoryRetainer) Check(ctx context.Context, pattern string) ([]*mqtt.PublishPacket, error) {
	msgs := make([]*mqtt.PublishPacket, 0)
	if err := face.ValidatePattern(pattern); err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	m.datas.Range(func(key, value any) bool {
		keys = append(keys, key.(string))
		return true
	})

	for i := 0; i < len(keys); i++ {
		if face.MatchTopic(pattern, keys[i]) == nil {
			// matched2 = append(matched2, keys[i])
			if data, found := m.datas.Load(keys[i]); found {
				msgs = append(msgs, data.(*mqtt.PublishPacket))
			}
		}
	}
	return msgs, nil
}

func (m *memoryRetainer) Keys(ctx context.Context) ([]string, error) {
	keys := make([]string, 0)
	m.datas.Range(func(key, value any) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}
