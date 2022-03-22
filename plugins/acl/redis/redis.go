package redis

import "fmt"

func New(options ...Option) (interface{}, error) {
	opts := Options{
		addr:      "localhost:6379",
		db:        3,
		prefix:    "mqtt.acl",
		blackTmpl: "$p:$u:$c",
		whiteTmpl: "$p:$u:*",
	}
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
	}
	return nil, fmt.Errorf("")
}
