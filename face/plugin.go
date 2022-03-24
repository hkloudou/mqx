package face

import (
	"fmt"
	"reflect"
	"sync"
)

// var DefaultAuths = map[string]func(Conf) Auth{}
// var DefaultRetains = map[string]func(Conf) Retain{}
// var DefaultAcls = map[string]func(Conf) ACL{}
var plugins = sync.Map{}

func AddPugin[T any](name string, fc func(Conf) T) {
	key := reflect.TypeOf(fc).String()
	actual, _ := plugins.LoadOrStore(key, make(map[string]func(Conf) T))
	actual.(map[string]func(Conf) T)[name] = fc
}

func GetPlugin[T any](name string) (ret func(Conf) T) {
	key := reflect.TypeOf(ret).String()
	actual, ok := plugins.Load(key)
	if !ok || actual == nil {
		return nil
	}
	if fc, found := actual.(map[string]func(Conf) T)[name]; found {
		return fc
	}
	return nil
}

func LoadPlugin[T any](name string, conf Conf) T {
	fc := GetPlugin[T](name)
	if fc == nil {
		panic(fmt.Sprintf("fail load plugin:[%v] %s", reflect.TypeOf(fc), name))
	}
	return fc(conf)
}
