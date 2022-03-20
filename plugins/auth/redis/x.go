package redis

import "github.com/hkloudou/mqx/face"

type token struct {
	// SET	mqtt.$u.$c
	Req *face.AuthReply
}

func New() face.Auth {
	return nil
}

// auth.cmd: GET	mqtt.$u.$c
// list.cmd: KEYS 	mqtt.$u.*

// auth.redis.auth_cmd = HMGET mqtt_user:%u password
//         Copied!

// 你可以在命令中使用以下占位符，执行时 EMQX 将自动填充为客户端信息：

// %u：用户名
// %c：Client ID
// %C：TLS 证书公用名（证书的域名或子域名），仅当 TLS 连接时有效
// %d：TLS 证书 subject，仅当 TLS 连接时有效
/*

 */
