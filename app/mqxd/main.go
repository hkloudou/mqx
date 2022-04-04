package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hkloudou/mqx/face"
	_ "github.com/hkloudou/mqx/plugins/acl/memory"
	_ "github.com/hkloudou/mqx/plugins/auth/memory"
	_ "github.com/hkloudou/mqx/plugins/auth/redis"
	_ "github.com/hkloudou/mqx/plugins/conf/ini"
	_ "github.com/hkloudou/mqx/plugins/retain/memory"
	_ "github.com/hkloudou/mqx/plugins/retain/redis"
	_ "github.com/hkloudou/mqx/plugins/session/memory"
	_ "github.com/hkloudou/mqx/plugins/session/redis"
	"github.com/hkloudou/xlib/xcert"
	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xlib/xruntime"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
	transport "github.com/hkloudou/xtransport/transports/tcp"
	wtransport "github.com/hkloudou/xtransport/transports/ws"
)

var _hook face.Hook

func loadPlugin() {
	// key := _conf.MustString("auth", "plugin", "")
}

func addTransport(protocol string, l xtransport.Listener) {
	if err := l.Accept(func(sock xtransport.Socket) {
		port, _ := strconv.ParseUint(strings.Split(sock.Remote()+":0", ":")[1], 10, 16)
		meta := &face.MetaInfo{
			Protocol:        protocol,
			ClientIP:        net.ParseIP(strings.Split(sock.Remote(), ":")[0]),
			ClientPort:      uint16(port),
			Connected:       false,
			ConnectionState: sock.ConnectionState(),
		}
		sock.Session().Set("meta", meta)
		sock.SetTimeOut(time.Second * 5)
		//fitst packet must be arrive in next 5 seconds

		defer func() {
			if r := recover(); r != nil {
				log.Println(xcolor.Red(fmt.Sprintf("accept panic%v", r)))
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
				// println(xcolor.Red(err.Error()))
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
					log.Println(fmt.Errorf("not support packet type:%d", request.Type()))
				}
			default:
				return
			}
		}
	}); err != nil {
		log.Println(xcolor.Red(fmt.Sprintf("listen err:%v", err.Error())))
		// panic(err)
		return
	}
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
	_hook = newHook(_auther, _retain, _session, _acl)
	xruntime.PrintInfo()
	println()
	wg := sync.WaitGroup{}
	xruntime.GoUnterminated(func() {
		if _conf.MustBool("tcp", "enable", false) {
			port := _conf.MustUint("tcp", "port", 1883)
			tran := transport.NewTransport("tcp", xtransport.Secure(false))
			l, err := tran.
				Listen(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("tcp listen err:%v", err.Error())))
				return
			}
			log.Println(xcolor.Green("tcp listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
			addTransport("tcp", l)
		}
	}, &wg, false, 1*time.Second)
	xruntime.GoUnterminated(func() {
		if _conf.MustBool("tls", "enable", false) {
			_cfg, err := xcert.ParseTlsConfig(_conf.MustString("tls", "tls_ca_path", ""), _conf.MustString("tls", "tls_cert_path", ""), _conf.MustString("tls", "tls_key_path", ""))
			port := _conf.MustUint("tls", "port", 8883)
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("tls parse err:%v", err.Error())))
				return
			}
			tran := transport.NewTransport("tcp", xtransport.Secure(true), xtransport.TLSConfig(_cfg))
			l, err := tran.
				Listen(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("tls listen err:%v", err.Error())))
				return
			}
			log.Println(xcolor.Green("tls listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
			addTransport("tls", l)
		}
	}, &wg, false, 1*time.Second)
	xruntime.GoUnterminated(func() {
		if _conf.MustBool("ws", "enable", false) {
			port := _conf.MustUint("ws", "port", 80)
			tran := wtransport.NewTransport("/ws", xtransport.Secure(false))
			l, err := tran.
				Listen(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("ws listen err:%v", err.Error())))
				return
			}
			log.Println(xcolor.Green("ws listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
			addTransport("ws", l)
		}
	}, &wg, false, 1*time.Second)
	xruntime.GoUnterminated(func() {
		if _conf.MustBool("wss", "enable", false) {
			_cfg, err := xcert.ParseTlsConfig(_conf.MustString("wss", "tls_ca_path", ""), _conf.MustString("wss", "tls_cert_path", ""), _conf.MustString("wss", "tls_key_path", ""))
			port := _conf.MustUint("wss", "port", 8883)
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("wss parse err:%v", err.Error())))
				return
			}
			tran := wtransport.NewTransport("/ws", xtransport.Secure(true), xtransport.TLSConfig(_cfg))
			l, err := tran.
				Listen(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("wss listen err:%v", err.Error())))
				return
			}
			log.Println(xcolor.Green("wss listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
			addTransport("wss", l)
		}
	}, &wg, false, 1*time.Second)
	wg.Wait()
}
