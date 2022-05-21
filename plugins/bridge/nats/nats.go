package nats

import (
	"github.com/hkloudou/xtransport/packets/mqtt"
	"github.com/nats-io/nats.go"
)

type config struct {
	Server string
}

type natsBridge struct {
	models []model
	cfg    config
	conn   *nats.Conn
	js     nats.JetStreamContext
}

func (m *natsBridge) Motion(cb func(obj *mqtt.PublishPacket)) {
	for i := 0; i < len(m.models); i++ {
		if m.models[i].Stream {
			patterns := m.models[i].Patterns
			for j := 0; j < len(patterns); j++ {
				m.motion(patterns[j], cb)
			}

		}
	}
}

func (m *natsBridge) motion(key string, cb func(obj *mqtt.PublishPacket)) {
	m.js.QueueSubscribe(key, "mqx.bridge", func(msg *nats.Msg) {
		pk := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
		pk.Payload = msg.Data
		pk.TopicName = msg.Subject
	})
}
