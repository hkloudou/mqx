package main

import (
	"context"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

func (m *app) OnClientSubcribe(s xtransport.Socket, p *mqtt.SubscribePacket) {
	meta := getMeta(s)
	if !meta.Connected {
		s.Close()
		return
	}

	// verify request
	if err := p.Validate(); err != nil {
		s.Close()
		return
	}
	if mqttConfig.StrictMode {
		if err := p.StrictValidate(); err != nil {
			s.Close()
			return
		}
	}

	// check acl
	var accessedTopics = make([]string, 0)
	reply := mqtt.NewControlPacket(mqtt.Suback).(*mqtt.SubackPacket)
	reply.MessageID = p.MessageID
	reply.ReturnCodes = make([]byte, len(p.Qoss))
	for i := 0; i < len(p.Qoss); i++ {
		if code := m._acl.GetSub(meta, p.Topics[i]); code == face.AclCodeAllow {
			reply.ReturnCodes[i] = p.Qoss[i]
			accessedTopics = append(accessedTopics, p.Topics[i])
		} else {
			reply.ReturnCodes[i] = 0x80
		}
	}
	s.Send(reply)

	// store session
	if err := m._session.Add(context.Background(), meta.ClientIdentifier, accessedTopics...); err != nil {
		//TODO: system error
		s.Close()
		return
	}

	// check retain on subscribe
	if retaineds, err := m.checkRetain(meta, accessedTopics); err != nil {
		//TODO: system error
		s.Close()
		return
	} else {
		for i := 0; i < len(retaineds); i++ {
			if err := s.Send(retaineds[i]); err != nil {
				s.Close()
			}
		}
	}
}

func (m *app) OnClientUnSubcribe(s xtransport.Socket, p *mqtt.UnsubscribePacket) {
	meta := getMeta(s)
	if !meta.Connected {
		s.Close()
		return
	}

	// verify request
	if err := p.Validate(); err != nil {
		s.Close()
		return
	}
	if mqttConfig.StrictMode {
		if err := p.StrictValidate(); err != nil {
			s.Close()
			return
		}
	}
	res := mqtt.NewControlPacket(mqtt.Unsuback).(*mqtt.UnsubackPacket)
	res.MessageID = p.MessageID
	if err := m._session.Remove(context.TODO(), meta.ClientIdentifier, p.Topics...); err != nil {
		//TODO: system error
	}
	s.Send(res)
}
