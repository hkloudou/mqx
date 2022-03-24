package face

type Conf interface {
	MapTo(section string, source interface{}) error
}
