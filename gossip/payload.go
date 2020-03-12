package gossip

type KeepAlivePayload struct {
	KnownNodes []nodeDescriptor
}

// utility function to make a separate list of descriptors to send to another node
func prepareKeepAlivePayload(nodes1 []nodeDescriptor, me nodeDescriptor) KeepAlivePayload {
	var newList []nodeDescriptor // nil
	for _, e := range nodes1 {
		if e.Address != nil && e.ID >= 100 {
			// 100 checks for sentinels
			newList = append(newList, e)
		}
	}
	newList = append(newList, me) // insert as last element
	kap := KeepAlivePayload{newList}
	return kap
}

type StringPayload struct {
	Message messageDescriptor
}

func newStringPayload(message messageDescriptor) StringPayload {
	return StringPayload{message}
}