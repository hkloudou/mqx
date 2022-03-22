package face

import (
	"github.com/hkloudou/xtransport"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

// Hook message io
type Hook interface {
	OnClientConnect(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.ConnectPacket)
	OnClientConnack(s xtransport.Socket[mqtt.ControlPacket], req *mqtt.ConnectPacket, ack *mqtt.ConnackPacket)
	OnClientPublish(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.PublishPacket)
	OnClientSubcribe(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.SubscribePacket)
	OnClientUnSubcribe(s xtransport.Socket[mqtt.ControlPacket], p *mqtt.UnsubscribePacket)

	OnClientConnected(s xtransport.Socket[mqtt.ControlPacket], req *mqtt.ConnectPacket)
	OnClientDisConnected(s xtransport.Socket[mqtt.ControlPacket])
	// OnClientConnect(evt *HookEvtClientConnect)
	// OnClientPublish(evt *HookEvtClientPublish)
	// OnClientSubcribe(evt *HookEvtClientSubcribe)
	// OnClientUnSubcribe(evt *HookEvtClientUnSubcribe)

	// //
	// OnClientConnected(evt *HookEvtClientConnect)
	// OnClientDisConnected(evt *HookEvtClientConnect)
}

// type HookAction uint8

// const (
// 	HookActionClientConnect HookAction = 1
// 	// HookActionClientConnack    HookAction = 2
// 	HookActionClientPublish    HookAction = 3
// 	HookActionClientSubcribe   HookAction = 8
// 	HookActionClientUnsubcribe HookAction = 9

// 	//Client
// 	HookActionClientConnected    HookAction = 21
// 	HookActionClientDisConnected HookAction = 22
// 	//Session
// 	HookActionSessionSubscribed   HookAction = 31
// 	HookActionSessionUnSubscribed HookAction = 32
// )

// type hookClientInfo struct {
// 	ClientId string `json:"clientid"`
// 	UserName string `json:"username"`
// }

// type HookEvtClientConnect struct {
// 	hookClientInfo `json:",inline"`
// 	Ipaddress      string `json:"ipaddress"`
// 	Keepalive      int    `json:"keepalive"`
// 	ProtoVersion   int    `json:"protoVersion"`
// }

// type HookEvtClientPublish struct {
// 	hookClientInfo `json:",inline"`
// }

// type HookEvtClientDisconnected struct {
// 	hookClientInfo `json:",inline"`
// }

// type HookEvtClientSubcribe struct {
// 	hookClientInfo `json:",inline"`
// 	Topic          string `json:"topic"`
// 	Qos            uint8  `json:"qos"`
// }

// type HookEvtClientUnSubcribe struct {
// 	hookClientInfo `json:",inline"`
// 	Topic          string `json:"topic"`
// }
