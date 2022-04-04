package main

var mqttConfig struct {

	// mqtt.strict_mode
	// default: false
	StrictMode bool

	// mqtt.wildcard_subscription
	// default: true
	WildcardSubscription bool

	// mqtt.retain_available
	// default: true
	RetainAvailable bool
}
