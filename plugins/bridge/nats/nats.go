package nats

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hkloudou/nrpc"
	"github.com/hkloudou/xlib/xruntime"
	"github.com/hkloudou/xtransport/packets/mqtt"
	"github.com/nats-io/nats.go"
)

type config struct {
	Server              string
	PublishKey          string
	PublishJetstreamKey string
	// publish_steam_key = "mqtt.bridge.stpublish"
}

type natsBridge struct {
	// models []model
	cfg  config
	st   *nrpc.Stream
	conn *nats.Conn
	// js   nats.JetStreamContext
}

func (m *natsBridge) Motion(cb func(obj *mqtt.PublishPacket)) error {
	if err := m.motion(m.cfg.PublishKey, cb); err != nil {
		return err
	} else if err := m.motionJs(m.cfg.PublishJetstreamKey, cb); err != nil {
		return err
	}
	return nil
}

func (m *natsBridge) motion(key string, cb func(obj *mqtt.PublishPacket)) error {
	log.Println("ready motion publish in", key)
	queue := base64.RawURLEncoding.EncodeToString([]byte(xruntime.HostName() + "-mqxbridge"))
	_, err := m.conn.QueueSubscribe(key, queue, func(msg *nats.Msg) {
		pk, err := mqtt.ReadPacket(bytes.NewReader(msg.Data))
		log.Println("data arrve", pk.String())
		if err != nil {
			msg.Respond([]byte(err.Error()))
			return
		}
		if pk.Type() != mqtt.Publish {
			msg.Respond([]byte(fmt.Sprintf("can't pushlish a non mqtt.publish package")))
			return
		}
		cb(pk.(*mqtt.PublishPacket))
	})
	return err
}
func (m *natsBridge) motionJs(key string, cb func(obj *mqtt.PublishPacket)) error {
	log.Println("ready motion jetstream in", key)
	if err := m.st.Config(&nats.StreamConfig{
		Name:         strings.ReplaceAll("mqtt.bridge.publish.jetstream", ".", "_"),
		Subjects:     []string{key},
		MaxMsgs:      -1,
		MaxAge:       5 * time.Second,
		Storage:      nats.MemoryStorage,
		NoAck:        true,
		MaxConsumers: -1,
	}); err != nil {
		panic(err)
	}
	queue := base64.RawURLEncoding.EncodeToString([]byte(xruntime.HostName() + "-mqxbridge"))
	_, err := m.st.Js().QueueSubscribe(key, queue, func(msg *nats.Msg) {
		pk, err := mqtt.ReadPacket(bytes.NewReader(msg.Data))
		log.Println("js data arrve", pk.String())
		if err != nil {
			msg.Respond([]byte(err.Error()))
			return
		}
		if pk.Type() != mqtt.Publish {
			msg.Respond([]byte(fmt.Sprintf("can't pushlish js a non mqtt.publish package")))
			return
		}

		cb(pk.(*mqtt.PublishPacket))
	})
	return err
}

// func (m *natsBridge) motion(key string, stream bool, cb func(obj *mqtt.PublishPacket)) error {
// 	log.Println("motion", key, stream)
// 	// m.conn.Drain()
// 	queue := base64.RawURLEncoding.EncodeToString([]byte(xruntime.HostName() + "-mqxbridge"))
// 	// os.Hostname()
// 	// log.Println("queue", queue)
// 	readHead := func(pk *mqtt.PublishPacket, msg *nats.Msg) {
// 		if strings.ToLower(msg.Header.Get("x-mqtt-retain")) == "true" {
// 			// log.Println("read retain", pk.TopicName)
// 			pk.Retain = true
// 		} else {
// 			pk.Retain = false
// 		}
// 	}
// 	if stream {
// 		// m.js
// 		// log.Println("delete", m.js.)
// 		_, err := m.js.QueueSubscribe(key, queue, func(msg *nats.Msg) {
// 			// log.Println("data", msg)
// 			pk := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
// 			pk.Payload = msg.Data
// 			pk.TopicName = msg.Subject
// 			readHead(pk, msg)
// 			cb(pk)
// 		})
// 		return err
// 	} else {
// 		_, err := m.conn.QueueSubscribe(key, queue, func(msg *nats.Msg) {
// 			// log.Println("data arrive", msg)
// 			pk := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
// 			pk.Payload = msg.Data
// 			pk.TopicName = msg.Subject
// 			readHead(pk, msg)
// 			cb(pk)
// 		})
// 		return err
// 	}
// 	return nil
// }
