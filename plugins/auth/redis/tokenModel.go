package redis

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
)

type tokenModel struct {
	Key           string
	CreateAt      uint64
	TokenPassword string
}

func (m *tokenModel) MarshalBinary() (data []byte, err error) {

	//data
	var _s = make([]byte, len(m.TokenPassword))
	copy(_s, m.TokenPassword)
	for i := 0; i < len(_s); i++ {
		_s[i] = _s[i] ^ byte(i&0xFF)
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf[0:], m.CreateAt)

	buf = append(buf, _s...)
	return buf, nil
}
func (m *tokenModel) Value() (driver.Value, error) {
	// encoding.BinaryMarshaler
	// bt, err := json.Marshal(m)
	return m.MarshalBinary()
}
func (j *tokenModel) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal tokenModel value:", value))
	}
	return j.UnmarshalBinary(bytes)
}
func (j *tokenModel) UnmarshalBinary(bytes []byte) error {
	if len(bytes) < 8 {
		return errors.New(fmt.Sprint("Failed to unmarshal tokenModel value:", bytes))
	}
	var obj tokenModel
	obj.CreateAt = binary.BigEndian.Uint64(bytes)
	bytes = bytes[8:]
	var _s = make([]byte, len(bytes))
	copy(_s, bytes)
	for i := 0; i < len(_s); i++ {
		_s[i] = _s[i] ^ byte(i&0xFF)
	}
	obj.TokenPassword = string(_s)
	*j = obj
	return nil
}

// ByAge implements sort.Interface based on the Age field.
type tokenModelByCreateAd []tokenModel

func (a tokenModelByCreateAd) Len() int           { return len(a) }
func (a tokenModelByCreateAd) Less(i, j int) bool { return a[i].CreateAt < a[j].CreateAt }
func (a tokenModelByCreateAd) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
