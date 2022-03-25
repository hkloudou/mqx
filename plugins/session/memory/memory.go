package redis

import (
	"context"
	"sync"

	"github.com/fatih/set"
	"github.com/hkloudou/mqx/face"
)

type memoryRession struct {
	_lock         sync.RWMutex
	_clientTopics sync.Map
	_topicClients sync.Map
}

func init() {
	face.RegisterPugin("memory", MustNew)
}

func MustNew(conf face.Conf) face.Session {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Session, error) {
	obj := &memoryRession{
		_clientTopics: sync.Map{},
		_topicClients: sync.Map{},
	}
	return obj, nil
}

func (m *memoryRession) getTopics(clientid string) set.Interface {
	acture, _ := m._clientTopics.LoadOrStore(clientid, set.New(set.ThreadSafe))
	return acture.(set.Interface)
}

func (m *memoryRession) getClients(topic string) set.Interface {
	acture, _ := m._topicClients.LoadOrStore(topic, set.New(set.ThreadSafe))
	return acture.(set.Interface)
}

func (m *memoryRession) Add(ctx context.Context, clientid string, patterns ...string) error {
	tmp := m.getTopics(clientid)
	for _, topic := range patterns {
		tmp.Add(topic)
		m.getClients(topic).Add(clientid)
	}
	return nil
}

func (m *memoryRession) Remove(ctx context.Context, clientid string, patterns ...string) error {
	tmp := m.getTopics(clientid)
	for _, topic := range patterns {
		tmp.Remove(topic)
		m.getClients(topic).Remove(clientid)
	}
	return nil
}

func (m *memoryRession) Clear(ctx context.Context, clientid string) {
	tmp := m.getTopics(clientid)
	for _, topic := range tmp.List() {
		m.getClients(topic.(string)).Remove(clientid)
	}
	tmp.Clear()
}

func (m *memoryRession) Patterns(ctx context.Context) ([]string, error) {
	keys := make([]string, 0)
	m._topicClients.Range(func(key, value any) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}

func (m *memoryRession) Match(ctx context.Context, topic string) ([]string, error) {
	patterns, err := m.Patterns(ctx)
	if err != nil {
		return nil, err
	}
	matched := make([]string, 0)
	for i := 0; i < len(patterns); i++ {
		if face.MatchTopic(patterns[i], topic) == nil {
			matched = append(matched, patterns[i])
		}
	}
	lists := set.New(set.NonThreadSafe)
	for i := 0; i < len(matched); i++ {
		clis := m.getClients(matched[i])
		lists.Merge(clis)
	}
	keys := make([]string, 0)
	for _, v := range lists.List() {
		keys = append(keys, v.(string))
	}
	return keys, nil
}

// func (m *memoryRession) List(ctx context.Context, pattern string) ([]string, error) {
// 	ids := make([]string, 0)
// 	// for _, v := range m.getClients(pattern).List() {
// 	// 	ids = append(ids, v.(string))
// 	// }
// 	return ids, nil
// }
