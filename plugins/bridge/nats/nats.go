package nats

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/hkloudou/xlib/xruntime"
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

func (m *natsBridge) Motion(cb func(obj *mqtt.PublishPacket)) error {
	for i := 0; i < len(m.models); i++ {
		patterns := m.models[i].Patterns
		for j := 0; j < len(patterns); j++ {
			if err := m.motion(patterns[j], m.models[i].Stream, cb); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *natsBridge) motion(key string, stream bool, cb func(obj *mqtt.PublishPacket)) error {
	log.Println("motion", key, stream)
	// m.conn.Drain()
	queue := base64.RawURLEncoding.EncodeToString([]byte(xruntime.HostName() + "-mqxbridge"))
	// os.Hostname()
	log.Println("queue", queue)
	readHead := func(pk *mqtt.PublishPacket, msg *nats.Msg) {
		if strings.ToLower(msg.Header.Get("x-mqtt-retain")) == "true" {
			pk.Retain = true
		} else {
			pk.Retain = false
		}
	}
	if stream {
		// m.js
		// log.Println("delete", m.js.)
		_, err := m.js.QueueSubscribe(key, queue, func(msg *nats.Msg) {
			// log.Println("data", msg)
			pk := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
			pk.Payload = msg.Data
			pk.TopicName = msg.Subject
			readHead(pk, msg)
			cb(pk)
		})
		return err
		// log.Println("err", err)
		// log.Println(sub.Drain())
	} else {
		m.conn.QueueSubscribe(key, queue, func(msg *nats.Msg) {
			pk := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
			pk.Payload = msg.Data
			pk.TopicName = msg.Subject
			readHead(pk, msg)
			cb(pk)
		})
	}
	return nil
}
