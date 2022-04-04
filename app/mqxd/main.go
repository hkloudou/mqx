package main

import "sync"

func main() {
	_a := &app{
		conns:      sync.Map{},
		topicConns: sync.Map{},
	}
	_a.init()
	_a.run()
}
