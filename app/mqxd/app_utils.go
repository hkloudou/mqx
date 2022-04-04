package main

import (
	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/xtransport"
)

func getMeta(s xtransport.Socket) *face.MetaInfo {
	return s.Session().MustGet("meta").(*face.MetaInfo)
}
