package face

// store interface:
// 1.store clientID<=>SUBSCRIBE
// 2.store qos>0 message
// 3.store
type Store interface {
	Subscribe(clientid string, pattern []string, maxQos uint8)
	UnSubscribe(clientid string, pattern []string)
	Clear(clientid string)
}
