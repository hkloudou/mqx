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

func newHook(aurher face.Auth) face.Hook {
	return &defaultHook{
		_auther: aurher,
		conns:   sync.Map{},
	}
}

type defaultHook struct {
	_subHooks []face.Hook
	_auther   face.Auth
	conns     sync.Map
}

const _maxKeepAlive = (18 * time.Hour) + (12 * time.Minute) + (15 * time.Second)

func (m *defaultHook) OnClientConnect(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.ConnectPacket) {
	// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html
	// After a Network Connection is established by a Client to a Server, the first Packet sent from the Client to the Server MUST be a CONNECT Packet [MQTT-3.1.0-1].
	// A Client can only send the CONNECT Packet once over a Network Connection. The Server MUST process a second CONNECT Packet sent from a Client as a protocol violation and disconnect the Client [MQTT-3.1.0-2].  See section 4.8 for information about handling errors.
	// The payload contains one or more encoded fields. They specify a unique Client identifier for the Client, a Will topic, Will Message, User Name and Password. All but the Client identifier are optional and their presence is determined based on flags in the variable header.
	// http://blog.mcxiaoke.com/mqtt/protocol/MQTT-3.1.1-CN.pdf
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
	}

	// TODO: will message

	//read serverName and subject commonName from connection
	res := mqtt.NewControlPacket(mqtt.Connack).(*mqtt.ConnackPacket)
	if m._auther == nil {
		res.ReturnCode = mqtt.ErrRefusedServerUnavailable
	} else {
		req := &face.AuthRequest{
			ClientId: p.ClientIdentifier,
			UserName: p.Username,
			PassWord: string(p.Password),
			ClientIp: s.Remote(),
		}
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
		res.ReturnCode = m._auther.Check(context.TODO(), req)
	}

	m.OnClientConnack(s, res)
	return
}

func (m *defaultHook) OnClientConnack(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.ConnackPacket) {
	s.Send(p)
	if p.ReturnCode == mqtt.Accepted {
		m.OnClientConnected(s)
	}
}

func (m *defaultHook) OnClientPublish(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.PublishPacket) {
	if s.Session().GetBool(_keyConnected) {
		s.Close()
		return
	}
}

func (m *defaultHook) OnClientSubcribe(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.SubscribePacket) {
	if s.Session().GetBool(_keyConnected) {
		s.Close()
		return
	}
}

func (m *defaultHook) OnClientUnSubcribe(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.UnsubscribePacket) {
	if s.Session().GetBool(_keyConnected) {
		s.Close()
		return
	}
}

func (m *defaultHook) OnClientConnected(s xtransport.Socket[mqtt.ControlPacket]) {
	s.Session().Set(_keyConnected, true)
	connID := uuid.New().String()
	s.Session().Set("connID", connID)
	m.conns.Store(connID, s)
	log.Println(connID, ">", "connected")
}

func (m *defaultHook) OnClientDisConnected(s xtransport.Socket[mqtt.ControlPacket]) {
	s.Session().Set(_keyConnected, false)
	se := s.Session()
	connID := se.GetString("connID")
	m.conns.Delete(connID)
	log.Println(connID, ">", "disConnected")
}