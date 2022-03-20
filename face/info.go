package face

import "crypto/tls"

type ClientInfo struct {
	IP              string
	Port            uint16
	Scheme          string /*tcp,tls,ws,wss*/
	ConnectionState *tls.ConnectionState
}

type NodeInfo struct {
	HostName string
	Port     uint16
}

type TransportInfo struct {
	ClientInfo
	NodeInfo
}
