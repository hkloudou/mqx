package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/hkloudou/mqx/face"
	_ "github.com/hkloudou/mqx/plugins/acl/memory"
	_ "github.com/hkloudou/mqx/plugins/auth/memory"
	_ "github.com/hkloudou/mqx/plugins/auth/redis"
	_ "github.com/hkloudou/mqx/plugins/conf/ini"
	_ "github.com/hkloudou/mqx/plugins/retain/memory"
	_ "github.com/hkloudou/mqx/plugins/retain/redis"
	_ "github.com/hkloudou/mqx/plugins/session/memory"
	_ "github.com/hkloudou/mqx/plugins/session/redis"
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
	println(xcolor.Green(fmt.Sprintf("%-15s", "load auth")), xcolor.Yellow(_conf.MustString("auth", "plugin", "momory")))
	println(xcolor.Green(fmt.Sprintf("%-15s", "load retain")), xcolor.Yellow(_conf.MustString("retain", "plugin", "momory")))
	println(xcolor.Green(fmt.Sprintf("%-15s", "load session")), xcolor.Yellow(_conf.MustString("session", "plugin", "momory")))
	println(xcolor.Green(fmt.Sprintf("%-15s", "load acl")), xcolor.Yellow(_conf.MustString("acl", "plugin", "momory")))
	_auther := face.LoadPlugin[face.Auth](_conf.MustString("auth", "plugin", "momory"), _conf)
	_retain := face.LoadPlugin[face.Retain](_conf.MustString("retain", "plugin", "memory"), _conf)
	_session := face.LoadPlugin[face.Session](_conf.MustString("session", "plugin", "memory"), _conf)
	_acl := face.LoadPlugin[face.Acl](_conf.MustString("acl", "plugin", "memory"), _conf)
	_hook := newHook(_auther, _retain, _session, _acl)
	tran := transport.NewTransport("tcp", xtransport.Secure(false))
	l, err := tran.
		Listen(":1883")
	if err != nil {
		println(xcolor.Red(fmt.Sprintf("%v", err)))
		return
	}
	log.Println("pls login")
	_auther.Update(context.TODO(), &face.AuthRequest{
		UserName: "a",
		PassWord: "a",
		ClientId: "a",
	})
	// go func() {
	// 	log.Println("ready update")
	// 	time.Sleep(10 * time.Second)
	// 	log.Println("go")
	// 	_auther.Update(context.TODO(), &face.AuthRequest{
	// 		UserName: "a",
	// 		PassWord: "a",
	// 		ClientId: "b",
	// 	})
	// }()
	xruntime.PrintInfo()
	println()
	println(xcolor.Green(tran.String() + " listen on [" + l.Addr() + "]"))
	if err := l.Accept(func(sock xtransport.Socket) {
		defer func() {
			if r := recover(); r != nil {
				println(xcolor.Red(fmt.Sprintf("%v", r)))
			}
			sock.Close()
			_hook.OnClientDisConnected(sock)
		}()
		for {
			request, err := sock.Recv(func(r io.Reader) (interface{}, error) {
				i, err := mqtt.ReadPacket(r)
				return i, err
			})
			if err != nil {
				println(xcolor.Red(err.Error()))
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
			default:
				return
			}
		}
	}); err != nil {
		println(xcolor.Red(fmt.Sprintf("%v", err)))
		return
	}
}
