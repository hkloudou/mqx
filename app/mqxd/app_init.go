package main

import (
	"fmt"

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
)

func (m *app) init() {
	m._cfg = face.LoadPlugin[face.Conf]("ini", "")
	println(xcolor.Green(fmt.Sprintf("%-11s", "auth")), ":", xcolor.Yellow(m._cfg.MustString("auth", "plugin", "momory")))
	println(xcolor.Green(fmt.Sprintf("%-11s", "retain")), ":", xcolor.Yellow(m._cfg.MustString("retain", "plugin", "momory")))
	println(xcolor.Green(fmt.Sprintf("%-11s", "session")), ":", xcolor.Yellow(m._cfg.MustString("session", "plugin", "momory")))
	println(xcolor.Green(fmt.Sprintf("%-11s", "acl")), ":", xcolor.Yellow(m._cfg.MustString("acl", "plugin", "momory")))

	m._auth = face.LoadPlugin[face.Auth](m._cfg.MustString("auth", "plugin", "momory"), m._cfg)
	m._retain = face.LoadPlugin[face.Retain](m._cfg.MustString("retain", "plugin", "memory"), m._cfg)
	m._session = face.LoadPlugin[face.Session](m._cfg.MustString("session", "plugin", "memory"), m._cfg)
	m._acl = face.LoadPlugin[face.Acl](m._cfg.MustString("acl", "plugin", "memory"), m._cfg)
	m._cfg.MapTo("", &m.cfg)
	println(m.cfg.String())
}
