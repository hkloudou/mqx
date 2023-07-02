package memory

import (
	"errors"

	"github.com/hkloudou/mqx/face"
)

// type model struct {
// 	Patterns []string `delim:","`
// 	Stream   bool
// }

func init() {
	face.RegisterPugin[face.Bridge]("memory", MustNew)
	// face.DefaultAuths["redis"] = MustNew
}

func MustNew(conf face.Conf) face.Bridge {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Bridge, error) {
	obj := &memoryBridge{}
	if conf == nil {
		return nil, errors.New("invalid conf")
	}
	err := conf.MapTo("bridge.plugin.memory", &obj.cfg)
	if err != nil {
		return nil, err
	}
	go obj.motion(nil)
	return obj, nil
}
