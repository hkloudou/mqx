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

func RegisterPugin[T any, P any](provider string, fc func(P) T) {
	key := reflect.TypeOf(fc).String()
	// log.Println("key", key)
	actual, _ := plugins.LoadOrStore(key, make(map[string]func(P) T))
	actual.(map[string]func(P) T)[provider] = fc
}

func GetPlugin[T any, P any](provider string) (ret func(P) T) {
	key := reflect.TypeOf(ret).String()
	actual, ok := plugins.Load(key)
	if !ok || actual == nil {
		return nil
	}
	if fc, found := actual.(map[string]func(P) T)[provider]; found {
		return fc
	}
	return nil
}

func LoadPlugin[T any, P any](provider string, conf P) T {
	fc := GetPlugin[T, P](provider)
	if fc == nil {
		panic(fmt.Sprintf("fail load plugin:[%v] %s", reflect.TypeOf(fc), provider))
	}
	return fc(conf)
}
