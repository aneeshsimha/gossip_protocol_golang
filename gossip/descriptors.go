package gossip

import "net"

type nodeDescriptor struct {
	address   net.IP
	timestamp uint64
}

type messageDescriptor struct {
	content   []byte
	timestamp uint64
}
