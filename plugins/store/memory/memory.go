package memory

import (
	"sync"

	"github.com/hkloudou/mqx/face"
)

type memoryStore struct {
	_map sync.Map
}

func MustNew(face.Conf) face.Store {
	return &memoryStore{
		_map: sync.Map{},
	}
}

func (m *memoryStore) Subscribe(clientid string, pattern []string, maxQos uint8) {

}
func (m *memoryStore) UnSubscribe(clientid string, pattern []string) {

}
func (m *memoryStore) Clear(clientid string) {
	// _topics := m.topic()
	// for _, _topic := range _topics {
	// 	// m._map
	// 	// list.List
	// }
}

func (m *memoryStore) topic() []string {
	_topics := make([]string, 0)
	m._map.Range(func(key, value any) bool {
		_topics = append(_topics, key.(string))
		return true
	})
	return _topics
}
