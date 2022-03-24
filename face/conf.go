package face

import "time"

type Conf interface {
	MapTo(section string, source interface{}) error
	MustString(section string, key string, defaultVal string) string
	MustBool(section string, key string, defaultVal bool) bool
	MustInt(section string, key string, defaultVal int) int
	MustUint(section string, key string, defaultVal uint) uint
	MustDuration(section string, key string, defaultVal time.Duration) time.Duration
}
