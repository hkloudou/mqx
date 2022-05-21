package main

import (
	"context"
	"log"
	"sync"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

func (m *app) OnClientPublish(s xtransport.Socket, p *mqtt.PublishPacket) {
	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	if !meta.Connected {
		log.Println("not connected")
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

	// retainer store
	// TODO: qos:2
	if p.Retain {
		if m._retain == nil {
			return
		}
		if err := m._retain.Store(context.TODO(), p); err != nil {
			log.Println(err)
			return
		}
	}
	m.publish(p)
	once := sync.Once{}
	once.Do(func() {
		if p.Qos == 1 {
			res := mqtt.NewControlPacket(mqtt.Puback).(*mqtt.PubackPacket)
			res.MessageID = p.MessageID
			res.Qos = p.Qos
			s.Send(res)
		}
	})
}

func (m *app) publish(p *mqtt.PublishPacket) {
	log.Println("publish", p)
	// TODO: publish data to client and other node(include zero byte payload packet)
	clients, err := m._session.Match(context.TODO(), p.TopicName)
	if err != nil {
		return
	}
	log.Println("match", clients)
	for i := 0; i < len(clients); i++ {
		go func(i2 int) {
			if _s, found := m.conns.Load(clients[i2]); found && _s != nil {
				if err2 := _s.(xtransport.Socket).Send(p); err2 != nil {
					// log.Println("err send msg to", clients[i2])
				}
			}
		}(i)
	}
}
