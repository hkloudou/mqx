package main

import (
	"context"
	"crypto/x509"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport"
	packet "github.com/hkloudou/xtransport/packets/mqtt"
)

func newHook() face.Hook {
	return &defaultHook{}
}

type defaultHook struct {
	_subHooks []face.Hook
	_auther   face.Auth
}

func (m *defaultHook) OnClientConnect(s xtransport.Socket[packet.ControlPacket], p *packet.ConnectPacket) {
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
	res := packet.NewControlPacket(packet.Connack).(*packet.ConnackPacket)
	if m._auther == nil {
		res.ReturnCode = packet.ErrRefusedServerUnavailable
	} else {
		res.ReturnCode = m._auther.Check(context.TODO(), req)
	}
	s.Send(p)
	m.OnClientConnack(s, res)
	return
}

func (m *defaultHook) OnClientConnack(s xtransport.Socket[packet.ControlPacket], p *packet.ConnackPacket) {

}

func (m *defaultHook) OnPublishConnect(s xtransport.Socket[packet.ControlPacket], p *packet.PublishPacket) {
}

func (m *defaultHook) OnClientSubcribe(s xtransport.Socket[packet.ControlPacket], p *packet.SubscribePacket) {
}

func (m *defaultHook) OnClientUnSubcribe(s xtransport.Socket[packet.ControlPacket], p *packet.UnsubscribePacket) {

}
