package main

import (
	"context"
	"crypto/x509"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

const _keyFirstConnPacket = "status.firstpacket"
const _keyConnected = "status.connected"

func newHook(aurher face.Auth, retainer face.Retain) face.Hook {
	return &defaultHook{
		_auther:    aurher,
		_retainer:  retainer,
		conns:      sync.Map{},
		topicConns: sync.Map{},
	}
}

type defaultHook struct {
	_subHooks  []face.Hook
	_auther    face.Auth
	_retainer  face.Retain
	conns      sync.Map
	topicConns sync.Map
}

const _maxKeepAlive = (18 * time.Hour) + (12 * time.Minute) + (15 * time.Second)

func (m *defaultHook) OnClientConnect(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.ConnectPacket) {
	// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html
	// After a Network Connection is established by a Client to a Server, the first Packet sent from the Client to the Server MUST be a CONNECT Packet [MQTT-3.1.0-1].
	// A Client can only send the CONNECT Packet once over a Network Connection. The Server MUST process a second CONNECT Packet sent from a Client as a protocol violation and disconnect the Client [MQTT-3.1.0-2].  See section 4.8 for information about handling errors.
	// The payload contains one or more encoded fields. They specify a unique Client identifier for the Client, a Will topic, Will Message, User Name and Password. All but the Client identifier are optional and their presence is determined based on flags in the variable header.
	if s.Session().GetBool(_keyFirstConnPacket) {
		s.Close()
		return
	}
	s.Session().Set(_keyFirstConnPacket, true)
	/* The Keep Alive is a time interval measured in seconds. Expressed as a 16-bit word, it is the maximum time interval that is permitted to elapse between the point at which the Client finishes transmitting one Control Packet and the point it starts sending the next. It is the responsibility of the Client to ensure that the interval between Control Packets being sent does not exceed the Keep Alive value. In the absence of sending any other Control Packets, the Client MUST send a PINGREQ Packet [MQTT-3.1.2-23].
	   The Client can send PINGREQ at any time, irrespective of the Keep Alive value, and use the PINGRESP to determine that the network and the Server are working.
	   If the Keep Alive value is non-zero and the Server does not receive a Control Packet from the Client within one and a half times the Keep Alive time period, it MUST disconnect the Network Connection to the Client as if the network had failed [MQTT-3.1.2-24].
	   If a Client does not receive a PINGRESP Packet within a reasonable amount of time after it has sent a PINGREQ, it SHOULD close the Network Connection to the Server.
	   A Keep Alive value of zero (0) has the effect of turning off the keep alive mechanism. This means that, in this case, the Server is not required to disconnect the Client on the grounds of inactivity.
	   Note that a Server is permitted to disconnect a Client that it determines to be inactive or non-responsive at any time, regardless of the Keep Alive value provided by that Client.
	   Non normative comment
	   The actual value of the Keep Alive is application specific; typically this is a few minutes. The maximum value is 18 hours 12 minutes and 15 seconds. */
	if p.Keepalive != 0 {
		k := (time.Duration(float32(p.Keepalive) * 1.5)) * time.Second
		if k > _maxKeepAlive {
			k = _maxKeepAlive
		}
		s.SetTimeOut(k)
	} else {
		s.SetTimeOut(_maxKeepAlive)
	}

	// TODO: will message

	res := mqtt.NewControlPacket(mqtt.Connack).(*mqtt.ConnackPacket)
	// only mqtt3.1.1 supported
	if p.ProtocolName != "MQTT" || p.ProtocolVersion != 4 {
		res.ReturnCode = mqtt.ErrProtocolViolation
		m.OnClientConnack(s, p, res)
		return
	}
	// deny all connection if no auth plugin provided
	if m._auther == nil {
		res.ReturnCode = mqtt.ErrRefusedServerUnavailable
		m.OnClientConnack(s, p, res)
		return
	}
	req := &face.AuthRequest{
		ClientId: p.ClientIdentifier,
		UserName: p.Username,
		PassWord: string(p.Password),
		ClientIp: s.Remote(),
	}
	// read serverName and subject commonName from connection
	_readTls := func() {
		state := s.ConnectionState()
		if state != nil {
			req.TlsServerName = state.ServerName
			for _, cert := range state.PeerCertificates {
				if cert.IsCA {
					continue
				}
				for i := 0; i < len(cert.ExtKeyUsage); i++ {
					if cert.ExtKeyUsage[i] == x509.ExtKeyUsageClientAuth || cert.ExtKeyUsage[i] == x509.ExtKeyUsageAny {
						req.TlsSubjectName = cert.Subject.CommonName
						return
					}
				}
			}
		}
	}

	if p.Username == "mqtt" && string(p.Password) == "public" {
		res.ReturnCode = 0x00
	} else {
		_readTls()
		res.ReturnCode = m._auther.Check(context.TODO(), req)
	}
	m.OnClientConnack(s, p, res)
}

func (m *defaultHook) OnClientConnack(s xtransport.Socket[mqtt.ControlPacket], req *mqtt.ConnectPacket, ack *mqtt.ConnackPacket) {
	s.Send(ack)
	if ack.ReturnCode == mqtt.Accepted {
		m.OnClientConnected(s, req)
	}
}

func (m *defaultHook) OnClientPublish(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.PublishPacket) {
	if !s.Session().GetBool(_keyConnected) {
		s.Close()
		return
	}
	// log.Println("OnClientPublish", p.String())
	// TODO: ACL interface
	// TODO: retainer store
	// TODO: qos:2
	if p.Qos == 1 {
		res := mqtt.NewControlPacket(mqtt.Puback).(*mqtt.PubackPacket)
		res.MessageID = p.MessageID
		res.Qos = p.Qos
		s.Send(res)
	}
	if p.Retain {
		if m._retainer == nil {
			return
		}
		if err := m._retainer.Store(context.TODO(), p); err != nil {
			log.Println(err)
			return
		}
	}
	// TODO: publish data to client and other node
}

func (m *defaultHook) OnClientSubcribe(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.SubscribePacket) {
	if !s.Session().GetBool(_keyConnected) {
		s.Close()
		return
	}
	res := mqtt.NewControlPacket(mqtt.Suback).(*mqtt.SubackPacket)
	res.MessageID = p.MessageID

	log.Println("OnClientSubcribe", p.String())
	// verify
	if len(p.Qoss) != len(p.Topics) || len(p.Qoss) == 0 {
		ma := len(p.Qoss)
		if len(p.Topics) > len(p.Qoss) {
			ma = len(p.Topics)
		}
		res.ReturnCodes = make([]byte, ma)
		for i := 0; i < len(res.ReturnCodes); i++ {
			res.ReturnCodes[i] = 0x80
		}
		s.Send(res)
		return
	}
	res.ReturnCodes = make([]byte, len(p.Qoss))

	// TODO: ACL interface
	// TODO: retainer read
	retaineds := make([]*mqtt.PublishPacket, 0)
	for i := 0; i < len(p.Topics); i++ {
		if m._retainer == nil {
			log.Println("no retainer define")
			res.ReturnCodes[i] = 0x80
			continue
		}
		objs, err := m._retainer.Check(context.TODO(), p.Topics[i])
		if err != nil {
			println("err", err.Error())
			res.ReturnCodes[i] = 0x80
			continue
		}
		res.ReturnCodes[i] = 0x00
		if objs != nil && len(objs) > 0 {
			retaineds = append(retaineds, objs...)
		}
	}
	s.Send(res)
	for i := 0; i < len(retaineds); i++ {
		s.Send(retaineds[i])
	}
}

func (m *defaultHook) OnClientUnSubcribe(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.UnsubscribePacket) {
	if !s.Session().GetBool(_keyConnected) {
		s.Close()
		return
	}

	res := mqtt.NewControlPacket(mqtt.Unsuback).(*mqtt.UnsubackPacket)
	res.MessageID = p.MessageID
	s.Send(res)
}

func (m *defaultHook) OnClientConnected(s xtransport.Socket[mqtt.ControlPacket], req *mqtt.ConnectPacket) {
	connid := uuid.New().String()
	s.Session().Set("status.connid", connid)
	s.Session().Set("auth.username", req.Username)
	s.Session().Set("auth.clientid", req.ClientIdentifier)
	m.conns.Store(connid, s)

	s.Session().Set(_keyConnected, true)
	log.Println(connid, ">", "connected")
}

func (m *defaultHook) OnClientDisConnected(s xtransport.Socket[mqtt.ControlPacket]) {
	s.Session().Set(_keyConnected, false)
	se := s.Session()
	connid := se.GetString("status.connid")
	m.conns.Delete(connid)
	log.Println(connid, ">", "disConnected")
}
