package face

import "github.com/hkloudou/xtransport"

type Acl interface {
	Subcribe(s xtransport.Socket, qos byte, retain bool, pattern string) (bool, error)
	Publish(s xtransport.Socket, qos byte, retain bool, pattern string) (bool, error)
}

/*
flow:
1、check allow
2、diacard deny
*/
