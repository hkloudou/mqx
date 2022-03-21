package redis

import (
	"strings"

	"github.com/hkloudou/mqx/face"
)

func (m *redisAuther) parseKeyFromReq(req *face.AuthRequest, tmpl string) string {
	str := tmpl
	str = strings.ReplaceAll(str, "$p", m.opts.prefix)
	str = strings.ReplaceAll(str, "$u", req.UserName)
	str = strings.ReplaceAll(str, "$c", req.ClientId)
	str = strings.ReplaceAll(str, "$h", req.TlsServerName)
	str = strings.ReplaceAll(str, "$d", req.TlsSubjectName)
	return str
}

// auth.cmd: GET	mqtt.$u.$c
// list.cmd: KEYS 	mqtt.$u.*

// auth.redis.auth_cmd = HMGET mqtt_user:%u password
//         Copied!

// 你可以在命令中使用以下占位符，执行时 EMQX 将自动填充为客户端信息：

// $p prefix
// %u：UserName	%U Encoded User
// %c：Client ID
// %h: TLS.ServerName(HostName) %H Encoded
// %s: TLS.subject.CommonName  $S encoded
// %C：TLS 证书公用名（证书的域名或子域名），仅当 TLS 连接时有效
// %d：TLS 证书 subject，仅当 TLS 连接时有效
/*

 */
