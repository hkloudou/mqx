package main

import (
	"container/list"
	"sync"

	"github.com/hkloudou/xtransport/packets/mqtt"
)

//https://github.com/eclipse/paho.mqtt.golang/blob/master/router.go

type MessageHandler func()

// route is a type which associates MQTT Topic strings with a
// callback to be executed upon the arrival of a message associated
// with a subscription to that topic.
type route struct {
	topic    string
	callback MessageHandler
}

type router struct {
	sync.RWMutex
	routes         *list.List
	defaultHandler MessageHandler
	messages       chan *mqtt.PublishPacket
}

func newRouter() *router {
	router := &router{routes: list.New(), messages: make(chan *mqtt.PublishPacket)}
	return router
}

// addRoute takes a topic string and MessageHandler callback. It looks in the current list of
// routes to see if there is already a matching Route. If there is it replaces the current
// callback with the new one. If not it add a new entry to the list of Routes.
func (r *router) addRoute(topic string, callback MessageHandler) {
	r.Lock()
	defer r.Unlock()
	for e := r.routes.Front(); e != nil; e = e.Next() {
		if e.Value.(*route).topic == topic {
			r := e.Value.(*route)
			r.callback = callback
			return
		}
	}
	r.routes.PushBack(&route{topic: topic, callback: callback})
}

// deleteRoute takes a route string, looks for a matching Route in the list of Routes. If
// found it removes the Route from the list.
func (r *router) deleteRoute(topic string) {
	r.Lock()
	defer r.Unlock()
	for e := r.routes.Front(); e != nil; e = e.Next() {
		if e.Value.(*route).topic == topic {
			r.routes.Remove(e)
			return
		}
	}
}

// setDefaultHandler assigns a default callback that will be called if no matching Route
// is found for an incoming Publish.
func (r *router) setDefaultHandler(handler MessageHandler) {
	r.Lock()
	defer r.Unlock()
	r.defaultHandler = handler
}
