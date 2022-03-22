package face

import (
	"github.com/hkloudou/xtransport"
	packet "github.com/hkloudou/xtransport/packets/mqtt"
)

// Hook message io
type Hook interface {
	OnClientConnect(s xtransport.Socket[packet.ControlPacket], p *packet.ConnectPacket)
	OnClientConnack(s xtransport.Socket[packet.ControlPacket], p *packet.ConnackPacket)
	OnPublishConnect(s xtransport.Socket[packet.ControlPacket], p *packet.PublishPacket)
	OnClientSubcribe(s xtransport.Socket[packet.ControlPacket], p *packet.SubscribePacket)
	OnClientUnSubcribe(s xtransport.Socket[packet.ControlPacket], p *packet.UnsubscribePacket)

	// OnClientConnected(s xtransport.Socket[packet.ControlPacket], p *packet.ConnectPacket)
	// OnClientDisConnected(s xtransport.Socket[packet.ControlPacket])
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
