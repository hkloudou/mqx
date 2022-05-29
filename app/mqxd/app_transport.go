package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

func (m *app) addTransport(protocol string, l xtransport.Listener) {
	if err := l.Accept(func(sock xtransport.Socket) {
		log.Println("accept connect from", sock.Remote())
		port, _ := strconv.ParseUint(strings.Split(sock.Remote()+":0", ":")[1], 10, 16)
		meta := &face.MetaInfo{
			Protocol:        protocol,
			ClientIP:        net.ParseIP(strings.Split(sock.Remote(), ":")[0]),
			ClientPort:      uint16(port),
			ConnectionState: sock.ConnectionState(),
			ConnID:          uuid.New().String(),
		}
		sock.Session().Set("meta", meta)
		m.conns.Store(meta.ConnID, sock)
		sock.SetTimeOut(time.Second * 5)
		//fitst packet must be arrive in next 5 seconds

		defer func() {
			if r := recover(); r != nil {
				log.Println(xcolor.Red(fmt.Sprintf("accept panic%v", r)))
			}
			sock.Close()
			m.OnClientDisConnected(sock)
		}()
		for {
			request, err := sock.Recv(func(r io.Reader) (interface{}, error) {
				i, err := mqtt.ReadPacket(r)
				return i, err
			})
			if err != nil {
				println("read err", xcolor.Red(err.Error()))
				return
			}
			if request == nil {
				continue
			}
			switch request := request.(type) {
			case mqtt.ControlPacket:
				if request.(mqtt.ControlPacket).Type() <= 0 || request.(mqtt.ControlPacket).Type() >= 14 {
					sock.Close()
					return
				}
				switch request.(mqtt.ControlPacket).Type() {
				case mqtt.Pingreq:
					sock.Send(mqtt.NewControlPacket(mqtt.Pingresp))
					break
				case mqtt.Connect:
					m.OnClientConnect(sock, request.(*mqtt.ConnectPacket))
					break
				case mqtt.Subscribe:
					m.OnClientSubcribe(sock, request.(*mqtt.SubscribePacket))
					break
				case mqtt.Unsubscribe:
					m.OnClientUnSubcribe(sock, request.(*mqtt.UnsubscribePacket))
					break
				case mqtt.Publish:
					m.OnClientPublish(sock, request.(*mqtt.PublishPacket))
					break
				default:
					log.Println(fmt.Errorf("not support packet type:%d", request.Type()))
				}
			default:
				return
			}
		}
	}); err != nil {
		log.Println(xcolor.Red(fmt.Sprintf("listen err:%v", err.Error())))
		return
	}
}
