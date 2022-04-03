package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

// const _keyFirstConnPacket = "status.firstpacket"

// const _keyConnected = "status.connected"

func newHook(aurh face.Auth, retain face.Retain, session face.Session, acl face.Acl) face.Hook {
	tmp := &defaultHook{
		_auth:    aurh,
		_retain:  retain,
		_session: session,
		_acl:     acl,
		conns:    sync.Map{},
	}
	// kick connect
	go func() {
		tmp._auth.MotionExpired(func(userName, clientId string) error {
			if actur, found := tmp.conns.Load(clientId); found {
				return actur.(xtransport.Socket).Close()
			}
			return nil
		})
	}()
	return tmp
}

type defaultHook struct {
	_auth      face.Auth
	_retain    face.Retain
	_session   face.Session
	_acl       face.Acl
	conns      sync.Map
	topicConns sync.Map
}

const _maxKeepAlive = (18 * time.Hour) + (12 * time.Minute) + (15 * time.Second)

func (m *defaultHook) OnClientConnect(s xtransport.Socket, p *mqtt.ConnectPacket) {
	// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html
	// After a Network Connection is established by a Client to a Server, the first Packet sent from the Client to the Server MUST be a CONNECT Packet [MQTT-3.1.0-1].
	// A Client can only send the CONNECT Packet once over a Network Connection. The Server MUST process a second CONNECT Packet sent from a Client as a protocol violation and disconnect the Client [MQTT-3.1.0-2].  See section 4.8 for information about handling errors.
	// The payload contains one or more encoded fields. They specify a unique Client identifier for the Client, a Will topic, Will Message, User Name and Password. All but the Client identifier are optional and their presence is determined based on flags in the variable header.
	if s.Session().GetBool("status.firstpacket") {
		s.Close()
		return
	}
	s.Session().Set("status.firstpacket", true)
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
	if m._auth == nil {
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

	_readTls()
	res.ReturnCode = m._auth.Check(context.TODO(), req)
	m.OnClientConnack(s, p, res)
}

func (m *defaultHook) OnClientConnack(s xtransport.Socket, req *mqtt.ConnectPacket, ack *mqtt.ConnackPacket) {
	s.Send(ack)
	if ack.ReturnCode == mqtt.Accepted {
		m.OnClientConnected(s, req)
	} else {
		time.Sleep(10 * time.Second)
		s.Close()
	}
}

func (m *defaultHook) OnClientPublish(s xtransport.Socket, p *mqtt.PublishPacket) {
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
	if enable, err := m._acl.Publish(s, p.Qos, p.Retain, p.TopicName); err != nil {
		log.Println("acl", err)
		s.Close()
		return
	} else if !enable {
		log.Println("enable", p.TopicName, enable)
		s.Close()
		return
	}
	once := sync.Once{}
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

	// TODO: publish data to client and other node(include zero byte payload packet)
	clients, err := m._session.Match(context.TODO(), p.TopicName)
	if err != nil {
		return
	}
	// log.Println("match", clients)
	for i := 0; i < len(clients); i++ {
		go func(i2 int) {
			if _s, found := m.conns.Load(clients[i2]); found && _s != nil {
				// match second times
				// don;t check again,let acl plugin do this
				// if face.IsPrivateTopic(p.TopicName) {
				// 	// retain private topic
				// 	if !face.MatchPrivateTopic(p.TopicName, "$uid", _s.(xtransport.Socket).Session().GetString("auth.clientid")) &&
				// 		!face.MatchPrivateTopic(p.TopicName, "$usr", _s.(xtransport.Socket).Session().GetString("auth.username")) {
				// 		// log.Println("un hit topic", p.TopicName)
				// 		return
				// 	}
				// }
				// log.Println("try send to", clients[i2])
				if err2 := _s.(xtransport.Socket).Send(p); err2 != nil {
					// log.Println("err send msg to", clients[i2])
				}
				once.Do(func() {
					if p.Qos == 1 {
						res := mqtt.NewControlPacket(mqtt.Puback).(*mqtt.PubackPacket)
						res.MessageID = p.MessageID
						res.Qos = p.Qos
						s.Send(res)
					}
				})
			}
		}(i)
	}
}

func (m *defaultHook) OnClientSubcribe(s xtransport.Socket, p *mqtt.SubscribePacket) {
	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	if !meta.Connected {
		s.Close()
		return
	}
	res := mqtt.NewControlPacket(mqtt.Suback).(*mqtt.SubackPacket)
	res.MessageID = p.MessageID

	// log.Println("OnClientSubcribe", p.String())
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

	// TODO: ACL interface
	// m._acl.Subcribe(s.Session().GetString("auth.username"))
	res.ReturnCodes = make([]byte, len(p.Qoss))
	var enable bool
	var err error
	var retaineds []*mqtt.PublishPacket
	for i := 0; i < len(p.Qoss); i++ {
		enable, err = m._acl.Subcribe(s, p.Qoss[i], p.Retain, p.Topics[i])
		if err != nil {
			break
		}
		if !enable {
			break
		}
	}

	// check retain on subscribe
	if err == nil {
		retaineds, err = m.checkRetain(s, p.Topics)
	}
	if err == nil {
		err = m._session.Add(context.Background(), meta.ClientIdentifier, p.Topics...)
	}
	if err != nil || !enable {
		res.ReturnCodes = make([]byte, len(p.Qoss))
		for i := 0; i < len(res.ReturnCodes); i++ {
			res.ReturnCodes[i] = 0x80
		}
	} else {
		for i := 0; i < len(res.ReturnCodes); i++ {
			res.ReturnCodes[i] = p.Qoss[i]
		}
	}

	s.Send(res)
	for i := 0; i < len(retaineds); i++ {
		if err := s.Send(retaineds[i]); err != nil {
			s.Close()
		}
	}
}

func (m *defaultHook) OnClientUnSubcribe(s xtransport.Socket, p *mqtt.UnsubscribePacket) {
	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	if !meta.Connected {
		s.Close()
		return
	}

	res := mqtt.NewControlPacket(mqtt.Unsuback).(*mqtt.UnsubackPacket)
	res.MessageID = p.MessageID
	if err := m._session.Remove(context.TODO(), meta.ClientIdentifier, p.Topics...); err != nil {
		log.Println("un suberr", err.Error())
	}
	s.Send(res)
}

func (m *defaultHook) OnClientConnected(s xtransport.Socket, req *mqtt.ConnectPacket) {
	// connid := uuid.New().String()
	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	meta.UserName = req.Username
	meta.ClientIdentifier = req.ClientIdentifier
	meta.Connected = true
	m.conns.Store(req.ClientIdentifier, s)
	fmt.Println(xcolor.Green("connected   "), meta.Stirng())
	if req.CleanSession {
		m._session.Clear(context.TODO(), req.ClientIdentifier)
	} else {
		patterns, err := m._session.ClientPatterns(context.TODO(), req.ClientIdentifier)
		if err != nil {
			s.Close()
		}
		// check retain on connected
		retaineds, err := m.checkRetain(s, patterns)
		if err != nil {
			s.Close()
		}
		for i := 0; i < len(retaineds); i++ {
			if err := s.Send(retaineds[i]); err != nil {
				s.Close()
			}
		}
	}
}

func (m *defaultHook) OnClientDisConnected(s xtransport.Socket) {
	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	meta.Connected = false
	m.conns.Delete(meta.ClientIdentifier)
	// log.Println(meta.ClientIdentifier, ">", "disConnected")
	// fmt.Println(xcolor.Green("connected"), meta.Stirng())
	fmt.Println(xcolor.Red("disconnected"), meta.Stirng())
	// disConnected
}

func (m *defaultHook) checkRetain(s xtransport.Socket, patterns []string) ([]*mqtt.PublishPacket, error) {
	retaineds := make([]*mqtt.PublishPacket, 0)
	for i := 0; i < len(patterns); i++ {
		objs, err := m._retain.Check(context.TODO(), patterns[i])
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(objs); i++ {
			obj := objs[i]
			if b, err := m._acl.Subcribe(s, 0, false, obj.TopicName); err == nil && b {
				retaineds = append(retaineds, obj)
			}
		}
	}
	return retaineds, nil
}
