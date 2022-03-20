package face

// Auth auth interface
// mqtt.ConnectPacket
type Auth interface {
	Init()
	Save(info *TransportInfo, ClientIdentifier string, userName string, passWord string)
	Check(info *TransportInfo, ClientIdentifier string, userName string, passWord string)
}
