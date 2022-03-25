package main

import (
	"fmt"

	"github.com/hkloudou/mqx/face"
	_ "github.com/hkloudou/mqx/plugins/auth/redis"
	_ "github.com/hkloudou/mqx/plugins/conf/ini"
	_ "github.com/hkloudou/mqx/plugins/retain/redis"
	_ "github.com/hkloudou/mqx/plugins/session/memory"
	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xlib/xruntime"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
	transport "github.com/hkloudou/xtransport/transports/tcp"
)

func loadPlugin() {
	// key := _conf.MustString("auth", "plugin", "")
}

func main() {
	// _conf := conf.MustNew("")
	_conf := face.LoadPlugin[face.Conf]("ini", "")
	_auther := face.LoadPlugin[face.Auth](_conf.MustString("auth", "plugin", "momory"), _conf)
	_retain := face.LoadPlugin[face.Retain](_conf.MustString("auth", "plugin", "memory"), _conf)
	_session := face.LoadPlugin[face.Session](_conf.MustString("auth", "plugin", "memory"), _conf)
	_hook := newHook(_auther, _retain, _session)
	tran := transport.NewTransport[mqtt.ControlPacket]("tcp", xtransport.Secure(false))
	l, err := tran.
		Listen(":1883")
	if err != nil {
		println(xcolor.Red(fmt.Sprintf("%v", err)))
		return
	}
	xruntime.PrintInfo()
	println()
	println(xcolor.Green(tran.String() + " listen on [" + l.Addr() + "]"))
	if err := l.Accept(func(sock xtransport.Socket[mqtt.ControlPacket]) {
		defer func() {
			if r := recover(); r != nil {
				println(xcolor.Red(fmt.Sprintf("%v", r)))
			}
			sock.Close()
			_hook.OnClientDisConnected(sock)
		}()
		for {
			request, err := sock.Recv(mqtt.ReadPacket)

			if err != nil {
				println(xcolor.Red(err.Error()))
				return
			}
			// log.Println("recv", request.String())
			if request.Type() <= 0 || request.Type() >= 14 {
				sock.Close()
				return
			}
			switch request.Type() {
			case mqtt.Pingreq:
				sock.Send(mqtt.NewControlPacket(mqtt.Pingresp))
				break
			case mqtt.Connect:
				_hook.OnClientConnect(sock, request.(*mqtt.ConnectPacket))
				break
			case mqtt.Subscribe:
				_hook.OnClientSubcribe(sock, request.(*mqtt.SubscribePacket))
				break
			case mqtt.Unsubscribe:
				_hook.OnClientUnSubcribe(sock, request.(*mqtt.UnsubscribePacket))
				break
			case mqtt.Publish:
				_hook.OnClientPublish(sock, request.(*mqtt.PublishPacket))
				break
			default:
				// return nil, fmt.Errorf("not support packet type:%d", data.Type())
			}
		}
	}); err != nil {
		println(xcolor.Red(fmt.Sprintf("%v", err)))
		return
	}
}
