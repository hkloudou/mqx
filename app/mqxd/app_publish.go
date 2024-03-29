package main

import (
	"context"
	"log"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

func (m *app) OnClientPublish(s xtransport.Socket, p *mqtt.PublishPacket) {
	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	if !meta.Logined() {
		log.Println("not Logined")
		s.Close()
		return
	}
	if err := face.ValidateTopic(p.TopicName); err != nil {
		log.Println("valid topic", p.TopicName, err)
		return
	}
	// ACL interface
	if code := m._acl.GetPub(meta, p.TopicName); code != face.AclCodeAllow {
		log.Println("acl", code)
		s.Close()
		return
	}
	// publish to distributed nodes
	err := m._bridge.Publish(p)
	if p.Qos == 1 {
		res := mqtt.NewControlPacket(mqtt.Puback).(*mqtt.PubackPacket)
		res.MessageID = p.MessageID
		if err != nil {
			res.Qos = 0x80
		} else {
			res.Qos = p.Qos
		}
		s.Send(res)
	}
}

func (m *app) onPublish(p *mqtt.PublishPacket) {
	// retain
	if p.Retain {
		if m._retain == nil {
			return
		}
		if err := m._retain.Store(context.TODO(), p); err != nil {
			log.Println(err)
			return
		}
	}
	// log.Println("publish", p)
	// TODO: publish data to client and other node(include zero byte payload packet)
	sessions, err := m._session.Match(context.TODO(), p.TopicName)
	if err != nil {
		return
	}
	// connIDs := set.New(set.ThreadSafe)

	for i := 0; i < len(sessions); i++ {
		connID, found := m.sessionConns.Load(sessions[i])
		if !found {
			continue
		}
		go func(i2 int) {
			if _s, found := m.conns.Load(connID); found && _s != nil {
				if err2 := _s.(xtransport.Socket).Send(p); err2 != nil {
					log.Println("err send msg to", connID)
				}
			}
		}(i)
	}
}
