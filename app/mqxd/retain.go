package main

import (
	"context"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport/packets/mqtt"
)

func (m *app) checkRetain(meta *face.MetaInfo, patterns []string) ([]*mqtt.PublishPacket, error) {
	retaineds := make([]*mqtt.PublishPacket, 0)
	for i := 0; i < len(patterns); i++ {
		objs, err := m._retain.Check(context.TODO(), patterns[i])
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(objs); i++ {
			obj := objs[i]
			if code := m._acl.GetSub(meta, obj.TopicName); code == face.AclCodeAllow {
				retaineds = append(retaineds, obj)
			}
		}
	}
	return retaineds, nil
}
