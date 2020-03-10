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

// utility functions to make a separate list of descriptors to send to another node
func preparePayload(nodes1 []nodeDescriptor, me nodeDescriptor) KeepAlivePayload {
	newList := make([]nodeDescriptor, len(nodes1)+1)
	copy(newList, nodes1)
	newList[len(nodes1)] = me // insert as last element
	kap := KeepAlivePayload{newList}
	return kap
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
