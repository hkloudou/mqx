package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xlib/xcert"
	"github.com/hkloudou/xlib/xcolor"
	"github.com/hkloudou/xlib/xruntime"
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
	transport "github.com/hkloudou/xtransport/transports/tcp"
	wtransport "github.com/hkloudou/xtransport/transports/ws"
)

type app struct {
	_cfg         face.Conf
	_auth        face.Auth
	_retain      face.Retain
	_session     face.Session
	_acl         face.Acl
	_bridge      face.Bridge
	conns        sync.Map
	sessionConns sync.Map
	// connUserBook sync.Map
}

func (m *app) run() {
	xruntime.PrintInfo()
	println()
	wg := sync.WaitGroup{}
	xruntime.GoUnterminated(func() {
		if m._cfg.MustBool("tcp", "enable", false) {
			port := m._cfg.MustUint("tcp", "port", 1883)
			tran := transport.NewTransport("tcp", xtransport.Secure(false))
			l, err := tran.
				Listen(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("tcp listen err:%v", err.Error())))
				return
			}
			log.Println(xcolor.Green("tcp listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
			m.addTransport("tcp", l)
		}
	}, &wg, false, 1*time.Second)
	xruntime.GoUnterminated(func() {
		if m._cfg.MustBool("tls", "enable", false) {
			_cfg, err := xcert.ParseTlsConfig(m._cfg.MustString("tls", "tls_ca_path", ""), m._cfg.MustString("tls", "tls_cert_path", ""), m._cfg.MustString("tls", "tls_key_path", ""))
			port := m._cfg.MustUint("tls", "port", 8883)
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
			m.addTransport("tls", l)
		}
	}, &wg, false, 1*time.Second)
	xruntime.GoUnterminated(func() {
		if m._cfg.MustBool("ws", "enable", false) {
			port := m._cfg.MustUint("ws", "port", 80)
			tran := wtransport.NewTransport("/ws", xtransport.Secure(false))
			l, err := tran.
				Listen(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Println(xcolor.Red(fmt.Sprintf("ws listen err:%v", err.Error())))
				return
			}
			log.Println(xcolor.Green("ws listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
			m.addTransport("ws", l)
		}
	}, &wg, false, 1*time.Second)
	xruntime.GoUnterminated(func() {
		if m._cfg.MustBool("wss", "enable", false) {
			_cfg, err := xcert.ParseTlsConfig(m._cfg.MustString("wss", "tls_ca_path", ""), m._cfg.MustString("wss", "tls_cert_path", ""), m._cfg.MustString("wss", "tls_key_path", ""))
			port := m._cfg.MustUint("wss", "port", 8883)
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
			m.addTransport("wss", l)
		}
	}, &wg, false, 1*time.Second)

	xruntime.GoUnterminated(func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello")
		})
		// log.Println("hello")
		// mux.ServeHTTP()
		port := 10081
		log.Println(xcolor.Green("metric listen on"), xcolor.Green(fmt.Sprintf(":%d", port)))
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
		if err != nil {
			log.Println(xcolor.Red(fmt.Sprintf("metric listen err:%v", err.Error())))
			return
		}

	}, &wg, false, 1*time.Second)

	if err := m._bridge.Motion(func(obj *mqtt.PublishPacket) {
		m.publish(obj)
	}); err != nil {
		log.Println(xcolor.Red(fmt.Sprintf("bridge err:%v", err.Error())))
		time.Sleep(10 * time.Second)
		panic("err")
	}
	wg.Wait()
}
