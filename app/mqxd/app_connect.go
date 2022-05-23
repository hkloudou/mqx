package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/fatih/set"
	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

const _maxKeepAlive = (18 * time.Hour) + (12 * time.Minute) + (15 * time.Second)

func (m *app) getSessionConnections(sessionKey string) set.Interface {
	actual, _ := m.sessionConns.LoadOrStore(sessionKey, set.New(set.ThreadSafe))
	return actual.(set.Interface)
}

func (m *app) OnClientConnect(s xtransport.Socket, p *mqtt.ConnectPacket) {
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

	res.ReturnCode = p.Validate()
	if res.ReturnCode == mqtt.Accepted {
		res.ReturnCode = m._auth.Check(context.TODO(), req)
	}
	m.OnClientConnack(s, p, res)
}

func (m *app) OnClientConnack(s xtransport.Socket, req *mqtt.ConnectPacket, ack *mqtt.ConnackPacket) {
	s.Send(ack)
	if ack.ReturnCode == mqtt.Accepted {
		m.OnClientConnected(s, req)
	} else {
		time.Sleep(10 * time.Second)
		s.Close()
	}
}

func (m *app) OnClientConnected(s xtransport.Socket, req *mqtt.ConnectPacket) {

	meta := s.Session().MustGet("meta").(*face.MetaInfo)
	meta.UserName = req.Username
	meta.ClientIdentifier = req.ClientIdentifier

	if len(meta.ClientIdentifier) == 0 {

	}

	meta.SessionKey = meta.ClientIdentifier

	// meta.Logined = true

	fmt.Println(xcolor.Green("logined   "), meta.Stirng())
	m.getSessionConnections(meta.SessionKey).Add(meta.ConnID)

	// check clean session on logined
	if req.CleanSession {
		m._session.Clear(context.TODO(), meta.SessionKey)
	} else {
		// 1. check the non cleanSesion's subscribe
		patterns, err := m._session.ClientPatterns(context.TODO(), meta.SessionKey)
		if err != nil {
			s.Close()
		}
		// 2. check retain message
		retaineds, err := m.checkRetain(meta, patterns)
		if err != nil {
			s.Close()
		}
		// 3. send the retained message to the client side
		for i := 0; i < len(retaineds); i++ {
			if err := s.Send(retaineds[i]); err != nil {
				s.Close()
			}
		}
	}
}

func (m *app) OnClientDisConnected(s xtransport.Socket) {
	meta := s.Session().MustGet("meta").(*face.MetaInfo)

	//1. remove from sessionBook
	if len(meta.UserName) > 0 {
		book := m.getSessionConnections(meta.SessionKey)
		book.Remove(meta.ConnID)
	}
	//2. dele from conn lists
	m.conns.Delete(meta.ConnID)

	fmt.Println(xcolor.Red("disconnected"), meta.Stirng())
}
