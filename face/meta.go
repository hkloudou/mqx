package face

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/hkloudou/xlib/xcolor"
)

type MetaInfo struct {
	ConnID           string
	Protocol         string
	UserName         string
	ClientIdentifier string
	ClientIP         net.IP
	ClientPort       uint16
	// Logined          bool
	SessionKey      string
	ConnectionState *tls.ConnectionState
}

func (m *MetaInfo) Stirng() string {
	return fmt.Sprintf("[%s] %-20s, connid = %s, cid = %s, user = %s", m.Protocol, xcolor.Yellow(fmt.Sprintf("%s:%d", m.ClientIP.String(), m.ClientPort)), xcolor.Blue(m.ConnID), xcolor.Blue(m.ClientIdentifier), xcolor.Blue(m.UserName))
}

func (m *MetaInfo) Logined() bool {
	return len(m.UserName) > 0
}
