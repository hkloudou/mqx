package face

type Auth interface {
	Init()
	Login(info *TransportInfo, clientId string, userName string, passWord string)
	Check(info *TransportInfo, clientId string, userName string, passWord string)
}
