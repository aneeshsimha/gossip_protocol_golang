package gossip

//const (
//	KeepAlive = -1
//	String    = -2
//)

//type Payload interface {
//	ProtocolCode() int
//	Timestamp() uint64
//}

type KeepAlivePayload struct {
	KnownNodes []nodeDescriptor
}

//func (p *KeepAlivePayload) ProtocolCode() int {
//	return KeepAlive
//}

type StringPayload struct {
	Message messageDescriptor
}

//func (p *StringPayload) ProtocolCode() int {
//	return String
//}
