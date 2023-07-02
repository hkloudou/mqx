package memory

import (
	"github.com/hkloudou/xtransport/packets/mqtt"
)

type config struct {
	// Server              string
	// PublishKey          string
	// PublishJetstreamKey string
	// publish_steam_key = "mqtt.bridge.stpublish"
}

type memoryBridge struct {
	// models []model
	cfg config
	ch  chan *mqtt.PublishPacket
	cb  func(obj *mqtt.PublishPacket)
	// st   *nrpc.Stream
	// conn *nats.Conn
	// js   nats.JetStreamContext
}

func (m *memoryBridge) Motion(cb func(obj *mqtt.PublishPacket)) error {
	// err := m.motion(cb)
	m.cb = cb
	// if err != nil {
	// 	return err
	// }
	// err = m.motionJs(cb)
	return nil
}

func (m *memoryBridge) Publish(pack *mqtt.PublishPacket) error {
	// err := m.ensureStream(m.cfg.PublishJetstreamKey)
	// if err != nil {
	// 	return err
	// }
	// _m := nats.NewMsg(m.cfg.PublishJetstreamKey)
	// var buf bytes.Buffer
	// _, err = pack.WriteTo(&buf)
	// if err != nil {
	// 	return err
	// }
	// _m.Data = buf.Bytes()
	// _, err = m.st.Js().PublishMsg(_m)
	go func() {
		m.ch <- pack
	}()
	return nil
}

func (m *memoryBridge) motion(cb func(obj *mqtt.PublishPacket)) error {
	// log.Println("ready motion publish in", m.cfg.PublishKey)
	// queue := base64.RawURLEncoding.EncodeToString([]byte(xruntime.HostName() + "-mqxbridge"))
	// _, err := m.conn.QueueSubscribe(m.cfg.PublishKey, queue, func(msg *nats.Msg) {
	// 	pk, err := mqtt.ReadPacket(bytes.NewReader(msg.Data))
	// 	// log.Println("da", pk.String())
	// 	if err != nil {
	// 		msg.Respond([]byte(err.Error()))
	// 		return
	// 	}
	// 	if pk.Type() != mqtt.Publish {
	// 		msg.Respond([]byte(fmt.Sprintf("can't pushlish a non mqtt.publish package")))
	// 		return
	// 	}
	// 	cb(pk.(*mqtt.PublishPacket))
	// })
	for {
		tmp := <-m.ch
		if m.cb != nil {
			m.cb(tmp)
		}
	}
	return nil
}

// func (m *memoryBridge) ensureStream(key string) error {
// 	return m.st.Config(&nats.StreamConfig{
// 		Name:         strings.ReplaceAll("mqtt.bridge.publish.jetstream", ".", "_"),
// 		Subjects:     []string{key},
// 		MaxMsgs:      -1,
// 		MaxAge:       5 * time.Second,
// 		Storage:      nats.MemoryStorage,
// 		NoAck:        true,
// 		MaxConsumers: -1,
// 	})
// }
// func (m *memoryBridge) motionJs(cb func(obj *mqtt.PublishPacket)) error {
// 	log.Println("ready motion jetstream in", m.cfg.PublishJetstreamKey)
// 	if err := m.ensureStream(m.cfg.PublishJetstreamKey); err != nil {
// 		panic(err)
// 	}
// 	queue := base64.RawURLEncoding.EncodeToString([]byte(xruntime.HostName() + "-mqxbridge"))
// 	_, err := m.st.Js().QueueSubscribe(m.cfg.PublishJetstreamKey, queue, func(msg *nats.Msg) {
// 		pk, err := mqtt.ReadPacket(bytes.NewReader(msg.Data))
// 		// log.Println("JS data arrve", pk.String())
// 		if err != nil {
// 			msg.Respond([]byte(err.Error()))
// 			return
// 		}
// 		if pk.Type() != mqtt.Publish {
// 			msg.Respond([]byte(fmt.Sprintf("can't pushlish js a non mqtt.publish package")))
// 			return
// 		}

// 		cb(pk.(*mqtt.PublishPacket))
// 	})
// 	return err
// }

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
